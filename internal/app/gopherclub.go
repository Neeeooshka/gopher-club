package app

import (
	"context"
	"github.com/Neeeooshka/gopher-club.git/internal/config"
	"github.com/Neeeooshka/gopher-club.git/internal/storage"
	"net/http"
)

type gopherClubApp struct {
	Options config.Options
	storage storage.Membership
	context ctx
}

type ctx struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewGopherClubAppInstance(opt config.Options, s storage.Membership) *gopherClubApp {

	c, cancel := context.WithCancel(context.Background())

	instance := &gopherClubApp{
		Options: opt,
		storage: s,
		context: ctx{ctx: c, cancel: cancel},
	}

	return instance
}

func (a *gopherClubApp) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *gopherClubApp) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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
