package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"github.com/Neeeooshka/gopher-club/internal/storage/postgres/sqlc"
	"github.com/Neeeooshka/gopher-club/pkg/logger/zap"
	"github.com/jackc/pgx/v5"
)

func (s *Postgres) GetUserByLogin(login string) (models.User, error) {

	var user models.User

	u, err := s.sqlc.GetUserByLogin(context.Background(), login)
	if err != nil {
		return user, err
	}

	user.ID = u.ID
	user.Login = u.Login
	user.Password = u.Password
	user.Balance = u.Balance
	user.Credentials = u.Credentials

	return user, nil
}

func (s *Postgres) AddUser(ctx context.Context, user models.User, salt string) error {

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

	u := sqlc.AddUserParams{
		Login:    user.Login,
		Password: user.Password,
	}

	result, err := qtx.AddUser(ctx, u)
	if err != nil {
		return fmt.Errorf("could not add user: %w", err)
	}

	if !result.IsNew {
		return storage.NewConflictUserError(result.ID, u.Login)
	}

	c := sqlc.AddCredentialsParams{
		UserID: result.ID,
		PValue: salt,
	}

	err = qtx.AddCredentials(ctx, c)
	if err != nil {
		return fmt.Errorf("could not add credentials: %w", err)
	}

	return tx.Commit(ctx)
}
