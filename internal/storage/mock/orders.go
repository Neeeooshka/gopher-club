package mock

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
)

func (m *MockRepository) AddOrder(_ string, _ int) (orders.Order, error) {
	var order orders.Order
	return order, nil
}

func (m *MockRepository) ListUserOrders(_ context.Context, _ users.User) ([]orders.Order, error) {
	var result []orders.Order
	return result, nil
}
