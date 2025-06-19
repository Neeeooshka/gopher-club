package mocks

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/models"
)

func (m *MockRepository) ListWaitingOrders(_ context.Context) ([]models.Order, error) {
	var result []models.Order
	return result, nil
}

func (m *MockRepository) UpdateOrders(_ context.Context, _ []models.Order) error {
	return nil
}
