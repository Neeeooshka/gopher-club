package postgres

import (
	"context"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/internal/storage/postgres/sqlc"
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

	defer tx.Rollback(ctx)

	rows := make([]sqlc.UpdateOrdersParams, len(orders))
	for i, order := range orders {
		rows[i] = sqlc.UpdateOrdersParams{
			Status:  order.Status,
			Accrual: order.Accrual,
			ID:      order.ID,
		}
	}

	qtx := s.sqlc.WithTx(tx)
	result := qtx.UpdateOrders(ctx, rows)
	if err = result.Close(); err != nil {
		return fmt.Errorf("could not update orders: %w", err)
	}

	for _, order := range orders {
		// get accrual before update
		orderMementoBefore := order.GetMemento("beforeUpdate")
		if orderMementoBefore != nil {
			return fmt.Errorf("cannot get accrual before update")
		}

		// add balance to user
		addBalance := order.Accrual - orderMementoBefore.GetAccrual()
		if addBalance != 0 {
			err = s.sqlc.UpdateBalance(ctx, sqlc.UpdateBalanceParams{
				Balance: addBalance,
				ID:      order.ID,
			})
			if err != nil {
				return fmt.Errorf("could not update balance with userID %d: %w", order.UserID, err)
			}
		}
	}

	return tx.Commit(ctx)
}
