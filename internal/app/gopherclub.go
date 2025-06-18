package app

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/config"
	"github.com/Neeeooshka/gopher-club/internal/services/balance"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"net/http"
)

type GopherClubApp struct {
	BalanceService balance.BalanceService
	Options        config.Options
	OrdersService  orders.OrdersService
	storage        storage.Storage
	ctx            Ctx
	UserService    users.UserService
}

type Ctx struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewGopherClubAppInstance(opt config.Options, s storage.Storage) *GopherClubApp {

	c, cancel := context.WithCancel(context.Background())
	ctx := &Ctx{Ctx: c, Cancel: cancel}
	us := users.NewUserService(c, s)
	os := orders.NewOrdersService(c, s, &us)

	instance := &GopherClubApp{
		BalanceService: balance.NewBalanceService(c, s, &us, &os),
		Options:        opt,
		OrdersService:  os,
		UserService:    us,
		ctx:            *ctx,
		storage:        s,
	}

	return instance
}

func (a *GopherClubApp) ServiceUnavailableHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
}
