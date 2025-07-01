package orders

import (
	"context"
	"errors"
	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"github.com/Neeeooshka/gopher-club/pkg/httputil"
	"github.com/Neeeooshka/gopher-club/pkg/logger/zap"
	"io"
	"net/http"
	"time"
)

func (o *OrdersService) AddUserOrderHandler(w http.ResponseWriter, r *http.Request) {

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
	if !CheckLuhn(orderNumber) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	order, err := o.storage.AddOrder(r.Context(), orderNumber, user.ID)
	var cue *storage.ConflictOrderError
	var coue *storage.ConflictOrderUserError
	if err != nil {
		if errors.As(err, &cue) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.As(err, &coue) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	o.updateService.AddWaitingOrder(order)

	w.WriteHeader(http.StatusAccepted)
}

func (o *OrdersService) GetUserOrdersHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := r.Context().Value(models.UserContextKey).(models.User)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	orders, err := o.storage.ListUserOrders(ctx, user)
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
