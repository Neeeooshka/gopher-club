package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/config"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"github.com/Neeeooshka/gopher-club/internal/users"
	"net/http"
)

type gopherClubApp struct {
	Options config.Options
	storage storage.Storage
	context ctx
}

type ctx struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type user struct {
	login string `json:"login"`
	pass  string `json:"password"`
}

func (u *user) validate() bool {
	return u.login != "" && u.pass != ""
}

func NewGopherClubAppInstance(opt config.Options, s storage.Storage) *gopherClubApp {

	c, cancel := context.WithCancel(context.Background())

	instance := &gopherClubApp{
		Options: opt,
		storage: s,
		context: ctx{ctx: c, cancel: cancel},
	}

	return instance
}

func (a *gopherClubApp) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {

	var u user

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !u.validate() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	password, err := users.CreatePassword(u.pass)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = a.storage.AddUser(users.User{Login: u.login, Password: password})
	var ce *users.ConflictUserError
	if err != nil {
		if errors.As(err, &ce) {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *gopherClubApp) LoginUserHandler(w http.ResponseWriter, r *http.Request) {

	var u user

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !u.validate() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := a.Authenticate(u)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

func (a *gopherClubApp) Authenticate(u user) error {

	userData, err := a.storage.GetUserByLogin(u.login)

	if err != nil {
		return fmt.Errorf("Ошибка аутентификации: %w", err)
	}

	key, err := a.storage.GetUserKey(userData.ID)

	if err != nil {
		return fmt.Errorf("Ошибка аутентификации: %w", err)
	}

	userData.Password, err = users.NewPassword(userData.Hash, key)

	if err != nil {
		return fmt.Errorf("Ошибка аутентификации: %w", err)
	}

	if !userData.Password.Verify(u.pass, userData.Token) {
		return fmt.Errorf("Ошибка аутентификации: %w", errors.New("Неверное имя пользователя или пароль"))
	}

	return nil
}
