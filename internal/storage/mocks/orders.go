package mocks

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/models"
)

func (m *MockRepository) AddOrder(_ string, _ int) (models.Order, error) {
	var order models.Order
	return order, nil
}

func (m *MockRepository) ListUserOrders(_ context.Context, _ models.User) ([]models.Order, error) {
	var result []models.Order
	return result, nil
}
