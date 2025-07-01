package balance

import (
	"context"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/models"
)

type BalanceRepository interface {
	GetWithdrawals(context.Context, models.User) ([]models.Withdraw, error)
	GetWithdrawn(context.Context, models.User) (float32, float32, error)
	WithdrawBalance(context.Context, models.Withdraw) error
}

type BalanceService struct {
	errors  []error
	init    bool
	storage BalanceRepository
}

func NewBalanceService(or interface{}) BalanceService {

	var bs BalanceService

	balanceRepo, ok := or.(BalanceRepository)

	if !ok {
		bs.errors = append(bs.errors, fmt.Errorf("2th argument expected BalanceRepository, got %T", or))
	}

	if len(bs.errors) > 0 {
		return bs
	}

	bs.storage = balanceRepo
	bs.init = true

	return bs
}

func (b *BalanceService) HealthCheck() ([]error, bool) {
	return b.errors, b.init
}

func (b *BalanceService) GetName() string {
	return "BalanceService"
}
