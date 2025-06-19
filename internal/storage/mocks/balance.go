package mocks

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/models"
)

func (m *MockRepository) GetWithdrawals(_ context.Context, _ models.User) ([]models.Withdraw, error) {
	var withdrawals []models.Withdraw
	return withdrawals, nil
}

func (m *MockRepository) WithdrawBalance(_ context.Context, _ models.Withdraw) error {
	return nil
}

func (m *MockRepository) GetWithdrawn(_ context.Context, _ models.User) (float64, error) {
	var withdrawn float64
	return withdrawn, nil
}
