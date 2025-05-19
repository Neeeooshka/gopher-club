package main

import (
	"errors"
	"flag"
	"github.com/Neeeooshka/gopher-club.git/internal/app"
	"github.com/Neeeooshka/gopher-club.git/internal/compressor"
	"github.com/Neeeooshka/gopher-club.git/internal/compressor/gzip"
	"github.com/Neeeooshka/gopher-club.git/internal/config"
	"github.com/Neeeooshka/gopher-club.git/internal/logger"
	"github.com/Neeeooshka/gopher-club.git/internal/logger/zap"
	"github.com/Neeeooshka/gopher-club.git/internal/storage"
	"github.com/Neeeooshka/gopher-club.git/internal/storage/postgres"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {

	var err error
	var store storage.Membership

	opt := getOptions()

	if opt.DB.String() != "" {
		store, err = postgres.NewPostgresLinksStorage(opt.DB.String())
		if err != nil {
			panic(err)
		}
	} else {
		panic(errors.New("DB connection is not set"))
	}

	defer store.Close()

	appInstance := app.NewGopherClubAppInstance(opt, store)

	zapLoger, err := zap.NewZapLogger("info")
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()
	router.Post("/api/user/register", logger.IncludeLogger(compressor.IncludeCompressor(appInstance.RegisterUserHandler, gzip.NewGzipCompressor()), zapLoger))
	router.Post("/api/user/login", logger.IncludeLogger(compressor.IncludeCompressor(appInstance.LoginUserHandler, gzip.NewGzipCompressor()), zapLoger))
	router.Post("/api/user/orders", logger.IncludeLogger(compressor.IncludeCompressor(appInstance.AddUserOrderHandler, gzip.NewGzipCompressor()), zapLoger))
	router.Get("/api/user/orders", logger.IncludeLogger(appInstance.GetUserOrdersHandler, zapLoger))
	router.Get("/api/user/balance", logger.IncludeLogger(appInstance.GetUserBalanceHandler, zapLoger))
	router.Post("/api/user/balance/withdraw", logger.IncludeLogger(compressor.IncludeCompressor(appInstance.WithdrawUserBalanceHandler, gzip.NewGzipCompressor()), zapLoger))
	router.Get("/api/user/withdrawals", logger.IncludeLogger(appInstance.GetUserWithdrawals, zapLoger))

	http.ListenAndServe(appInstance.Options.GetServer(), router)
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
