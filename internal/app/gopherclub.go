package app

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/config"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"net/http"
)

type GopherClubApp struct {
	Options       config.Options
	storage       storage.Storage
	ctx           Ctx
	UserService   users.UserService
	OrdersService orders.OrdersService
}

type Ctx struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewGopherClubAppInstance(opt config.Options, s storage.Storage) *GopherClubApp {

	c, cancel := context.WithCancel(context.Background())
	ctx := &Ctx{Ctx: c, Cancel: cancel}
	us := users.NewUserService(c, s)

	instance := &GopherClubApp{
		UserService:   us,
		OrdersService: orders.NewOrdersService(c, s, &us),
		Options:       opt,
		ctx:           *ctx,
		storage:       s,
	}

	return instance
}

func (a *GopherClubApp) ServiceUnavailableHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
}

func (a *GopherClubApp) AddUserOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *GopherClubApp) GetUserOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *GopherClubApp) GetUserBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *GopherClubApp) WithdrawUserBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a *GopherClubApp) GetUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
