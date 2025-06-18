package mocks

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
)

func (m *MockRepository) ListWaitingOrders(_ context.Context) ([]orders.Order, error) {
	var result []orders.Order
	return result, nil
}

func (m *MockRepository) UpdateOrders(_ context.Context, _ []orders.Order) error {
	return nil
}
