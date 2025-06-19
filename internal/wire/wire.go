//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/Neeeooshka/gopher-club/internal/app"
	"github.com/google/wire"
)

func InitializeApp() (*app.GopherClubApp, func(), error) {
	wire.Build(
		ProvideConfig,
		ProvidePostgresStorage,
		ProvideZapLogger,
		ProvideGzipCompressor,
		ProvideApp,
	)
	return nil, nil, nil
}
