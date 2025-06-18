package mock

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/services/balance"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
)

func (m *MockRepository) GetWithdrawals(_ context.Context, _ users.User) ([]balance.Withdraw, error) {
	var withdrawals []balance.Withdraw
	return withdrawals, nil
}

func (m *MockRepository) WithdrawBalance(_ context.Context, _ balance.Withdraw) error {
	return nil
}

func (m *MockRepository) GetWithdrawn(_ context.Context, _ users.User) (float64, error) {
	var withdrawn float64
	return withdrawn, nil
}
