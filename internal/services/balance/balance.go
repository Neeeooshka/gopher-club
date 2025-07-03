package balance

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/models"
)

type BalanceRepository interface {
	GetWithdrawals(context.Context, models.User) ([]models.Withdraw, error)
	GetWithdrawn(context.Context, models.User) (float32, error)
	WithdrawBalance(context.Context, models.Withdraw) error
}

type BalanceService struct {
	errors  []error
	init    bool
	storage BalanceRepository
}

func NewBalanceService(repo BalanceRepository) BalanceService {

	var bs BalanceService

	bs.storage = repo
	bs.init = true

	return bs
}

func (b *BalanceService) HealthCheck() ([]error, bool) {
	return b.errors, b.init
}

func (b *BalanceService) GetName() string {
	return "BalanceService"
}
