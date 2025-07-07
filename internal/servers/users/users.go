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

type UserServer struct {
	init    bool
	storage UserRepository
}

func NewUserServer(repo UserRepository) UserServer {
	return UserServer{storage: repo, init: true}
}

func (u *UserServer) AuthMiddleware(next http.Handler) http.Handler {
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

func (u *UserServer) Authenticate(jwtToken string) (models.User, error) {

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

func (u *UserServer) Authorize(cr credentials) (string, error) {

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

func (u *UserServer) HealthCheck() bool {
	return u.init
}

func (u *UserServer) GetName() string {
	return "UserServer"
}
