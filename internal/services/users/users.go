package users

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Neeeooshka/gopher-club/internal/models"
)

type UserRepository interface {
	AddUser(context.Context, models.User, string) error
	GetUserByLogin(string) (models.User, error)
}

type UserService struct {
	errors  []error
	init    bool
	storage UserRepository
}

func NewUserService(repo UserRepository) UserService {

	var us UserService

	us.storage = repo
	us.init = true

	return us
}

func (u *UserService) AuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)

		user, err := u.Authenticate(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), models.UserContextKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
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

func (u *UserService) HealthCheck() ([]error, bool) {
	return u.errors, u.init
}

func (u *UserService) GetName() string {
	return "UserService"
}
