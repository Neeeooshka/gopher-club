package users

import (
	"context"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/models"
)

type UserRepository interface {
	AddUser(context.Context, models.User, string) error
	GetUserByLogin(string) (models.User, error)
}

type UserService struct {
	Errors  []error
	Inited  bool
	storage UserRepository
}

func NewUserService(ur interface{}) UserService {

	var us UserService

	userRepo, ok := ur.(UserRepository)

	if !ok {
		us.Errors = append(us.Errors, fmt.Errorf("2th argument expected UserRepository, got %T", ur))
	}

	if len(us.Errors) > 0 {
		return us
	}

	us.storage = userRepo
	us.Inited = true

	return us
}

func (u *UserService) Authenticate(jwtToken string) (models.User, error) {

	var user models.User

	login, err := VerifyJWTToken(jwtToken)
	if err != nil {
		return user, err
	}

	user, err = u.storage.GetUserByLogin(login)
	if err != nil {
		return user, fmt.Errorf("error authentication: %w", err)
	}

	return user, nil
}

func (u *UserService) Authorize(cr credentials) (string, error) {

	user, err := u.storage.GetUserByLogin(cr.Login)
	if err != nil {
		return "", fmt.Errorf("error authorization: %w", err)
	}

	err = cr.verifyPassword(user)
	if err != nil {
		return "", fmt.Errorf("error authorization: %w", err)
	}

	return CreateJWTToken(user.Login)
}
