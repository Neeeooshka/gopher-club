package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/logger/zap"
	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/internal/storage/postgres/sqlc"
)

func (s *Postgres) GetWithdrawals(ctx context.Context, user models.User) ([]models.Withdraw, error) {

	results, err := s.sqlc.GetWithdrawals(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting withdrawals: %w", err)
	}

	withdrawals := make([]models.Withdraw, len(results))
	for i, result := range results {
		withdrawals[i] = models.Withdraw{
			ID:           result.ID,
			UserID:       result.UserID,
			OrderNum:     result.Num,
			DateWithdraw: result.DateWithdraw,
			Sum:          result.Sum,
		}
	}

	return withdrawals, nil
}

func (s *Postgres) WithdrawBalance(ctx context.Context, w models.Withdraw) error {

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, sql.ErrTxDone) {
			logger, _ := zap.NewZapLogger("debug")
			logger.Debug("failed to rollback transaction", logger.Error(err))
		}
	}()

	qtx := s.sqlc.WithTx(tx)

	err = qtx.WithdrawBalance(ctx, sqlc.WithdrawBalanceParams{
		UserID: w.UserID,
		Num:    w.OrderNum,
		Sum:    w.Sum,
	})

	if err != nil {
		return fmt.Errorf("could not withdraw balance: %w", err)
	}

	err = qtx.UpdateBalance(ctx, sqlc.UpdateBalanceParams{Balance: w.Sum * -1, ID: w.UserID})
	if err != nil {
		return fmt.Errorf("could not update user balance: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Postgres) GetWithdrawn(ctx context.Context, user models.User) (float64, error) {
	return s.sqlc.GetWithdrawn(ctx, user.ID)
}
