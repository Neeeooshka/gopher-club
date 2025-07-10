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

		r.Use(a.UserServer.AuthMiddleware)

		// OrdersServer handlers
		r.Group(func(r chi.Router) {

			r.Use(a.HealthCheckMiddleware(&a.OrdersServer))

			r.Get("/api/user/orders", a.OrdersServer.GetUserOrdersHandler)
			r.Post("/api/user/orders", a.OrdersServer.AddUserOrderHandler)
		})

		// BalanceServer handlers
		r.Group(func(r chi.Router) {

			r.Use(a.HealthCheckMiddleware(&a.BalanceServer))

			r.Get("/api/user/balance", a.BalanceServer.GetUserBalanceHandler)
			r.Get("/api/user/withdrawals", a.BalanceServer.GetUserWithdrawalsHandler)
			r.Post("/api/user/balance/withdraw", a.BalanceServer.WithdrawBalanceHandler)
		})
	})

	a.Router.Post("/api/user/register", a.UserServer.RegisterUserHandler)
	a.Router.Post("/api/user/login", a.UserServer.LoginUserHandler)
}

func (a *GopherClubApp) getMiddlewares() []func(http.Handler) http.Handler {

	var middlewares []func(http.Handler) http.Handler

	// logger
	if a.logger != nil {
		middlewares = append(middlewares, a.logger.Middleware)
	}

	// HealthChecker UserServer for all requests
	middlewares = append(middlewares, a.HealthCheckMiddleware(&a.UserServer))

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
