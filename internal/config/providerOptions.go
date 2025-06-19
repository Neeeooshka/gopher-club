package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

func GetOptions() Options {
	opt := NewOptions()
	cfg := NewConfig()

	flag.Var(&opt.ServerAddress, "a", "Server address - host:port")
	flag.Var(&opt.AccrualAddress, "r", "Accrual system address - protocol://host:port")
	flag.Var(&opt.DB, "d", "postgres connection string")

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
