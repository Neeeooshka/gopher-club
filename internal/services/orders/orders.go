package orders

import (
	"context"
	"errors"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/config"
	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"github.com/Neeeooshka/gopher-club/pkg/httputil"
	"github.com/Neeeooshka/gopher-club/pkg/logger/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	StatusNew        = "NEW"
	StatusProcessing = "PROCESSING"
	StatusInvalid    = "INVALID"
	StatusProcessed  = "PROCESSED"
)

type OrdersRepository interface {
	AddOrder(string, int) (models.Order, error)
	ListUserOrders(context.Context, models.User) ([]models.Order, error)
}

type OrdersService struct {
	Errors        []error
	Inited        bool
	storage       OrdersRepository
	updateService *OrdersUpdateService
}

func NewOrdersService(or interface{}, opt config.Options) OrdersService {

	var os OrdersService

	ordersRepo, ok := or.(OrdersRepository)

	if !ok {
		os.Errors = append(os.Errors, fmt.Errorf("2th argument expected OrdersRepository, got %T", or))
	}

	ous, err := NewOrdersUpdateService(or, opt)

	if err != nil {
		os.Errors = append(os.Errors, errors.New("cannot initialize OrdersUpdateService"))
	}

	if len(os.Errors) > 0 {
		return os
	}

	os.storage = ordersRepo
	os.updateService = ous
	os.Inited = true

	return os
}

func (o *OrdersService) AddUserOrderHandler(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user")
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	u := user.(models.User)

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

	order, err := o.storage.AddOrder(orderNumber, u.ID)
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

	user := r.Context().Value("user")
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	u := user.(models.User)

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	orders, err := o.storage.ListUserOrders(ctx, u)
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

func CheckLuhn(orderNumber string) bool {

	var sum int

	parity := len(orderNumber) % 2
	for i := 0; i < len(orderNumber); i++ {

		digit, err := strconv.Atoi(string(orderNumber[i]))
		if err != nil {
			return false
		}

		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}
