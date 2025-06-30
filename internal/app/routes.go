package app

import (
	"github.com/Neeeooshka/gopher-club/pkg/compressor"
	"github.com/Neeeooshka/gopher-club/pkg/logger"
	"github.com/go-chi/chi/v5"
)

func (a *GopherClubApp) InitializeRoutes(log logger.Logger, comp compressor.Compressor) {

	l := logger.NewLogger(log)
	c := compressor.NewCompressor(comp)

	a.Router.Use(l.Middleware, a.UserService.AuthMiddleware)

	a.Router.Group(func(r chi.Router) {

		r.Use(c.Middleware)

		r.Post("/api/user/register", a.UserService.RegisterUserHandler)
		r.Post("/api/user/login", a.UserService.LoginUserHandler)
		r.Post("/api/user/orders", a.OrdersService.AddUserOrderHandler)
		r.Post("/api/user/balance/withdraw", a.BalanceService.WithdrawBalanceHandler)
	})

	a.Router.Get("/api/user/orders", a.OrdersService.GetUserOrdersHandler)
	a.Router.Get("/api/user/balance", a.BalanceService.GetUserBalanceHandler)
	a.Router.Get("/api/user/withdrawals", a.BalanceService.GetUserWithdrawalsHandler)
}
