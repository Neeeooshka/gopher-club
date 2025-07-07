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

type BalanceServer struct {
	init    bool
	storage BalanceRepository
}

func NewBalanceServer(repo BalanceRepository) BalanceServer {
	return BalanceServer{storage: repo, init: true}
}

func (b *BalanceServer) HealthCheck() bool {
	return b.init
}

func (b *BalanceServer) GetName() string {
	return "BalanceServer"
}
