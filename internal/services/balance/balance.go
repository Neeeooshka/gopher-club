package balance

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
	"github.com/Neeeooshka/gopher-club/pkg/httputil"
	"github.com/Neeeooshka/gopher-club/pkg/logger/zap"
	"net/http"
	"time"
)

type BalanceRepository interface {
	WithdrawBalance(context.Context, models.Withdraw) error
	GetWithdrawals(context.Context, models.User) ([]models.Withdraw, error)
	GetWithdrawn(context.Context, models.User) (float32, error)
}

type BalanceService struct {
	Inited  bool
	storage BalanceRepository
	Errors  []error
}

func NewBalanceService(or interface{}) BalanceService {

	var bs BalanceService

	balanceRepo, ok := or.(BalanceRepository)

	if !ok {
		bs.Errors = append(bs.Errors, fmt.Errorf("2th argument expected BalanceRepository, got %T", or))
	}

	if len(bs.Errors) > 0 {
		return bs
	}

	bs.storage = balanceRepo
	bs.Inited = true

	return bs
}

func (b *BalanceService) WithdrawBalanceHandler(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user")
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	u := user.(models.User)

	var withdraw models.Withdraw

	if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !orders.CheckLuhn(withdraw.OrderNum) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if u.Balance < withdraw.Sum {
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	withdraw.UserID = u.ID

	go func() {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err := b.storage.WithdrawBalance(ctx, withdraw)

		if err != nil {
			zap.Log.Debug(fmt.Sprintf("cannot withdraw balance for user %d", u.ID), zap.Log.Error(err))
		}
	}()

	w.WriteHeader(http.StatusOK)
}

func (b *BalanceService) GetUserBalanceHandler(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user")
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	u := user.(models.User)

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	withdrawn, err := b.storage.GetWithdrawn(ctx, u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balance := struct {
		Balance  float32 `json:"current"`
		Withdraw float32 `json:"withdrawn"`
	}{
		Balance:  u.Balance,
		Withdraw: withdrawn,
	}

	httputil.WriteJSON(w, balance)
}

func (b *BalanceService) GetUserWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user")
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	u := user.(models.User)

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	withdrawals, err := b.storage.GetWithdrawals(ctx, u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	httputil.WriteJSON(w, withdrawals)
}
