package app

import (
	"github.com/Neeeooshka/gopher-club/internal/config"
	"github.com/Neeeooshka/gopher-club/internal/services/balance"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"github.com/go-chi/chi/v5"
)

type GopherClubApp struct {
	Options config.Options
	storage storage.Storage
	Router  *chi.Mux

	BalanceService balance.BalanceService
	OrdersService  orders.OrdersService
	UserService    users.UserService
}

func NewGopherClubAppInstance(opt config.Options, s storage.Storage) *GopherClubApp {

	instance := &GopherClubApp{
		Options: opt,
		storage: s,
		Router:  chi.NewRouter(),

		BalanceService: balance.NewBalanceService(s),
		OrdersService:  orders.NewOrdersService(s, opt),
		UserService:    users.NewUserService(s),
	}

	return instance
}
