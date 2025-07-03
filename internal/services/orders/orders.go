package orders

import (
	"context"
	"errors"
	"strconv"

	"github.com/Neeeooshka/gopher-club/internal/config"
	"github.com/Neeeooshka/gopher-club/internal/models"
)

type OrdersRepository interface {
	AddOrder(context.Context, string, int) (models.Order, error)
	ListUserOrders(context.Context, models.User) ([]models.Order, error)
}

type OrdersService struct {
	errors        []error
	init          bool
	storage       OrdersRepository
	updateService *OrdersUpdateService
}

func NewOrdersService(repo OrdersRepository, opt config.Options) OrdersService {

	var os OrdersService

	ous, err := NewOrdersUpdateService(repo.(OrdersUpdateRepository), opt)

	if err != nil {
		os.errors = append(os.errors, errors.New("cannot initialize OrdersUpdateService"))
	}

	if len(os.errors) > 0 {
		return os
	}

	os.storage = repo
	os.updateService = ous
	os.init = true

	return os
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

func (o *OrdersService) HealthCheck() ([]error, bool) {
	return o.errors, o.init
}

func (o *OrdersService) GetName() string {
	return "OrdersService"
}
