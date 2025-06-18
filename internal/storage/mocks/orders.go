package mocks

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/services/models"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
)

func (m *MockRepository) AddOrder(_ string, _ int) (orders.Order, error) {
	var order orders.Order
	return order, nil
}

func (m *MockRepository) ListUserOrders(_ context.Context, _ models.User) ([]orders.Order, error) {
	var result []orders.Order
	return result, nil
}
