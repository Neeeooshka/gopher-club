package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *GopherClubApp) InitializeRoutes() {

	// middlewares for all handlers
	a.Router.Use(a.getMiddlewares()...)

	// auth only handlers
	a.Router.Group(func(r chi.Router) {

		r.Use(a.UserService.AuthMiddleware)

		// OrdersService handlers
		r.Group(func(r chi.Router) {

			r.Use(a.HealthCheckMiddleware(&a.OrdersService))

			r.Get("/api/user/orders", a.OrdersService.GetUserOrdersHandler)
			r.Post("/api/user/orders", a.OrdersService.AddUserOrderHandler)
		})

		// BalanceService handlers
		r.Group(func(r chi.Router) {

			r.Use(a.HealthCheckMiddleware(&a.BalanceService))

			r.Get("/api/user/balance", a.BalanceService.GetUserBalanceHandler)
			r.Get("/api/user/withdrawals", a.BalanceService.GetUserWithdrawalsHandler)
			r.Post("/api/user/balance/withdraw", a.BalanceService.WithdrawBalanceHandler)
		})
	})

	a.Router.Post("/api/user/register", a.UserService.RegisterUserHandler)
	a.Router.Post("/api/user/login", a.UserService.LoginUserHandler)
}

func (a *GopherClubApp) getMiddlewares() []func(http.Handler) http.Handler {

	var middlewares []func(http.Handler) http.Handler

	// logger
	if a.logger != nil {
		middlewares = append(middlewares, a.logger.Middleware)
	}

	// HealthChecker UserService for all requests
	middlewares = append(middlewares, a.HealthCheckMiddleware(&a.UserService))

	if a.compressor != nil {
		// compressor reader
		middlewares = append(middlewares, a.compressor.Middleware)
		// compressor writer
		middlewares = append(middlewares, middleware.Compress(5, a.compressor.GetEncoding()))
	}

	// set timeout for all requests
	middlewares = append(middlewares, a.TimeoutMiddleware)

	return middlewares
}
