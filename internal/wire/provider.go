package wire

import (
	"github.com/Neeeooshka/gopher-club/internal/app"
	"github.com/Neeeooshka/gopher-club/internal/compressor"
	"github.com/Neeeooshka/gopher-club/internal/compressor/gzip"
	"github.com/Neeeooshka/gopher-club/internal/config"
	"github.com/Neeeooshka/gopher-club/internal/logger"
	"github.com/Neeeooshka/gopher-club/internal/logger/zap"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"github.com/Neeeooshka/gopher-club/internal/storage/mocks"
	"github.com/Neeeooshka/gopher-club/internal/storage/postgres"
)

func ProvideConfig() config.Options {
	return config.GetOptions()
}

func ProvidePostgresStorage(cfg config.Options) (storage.Storage, error) {
	return postgres.NewPostgresStorage(cfg.DB.String())
}

func ProvideMockStorage() (storage.Storage, error) {
	return &mocks.MockRepository{}, nil
}

func ProvideZapLogger() (logger.Logger, error) {
	return zap.NewZapLogger("info")
}

func ProvideGzipCompressor() compressor.Compressor {
	return gzip.NewGzipCompressor()
}

func ProvideApp(
	cfg config.Options,
	store storage.Storage,
	log logger.Logger,
	comp compressor.Compressor,
) *app.GopherClubApp {
	appInstance := app.NewGopherClubAppInstance(cfg, store)
	appInstance.InitializeRoutes(log, comp)

	return appInstance
}
