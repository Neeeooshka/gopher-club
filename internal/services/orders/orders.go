package orders

import (
	"context"
	"errors"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
	"io"
	"net/http"
	"strconv"
)

type OrdersRepository interface {
	AddOrder(string, int) error
}

type OrdersService struct {
	Inited      bool
	storage     OrdersRepository
	ctx         context.Context
	Errors      []error
	UserService *users.UserService
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

	os.ctx = ctx
	os.storage = ordersRepo
	os.UserService = us

	if len(os.Errors) == 0 {
		os.Inited = true
	}

	return os
}

func (o *OrdersService) AddUserOrderHandler(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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
	if !o.checkLuhn(orderNumber) {
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

func (o *OrdersService) checkLuhn(orderNumber string) bool {

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
