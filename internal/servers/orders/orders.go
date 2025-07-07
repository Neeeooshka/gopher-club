package orders

import (
	"context"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/config"
	"github.com/Neeeooshka/gopher-club/internal/models"
	"time"
)

type OrdersRepository interface {
	AddOrder(context.Context, string, int) (models.Order, error)
	ListUserOrders(context.Context, models.User) ([]models.Order, error)
	ListWaitingOrders(context.Context) ([]models.Order, error)
}

type OrdersServer struct {
	init          bool
	storage       OrdersRepository
	updateService *OrdersUpdateService
}

func NewOrdersServer(repo OrdersRepository, opt config.Options) (OrdersServer, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	orders, err := repo.ListWaitingOrders(ctx)
	if err != nil {
		return OrdersServer{}, fmt.Errorf("unable to pull order list: %w", err)
	}

	updateService, err := NewOrdersUpdateService(repo.(OrdersUpdateRepository), opt, orders)

	if err != nil {
		return OrdersServer{}, fmt.Errorf("cannot initialize OrdersUpdateService: %v", err)
	}

	return OrdersServer{storage: repo, updateService: updateService, init: true}, nil
}

func (o *OrdersServer) HealthCheck() bool {
	return o.init
}

func (o *OrdersServer) GetName() string {
	return "OrdersServer"
}
