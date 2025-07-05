package balance

import (
	"encoding/json"
	"github.com/Neeeooshka/gopher-club/internal/dto"
	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
	"github.com/Neeeooshka/gopher-club/pkg/httputil"
	"net/http"
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

	err := b.storage.WithdrawBalance(r.Context(), withdraw)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (b *BalanceService) GetUserBalanceHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(models.UserContextKey).(models.User)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	withdrawn, err := b.storage.GetWithdrawn(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	httputil.WriteJSON(w, dto.Balance{
		Balance:  user.Balance,
		Withdraw: withdrawn,
	})
}

func (b *BalanceService) GetUserWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(models.UserContextKey).(models.User)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	withdrawals, err := b.storage.GetWithdrawals(r.Context(), user)
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
