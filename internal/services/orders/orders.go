package orders

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
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
	AddOrder(string, int) error
	ListOrders(context.Context) ([]Order, error)
	ListUserOrders(context.Context, users.User) ([]Order, error)
	UpdateOrders(context.Context, []Order) error
}

type OrdersService struct {
	Errors         []error
	Inited         bool
	UserService    *users.UserService
	ctx            context.Context
	storage        OrdersRepository
	updateInterval time.Duration
}

func NewOrdersService(ctx context.Context, or interface{}, us *users.UserService) OrdersService {

	var os OrdersService

	ordersRepo, ok := or.(OrdersRepository)

	if !ok {
		os.Errors = append(os.Errors, fmt.Errorf("2th argument expected OrdersRepository, got %T", or))
	}

	if !us.Inited {
		os.Errors = append(os.Errors, errors.New("UserService is unavailable"))
	}

	if len(os.Errors) > 0 {
		return os
	}

	os.ctx = ctx
	os.storage = ordersRepo
	os.updateInterval = time.Minute * 5
	os.UserService = us
	os.Inited = true

	return os
}

func (o *OrdersService) AddUserOrderHandler(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	err := o.UserService.Authenticate(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	orderNumber := string(body)
	if !o.CheckLuhn(orderNumber) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = o.storage.AddOrder(orderNumber, o.UserService.User.ID)
	var cue *ConflictOrderError
	var coue *ConflictOrderUserError
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

	w.WriteHeader(http.StatusAccepted)
}

func (o *OrdersService) GetUserOrdersHandler(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	err := o.UserService.Authenticate(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orders, err := o.storage.ListUserOrders(o.ctx, o.UserService.User)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(w).Encode(orders); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (o *OrdersService) CheckLuhn(orderNumber string) bool {

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

func (o *OrdersService) updateOrders() {

	_ = time.NewTicker(o.updateInterval)

}

type Order struct {
	ID         int       `db:"id"`
	UserID     int       `db:"user_id"`
	Number     string    `db:"num" json:"number"`
	DateInsert time.Time `db:"date_insert" json:"uploaded_at"`
	Accrual    float64   `db:"accrual" json:"accrual,omitempty"`
	Status     string    `db:"status" json:"status"`
}

type ConflictOrderError struct {
	orderID string
}

func NewConflictOrderError(orderID string) *ConflictOrderError {
	return &ConflictOrderError{orderID}
}

func (e *ConflictOrderError) Error() string {
	return fmt.Sprintf("Order %s already exsists", e.orderID)
}

type ConflictOrderUserError struct {
	userID  int
	orderID string
}

func NewConflictOrderUserError(userID int, orderID string) *ConflictOrderUserError {
	return &ConflictOrderUserError{userID, orderID}
}

func (e *ConflictOrderUserError) Error() string {
	return fmt.Sprintf("Order %s already exsists belongs to another user", e.orderID)
}
