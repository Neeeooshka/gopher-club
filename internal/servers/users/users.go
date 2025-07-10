package users

import (
	"context"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/dto"
	"net/http"
	"strings"

	"github.com/Neeeooshka/gopher-club/internal/models"
)

type UserRepository interface {
	AddUser(context.Context, models.User) error
	GetUserByLogin(string) (models.User, error)
	AuthUser(*models.User) error
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

func (u *UserServer) Authorize(auth dto.AuthData) (string, error) {

	user := models.User{
		Login:    auth.Login,
		Password: auth.Password,
	}

	if err := u.storage.AuthUser(&user); err != nil {
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
