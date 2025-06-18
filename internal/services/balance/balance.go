package balance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/logger/zap"
	"github.com/Neeeooshka/gopher-club/internal/services/models"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
	"net/http"
	"time"
)

type BalanceRepository interface {
	WithdrawBalance(context.Context, Withdraw) error
	GetWithdrawals(context.Context, models.User) ([]Withdraw, error)
	GetWithdrawn(context.Context, models.User) (float64, error)
}

type BalanceService struct {
	Inited        bool
	storage       BalanceRepository
	ctx           context.Context
	Errors        []error
	UserService   *users.UserService
	OrdersService *orders.OrdersService
}

func NewBalanceService(ctx context.Context, or interface{}, us *users.UserService, os *orders.OrdersService) BalanceService {

	var bs BalanceService

	balanceRepo, ok := or.(BalanceRepository)

	if !ok {
		bs.Errors = append(bs.Errors, fmt.Errorf("2th argument expected BalanceRepository, got %T", or))
	}

	if !us.Inited {
		bs.Errors = append(bs.Errors, errors.New("UserService is unavailable"))
	}

	if !os.Inited {
		bs.Errors = append(bs.Errors, errors.New("OrdersService is unavailable"))
	}

	if len(bs.Errors) > 0 {
		return bs
	}

	bs.ctx = ctx
	bs.storage = balanceRepo
	bs.UserService = us
	bs.OrdersService = os
	bs.Inited = true

	return bs
}

func (b *BalanceService) WithdrawBalanceHandler(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	user, err := b.UserService.Authenticate(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var withdraw Withdraw

	if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !b.OrdersService.CheckLuhn(withdraw.OrderID) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if user.Balance < withdraw.Sum {
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	go func() {

		err = b.storage.WithdrawBalance(b.ctx, withdraw)

		if err != nil {
			logger, _ := zap.NewZapLogger("debug")
			logger.Debug(fmt.Sprintf("cannot withdraw balance for user %d", user.ID), logger.Error(err))
		}
	}()

	w.WriteHeader(http.StatusOK)
}

func (b *BalanceService) GetUserBalanceHandler(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	user, err := b.UserService.Authenticate(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	withdrawn, err := b.storage.GetWithdrawn(b.ctx, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balance := struct {
		Balance  float64 `json:"current"`
		Withdraw float64 `json:"withdrawn"`
	}{
		Balance:  user.Balance,
		Withdraw: withdrawn,
	}

	if err := json.NewEncoder(w).Encode(balance); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (b *BalanceService) GetUserWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	user, err := b.UserService.Authenticate(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	withdrawals, err := b.storage.GetWithdrawals(b.ctx, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(w).Encode(withdrawals); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type Withdraw struct {
	ID           int       `db:"id"`
	UserID       int       `db:"user_id"`
	OrderID      string    `db:"order_id" json:"order"`
	DateWithdraw time.Time `db:"date_withdraw" json:"processed_at"`
	Sum          float64   `db:"sum" json:"sum"`
}
