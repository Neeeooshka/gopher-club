package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Neeeooshka/gopher-club/internal/config"
	"github.com/Neeeooshka/gopher-club/internal/services/balance"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"github.com/Neeeooshka/gopher-club/pkg/compressor"
	"github.com/Neeeooshka/gopher-club/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type HealthChecker interface {
	HealthCheck() ([]error, bool)
	GetName() string
}

type GopherClubApp struct {
	Options config.Options
	storage storage.Storage
	Router  *chi.Mux

	BalanceService balance.BalanceService
	OrdersService  orders.OrdersService
	UserService    users.UserService

	logger     *logger.LoggerWrap
	compressor *compressor.CompressorWrap
}

func NewGopherClubAppInstance(opt config.Options, s storage.Storage) *GopherClubApp {

	instance := &GopherClubApp{
		Options: opt,
		storage: s,
		Router:  chi.NewRouter(),

		BalanceService: balance.NewBalanceService(s.(balance.BalanceRepository)),
		OrdersService:  orders.NewOrdersService(s.(orders.OrdersRepository), opt),
		UserService:    users.NewUserService(s.(users.UserRepository)),
	}

	return instance
}

func (a *GopherClubApp) WithLogger(l logger.Logger) *GopherClubApp {
	a.logger = logger.NewLogger(l)
	return a
}

func (a *GopherClubApp) WithCompressor(c compressor.Compressor) *GopherClubApp {
	a.compressor = compressor.NewCompressor(c)
	return a
}

// HealthCheckMiddleware checks the service for init and errors
func (a *GopherClubApp) HealthCheckMiddleware(service HealthChecker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			if errs, init := service.HealthCheck(); !init || errs != nil {
				log.Printf("Serice %s unavailabale: %v", service.GetName(), errs)
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func (a *GopherClubApp) TimeoutMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
		defer cancel()

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
