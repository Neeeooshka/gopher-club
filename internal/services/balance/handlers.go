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

func (b *BalanceService) WithdrawBalanceHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(models.UserContextKey).(models.User)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	withdraw := models.Withdraw{UserID: user.ID}

	if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !orders.CheckLuhn(withdraw.OrderNum) {
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

		err := b.storage.WithdrawBalance(ctx, withdraw)

		if err != nil {
			zap.Log.Debug(fmt.Sprintf("cannot withdraw balance for user %d", user.ID), zap.Log.Error(err))
		}
	}()

	w.WriteHeader(http.StatusOK)
}

func (b *BalanceService) GetUserBalanceHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(models.UserContextKey).(models.User)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	withdrawn, err := b.storage.GetWithdrawn(ctx, user)
	if err != nil {
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

	httputil.WriteJSON(w, balance)
}

func (b *BalanceService) GetUserWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(models.UserContextKey).(models.User)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
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

	httputil.WriteJSON(w, withdrawals)
}
