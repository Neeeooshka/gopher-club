package balance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/logger/zap"
	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
	"net/http"
	"time"
)

type BalanceRepository interface {
	WithdrawBalance(context.Context, models.Withdraw) error
	GetWithdrawals(context.Context, models.User) ([]models.Withdraw, error)
	GetWithdrawn(context.Context, models.User) (float32, error)
}

type BalanceService struct {
	Inited        bool
	storage       BalanceRepository
	Errors        []error
	UserService   *users.UserService
	OrdersService *orders.OrdersService
}

func NewBalanceService(or interface{}, us *users.UserService, os *orders.OrdersService) BalanceService {

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

	var withdraw models.Withdraw

	if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !b.OrdersService.CheckLuhn(withdraw.OrderNum) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if user.Balance < withdraw.Sum {
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	go func() {

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
		defer cancel()

		err = b.storage.WithdrawBalance(ctx, withdraw)

		if err != nil {
			logger, _ := zap.NewZapLogger("debug")
			logger.Debug(fmt.Sprintf("cannot withdraw balance for user %d", user.ID), logger.Error(err))
		}
	}()
}

func (b *BalanceService) GetUserBalanceHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	user, err := b.UserService.Authenticate(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	withdrawn, err := b.storage.GetWithdrawn(ctx, user)
	if err != nil {
		logger, _ := zap.NewZapLogger("debug")
		logger.Debug("internal error", logger.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balance := struct {
		Balance  float32 `json:"current"`
		Withdraw float32 `json:"withdrawn"`
	}{
		Balance:  user.Balance,
		Withdraw: withdrawn,
	}

	if err := json.NewEncoder(w).Encode(balance); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (b *BalanceService) GetUserWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	user, err := b.UserService.Authenticate(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	withdrawals, err := b.storage.GetWithdrawals(ctx, user)
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
	}
}
