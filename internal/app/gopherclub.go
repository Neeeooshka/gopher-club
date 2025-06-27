package app

import (
	"github.com/Neeeooshka/gopher-club/internal/config"
	"github.com/Neeeooshka/gopher-club/internal/services/balance"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
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

	us := users.NewUserService(s)
	os := orders.NewOrdersService(s, &us, opt)

	instance := &GopherClubApp{
		Options: opt,
		storage: s,
		Router:  chi.NewRouter(),

		BalanceService: balance.NewBalanceService(s, &us, &os),
		OrdersService:  os,
		UserService:    us,
	}

	return instance
}

func (a *GopherClubApp) ServiceUnavailableHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
}
