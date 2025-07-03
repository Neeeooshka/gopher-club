package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/internal/storage/postgres/sqlc"
	"github.com/Neeeooshka/gopher-club/pkg/logger/zap"
	"github.com/jackc/pgx/v5"
)

func (s *Postgres) ListWaitingOrders(ctx context.Context) ([]models.Order, error) {

	results, err := s.sqlc.ListWaitingOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing waiting orders: %w", err)
	}

	return s.extractOrders(results), nil
}

func (s *Postgres) UpdateOrders(ctx context.Context, orders []models.Order) error {

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			zap.Log.Debug("failed to rollback transaction", zap.Log.Error(err))
		}
	}()

	qtx := s.sqlc.WithTx(tx)

	// add balance to users
	for _, order := range orders {
		// get accrual before update
		orderMementoBefore, ok := order.GetMemento("beforeUpdate")
		if !ok {
			return fmt.Errorf("cannot get accrual before update")
		}

		addBalance := order.Accrual - orderMementoBefore.GetAccrual()
		if addBalance != 0 {
			err = qtx.UpdateBalance(ctx, sqlc.UpdateBalanceParams{
				Balance: addBalance,
				ID:      order.UserID,
			})
			if err != nil {
				return fmt.Errorf("could not update balance with userID %d: %w", order.UserID, err)
			}
		}
	}

	// update orders
	rows := make([]sqlc.UpdateOrdersParams, len(orders))
	for i, order := range orders {
		rows[i] = sqlc.UpdateOrdersParams{
			Status:  sqlc.NullOrderStatus{OrderStatus: sqlc.OrderStatus(order.Status), Valid: true},
			Accrual: order.Accrual,
			ID:      order.ID,
		}
	}

	result := qtx.UpdateOrders(ctx, rows)
	if err = result.Close(); err != nil {
		return fmt.Errorf("could not update orders: %w", err)
	}

	return tx.Commit(ctx)
}
