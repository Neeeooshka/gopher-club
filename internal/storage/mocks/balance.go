package mocks

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/services/balance"
	"github.com/Neeeooshka/gopher-club/internal/services/models"
)

func (m *MockRepository) GetWithdrawals(_ context.Context, _ models.User) ([]balance.Withdraw, error) {
	var withdrawals []balance.Withdraw
	return withdrawals, nil
}

func (m *MockRepository) WithdrawBalance(_ context.Context, _ balance.Withdraw) error {
	return nil
}

func (m *MockRepository) GetWithdrawn(_ context.Context, _ models.User) (float64, error) {
	var withdrawn float64
	return withdrawn, nil
}
