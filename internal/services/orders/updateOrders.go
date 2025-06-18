package orders

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/config"
	"github.com/Neeeooshka/gopher-club/internal/logger/zap"
	"net/http"
	"resty.dev/v3"
	"time"
)

type OrdersUpdateRepository interface {
	UpdateOrders(context.Context, []Order) error
	ListWaitingOrders(context.Context) ([]Order, error)
}

type OrdersUpdateService struct {
	ctx            context.Context
	opt            config.Options
	storage        OrdersUpdateRepository
	updateInterval time.Duration
	waitingOrders  []Order
}

func NewOrdersUpdateService(ctx context.Context, our interface{}, opt config.Options) (OrdersUpdateService, error) {

	var ous OrdersUpdateService

	repo, ok := our.(OrdersUpdateRepository)

	if !ok {
		return ous, fmt.Errorf("2th argument expected OrdersUpdateRepository, got %T", our)
	}

	orders, err := repo.ListWaitingOrders(ctx)
	if err != nil {
		return ous, fmt.Errorf("unable to request order details: %w", err)
	}

	ous.ctx = ctx
	ous.opt = opt
	ous.storage = repo
	ous.updateInterval = time.Minute * 5
	ous.waitingOrders = orders

	go ous.ordersUpdater()

	return ous, nil
}

func (o *OrdersUpdateService) AddWaitingOrder(order Order) {
	o.waitingOrders = append(o.waitingOrders, order)
}

func (o *OrdersUpdateService) ordersUpdater() {

	timer := time.NewTicker(o.updateInterval)
	defer timer.Stop()

	for {
		select {
		case <-o.ctx.Done():
			return
		case <-timer.C:
			o.updateOrders()
		}
	}
}

func (o *OrdersUpdateService) updateOrders() {

	// do nothing
	if len(o.waitingOrders) == 0 {
		return
	}

	logger, _ := zap.NewZapLogger("debug")

	ordersForUpdate := make([]Order, 0, len(o.waitingOrders))  // slice for update
	newWaitingOrders := make([]Order, 0, len(o.waitingOrders)) // new slice of waitingOrders

	// save updated orders
	defer func() {
		err := o.storage.UpdateOrders(o.ctx, ordersForUpdate)
		if err != nil {
			logger.Debug(fmt.Sprintf("cannot update orders from the Loyalty calculation system"), logger.Error(err))
			return
		}

		o.waitingOrders = newWaitingOrders
	}()

	type orderInfo struct {
		Number  string  `json:"order"`
		Status  string  `json:"status"`
		Accrual float64 `json:"accrual"`
	}

	client := resty.New()
	defer client.Close()

	r := client.R()

	for _, order := range o.waitingOrders {
		res, err := r.Get(fmt.Sprintf(o.opt.GetAccrualSystem()+"/api/orders/%s", order.Number))
		if err != nil {
			logger.Debug(fmt.Sprintf("cannot connect to the Loyalty calculation system"), logger.Error(err))
			return
		}

		if res.StatusCode() == http.StatusNoContent {
			logger.Debug(fmt.Sprintf("order is not find in the Loyalty calculation system: %s", order.Number))
			continue
		}

		if res.StatusCode() == http.StatusTooManyRequests {
			return
		}

		if res.StatusCode() != http.StatusOK {
			logger.Debug(fmt.Sprintf("the Loyalty calculation system return an unexpected status code: %d", res.StatusCode()))
			return
		}

		o := orderInfo{}

		if err := json.NewDecoder(res.Body).Decode(&o); err != nil {
			logger.Debug(fmt.Sprintf("cannot deserialize response from the Loyalty calculation system: %v", err), logger.Error(err))
			return
		}

		// if there changes
		if order.Accrual != o.Accrual || order.Status != o.Status {

			order.Accrual = o.Accrual
			order.Status = o.Status

			ordersForUpdate = append(ordersForUpdate, order)

			// exclude finish states for next time
			if order.Status != "PROCESSED" && order.Status != "INVALID" {
				newWaitingOrders = append(newWaitingOrders, order)
			}
		}
	}
}
