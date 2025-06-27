package app

import (
	"github.com/Neeeooshka/gopher-club/internal/compressor"
	"github.com/Neeeooshka/gopher-club/internal/logger"
)

func (a *GopherClubApp) InitializeRoutes(log logger.Logger, comp compressor.Compressor) {
	if a.UserService.Inited {
		a.Router.Post("/api/user/register", logger.IncludeLogger(compressor.IncludeCompressor(a.UserService.RegisterUserHandler, comp), log))
		a.Router.Post("/api/user/login", logger.IncludeLogger(compressor.IncludeCompressor(a.UserService.LoginUserHandler, comp), log))
	} else {
		a.Router.Post("/api/user/register", a.ServiceUnavailableHandler)
		a.Router.Post("/api/user/login", a.ServiceUnavailableHandler)
	}

	if a.OrdersService.Inited {
		a.Router.Post("/api/user/orders", logger.IncludeLogger(compressor.IncludeCompressor(a.OrdersService.AddUserOrderHandler, comp), log))
		a.Router.Get("/api/user/orders", logger.IncludeLogger(a.OrdersService.GetUserOrdersHandler, log))
	} else {
		a.Router.Post("/api/user/orders", a.ServiceUnavailableHandler)
		a.Router.Get("/api/user/orders", a.ServiceUnavailableHandler)
	}

	if a.BalanceService.Inited {
		a.Router.Get("/api/user/balance", logger.IncludeLogger(a.BalanceService.GetUserBalanceHandler, log))
		a.Router.Post("/api/user/balance/withdraw", logger.IncludeLogger(compressor.IncludeCompressor(a.BalanceService.WithdrawBalanceHandler, comp), log))
		a.Router.Get("/api/user/withdrawals", logger.IncludeLogger(a.BalanceService.GetUserWithdrawalsHandler, log))
	} else {
		a.Router.Get("/api/user/balance", a.ServiceUnavailableHandler)
		a.Router.Post("/api/user/balance/withdraw", a.ServiceUnavailableHandler)
		a.Router.Get("/api/user/withdrawals", a.ServiceUnavailableHandler)
	}
}
