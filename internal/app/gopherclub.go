package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Neeeooshka/gopher-club.git/internal/auth"
	"github.com/Neeeooshka/gopher-club.git/internal/config"
	"github.com/Neeeooshka/gopher-club.git/internal/storage"
	"github.com/Neeeooshka/gopher-club.git/internal/storage/postgres"
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

	var user struct {
		login string `json:"login"`
		pass  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	password, err := auth.CreatePassword(user.pass)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = a.storage.AddUser(auth.User{Login: user.login, Password: *password})
	var ce *postgres.ConflictUserError
	if err != nil {
		if errors.As(err, &ce) {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, ce.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
