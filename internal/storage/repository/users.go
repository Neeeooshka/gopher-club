package repository

import (
	"context"
	"fmt"

	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"github.com/Neeeooshka/gopher-club/internal/storage/repository/sqlc"
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

	return user, nil
}

func (s *Postgres) AuthUser(user *models.User) error {

	authUserParam := sqlc.AuthUserParams{
		Login: user.Login,
		Crypt: user.Password,
	}

	u, err := s.sqlc.AuthUser(context.Background(), authUserParam)
	if err != nil {
		return err
	}

	user.ID = u.ID
	user.Password = u.Password
	user.Balance = u.Balance

	return nil
}

func (s *Postgres) AddUser(ctx context.Context, user models.User) error {

	u := sqlc.AddUserParams{
		Login: user.Login,
		Crypt: user.Password,
	}

	result, err := s.sqlc.AddUser(ctx, u)
	if err != nil {
		return fmt.Errorf("could not add user: %w", err)
	}

	if !result.IsNew {
		return storage.NewConflictUserError(result.ID, u.Login)
	}

	return nil
}
