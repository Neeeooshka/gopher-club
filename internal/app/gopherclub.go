package app

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/config"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"net/http"
)

type Service interface {
	Init(Ctx, interface{}) error
}

type gopherClubApp struct {
	Options     config.Options
	storage     storage.Storage
	ctx         Ctx
	UserService users.UserService
}

type Ctx struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewGopherClubAppInstance(opt config.Options, s storage.Storage) *gopherClubApp {

	c, cancel := context.WithCancel(context.Background())
	ctx := &Ctx{Ctx: c, Cancel: cancel}

	instance := &gopherClubApp{
		UserService: users.NewUserService(c, s),
		Options:     opt,
		ctx:         *ctx,
		storage:     s,
	}

	return instance
}

func (a *gopherClubApp) ServiceUnavialableHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
}

func (a *gopherClubApp) AddUserOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *gopherClubApp) GetUserOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *gopherClubApp) GetUserBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *gopherClubApp) WithdrawUserBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *gopherClubApp) GetUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
