package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

func GetOptions() Options {
	opt := NewOptions()
	cfg := NewConfig()

	flag.Var(&opt.ServerAddress, "a", "Server address - host:port")
	flag.Var(&opt.AccrualAddress, "r", "Accrual system address - protocol://host:port")
	flag.Var(&opt.DB, "d", "postgres connection string")

	flag.Parse()
	err := env.Parse(&cfg)
	if err != nil {
		panic(fmt.Errorf("error parsing config: %s", err))
	}

	if cfg.ServerAddress != "" {
		err = opt.ServerAddress.Set(cfg.ServerAddress)
		if err != nil {
			panic(fmt.Errorf("error setting ServerAddress: %s", err))
		}
	}

	if cfg.AccrualAddress != "" {
		err = opt.AccrualAddress.Set(cfg.AccrualAddress)
		if err != nil {
			panic(fmt.Errorf("error setting AccrualAddress: %s", err))
		}
	}

	if cfg.DB != "" {
		err = opt.DB.Set(cfg.DB)
		if err != nil {
			panic(fmt.Errorf("error setting DB: %s", err))
		}
	}

	return opt
}
