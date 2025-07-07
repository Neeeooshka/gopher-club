package orders

import (
	"errors"
	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"github.com/Neeeooshka/gopher-club/pkg/httputil"
	"github.com/Neeeooshka/gopher-club/pkg/logger/zap"
	"github.com/Neeeooshka/gopher-club/pkg/utils"
	"io"
	"net/http"
)

func (o *OrdersServer) AddUserOrderHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(models.UserContextKey).(models.User)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			zap.Log.Debug("failed to close request body reader", zap.Log.Error(err))
		}
	}()

	orderNumber := string(body)
	if !utils.CheckLuhn(orderNumber) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	order, err := o.storage.AddOrder(r.Context(), orderNumber, user.ID)
	var conflictOrderError *storage.ConflictOrderError
	var conflictOrderUserError *storage.ConflictOrderUserError
	if err != nil {
		if errors.As(err, &conflictOrderError) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.As(err, &conflictOrderUserError) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	o.updateService.AddWaitingOrder(order)

	w.WriteHeader(http.StatusAccepted)
}

func (o *OrdersServer) GetUserOrdersHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(models.UserContextKey).(models.User)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	orders, err := o.storage.ListUserOrders(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	httputil.WriteJSON(w, orders)
}
