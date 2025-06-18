package main

import (
	"errors"
	"flag"
	"github.com/Neeeooshka/gopher-club/internal/app"
	"github.com/Neeeooshka/gopher-club/internal/compressor"
	"github.com/Neeeooshka/gopher-club/internal/compressor/gzip"
	"github.com/Neeeooshka/gopher-club/internal/config"
	"github.com/Neeeooshka/gopher-club/internal/logger"
	"github.com/Neeeooshka/gopher-club/internal/logger/zap"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"github.com/Neeeooshka/gopher-club/internal/storage/postgres"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {

	var err error
	var store storage.Storage

	opt := getOptions()

	if opt.DB.String() != "" {
		store, err = postgres.NewPostgresStorage(opt.DB.String())
		if err != nil {
			panic(err)
		}
	} else {
		panic(errors.New("DB connection is not set"))
	}

	defer store.Close()

	app := app.NewGopherClubAppInstance(opt, store)

	zapLoger, err := zap.NewZapLogger("info")
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	if app.UserService.Inited {
		router.Post("/api/user/register", logger.IncludeLogger(compressor.IncludeCompressor(app.UserService.RegisterUserHandler, gzip.NewGzipCompressor()), zapLoger))
		router.Post("/api/user/login", logger.IncludeLogger(compressor.IncludeCompressor(app.UserService.LoginUserHandler, gzip.NewGzipCompressor()), zapLoger))
	} else {
		router.Post("/api/user/register", app.ServiceUnavialableHandler)
		router.Post("/api/user/login", app.ServiceUnavialableHandler)
	}

	router.Get("/api/user/balance", logger.IncludeLogger(app.GetUserBalanceHandler, zapLoger))
	router.Post("/api/user/orders", logger.IncludeLogger(compressor.IncludeCompressor(app.AddUserOrderHandler, gzip.NewGzipCompressor()), zapLoger))
	router.Get("/api/user/orders", logger.IncludeLogger(app.GetUserOrdersHandler, zapLoger))
	router.Post("/api/user/balance/withdraw", logger.IncludeLogger(compressor.IncludeCompressor(app.WithdrawUserBalanceHandler, gzip.NewGzipCompressor()), zapLoger))
	router.Get("/api/user/withdrawals", logger.IncludeLogger(app.GetUserWithdrawals, zapLoger))

	http.ListenAndServe(app.Options.GetServer(), router)
}

// init options
func getOptions() config.Options {
	opt := config.NewOptions()
	cfg := config.NewConfig()

	flag.Var(&opt.ServerAddress, "a", "Server address - host:port")
	flag.Var(&opt.AccrualAddress, "r", "Accrual system address - protocol://host:port")
	flag.Var(&opt.DB, "d", "postrgres connection string")

	flag.Parse()
	env.Parse(&cfg)

	if cfg.ServerAddress != "" {
		opt.ServerAddress.Set(cfg.ServerAddress)
	}

	if cfg.AccrualAddress != "" {
		opt.AccrualAddress.Set(cfg.AccrualAddress)
	}

	if cfg.DB != "" {
		opt.DB.Set(cfg.DB)
	}

	return opt
}
