package orders

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/config"
	"github.com/Neeeooshka/gopher-club/internal/logger/zap"
	"github.com/Neeeooshka/gopher-club/internal/models"
	"net/http"
	"resty.dev/v3"
	"time"
)

type OrdersUpdateRepository interface {
	UpdateOrders(context.Context, []models.Order) error
	ListWaitingOrders(context.Context) ([]models.Order, error)
}

type OrdersUpdateService struct {
	logger         *zap.ZapLogger
	opt            config.Options
	storage        OrdersUpdateRepository
	updateInterval time.Duration
	waitingOrders  []models.Order
	isRunning      bool
}

func NewOrdersUpdateService(our interface{}, opt config.Options) (*OrdersUpdateService, error) {

	var ous OrdersUpdateService

	repo, ok := our.(OrdersUpdateRepository)

	if !ok {
		return nil, fmt.Errorf("2th argument expected OrdersUpdateRepository, got %T", our)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	orders, err := repo.ListWaitingOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to request order details: %w", err)
	}

	logger, err := zap.NewZapLogger("debug")
	if err != nil {
		return nil, fmt.Errorf("unable to initialize logger: %w", err)
	}

	ous.logger = logger
	ous.opt = opt
	ous.storage = repo
	ous.updateInterval = time.Second
	ous.waitingOrders = orders

	go ous.ordersUpdater()

	return &ous, nil
}

func (o *OrdersUpdateService) AddWaitingOrder(order models.Order) {
	o.waitingOrders = append(o.waitingOrders, order)
}

func (o *OrdersUpdateService) ordersUpdater() {

	timer := time.NewTicker(o.updateInterval)
	defer timer.Stop()

	for range timer.C {
		o.updateOrders()
	}
}

// updateOrders pull new status and accrual for all waiting orders from the Loyalty calculation system
// and send it to updateOrdersProcessor if there changes
func (o *OrdersUpdateService) updateOrders() {

	// do nothing
	if len(o.waitingOrders) == 0 || o.isRunning {
		return
	}

	// lock process
	o.start()

	dataCh := make(chan models.Order)
	defer close(dataCh)

	go o.updateOrdersProcessor(dataCh)

	type orderInfo struct {
		Number  string  `json:"order"`
		Status  string  `json:"status"`
		Accrual float32 `json:"accrual"`
	}

	client := resty.New()
	defer func() {
		if err := client.Close(); err != nil {
			logger, _ := zap.NewZapLogger("debug")
			logger.Debug("failed to close resty client", logger.Error(err))
		}
	}()

	r := client.R()

	for _, order := range o.waitingOrders {
		res, err := r.Get(fmt.Sprintf(o.opt.GetAccrualSystem()+"/api/orders/%s", order.Number))
		if err != nil {
			o.logger.Debug("cannot connect to the Loyalty calculation system", o.logger.Error(err))
			return
		}

		if res.StatusCode() == http.StatusNoContent {
			o.logger.Debug(fmt.Sprintf("order is not find in the Loyalty calculation system: %s", order.Number))
			continue
		}

		if res.StatusCode() == http.StatusTooManyRequests {
			return
		}

		if res.StatusCode() != http.StatusOK {
			o.logger.Debug(fmt.Sprintf("the Loyalty calculation system return an unexpected status code: %d", res.StatusCode()))
			return
		}

		oi := orderInfo{}

		if err := json.NewDecoder(res.Body).Decode(&oi); err != nil {
			o.logger.Debug(fmt.Sprintf("cannot deserialize response from the Loyalty calculation system: %v", err), o.logger.Error(err))
			return
		}

		// if there changes
		if order.Accrual != oi.Accrual || order.Status != oi.Status {

			// save order memento
			order.CreateMemento("beforeUpdate")

			order.Accrual = oi.Accrual
			order.Status = oi.Status

			dataCh <- order
		}
	}
}

func (o *OrdersUpdateService) updateOrdersProcessor(dataCh chan models.Order) {

	ordersForUpdateMap := make(map[string]models.Order)

	for order := range dataCh {
		ordersForUpdateMap[order.Number] = order
	}

	if len(ordersForUpdateMap) > 0 {
		o.applyUpdates(ordersForUpdateMap)
	}

	// unlock process
	o.stop()
}

func (o *OrdersUpdateService) applyUpdates(ordersForUpdateMap map[string]models.Order) {

	var ordersForUpdate []models.Order
	newWaitingOrders := make([]models.Order, 0, len(o.waitingOrders))

	for _, order := range o.waitingOrders {
		if ord, ok := ordersForUpdateMap[order.Number]; ok {
			ordersForUpdate = append(ordersForUpdate, ord)
			if ord.Status == StatusProcessed || ord.Status == StatusInvalid {
				continue
			}

			order = ord
		}

		newWaitingOrders = append(newWaitingOrders, order)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := o.storage.UpdateOrders(ctx, ordersForUpdate); err != nil {
		o.logger.Debug("cannot update orders from the Loyalty calculation system", o.logger.Error(err))
		return
	}

	o.waitingOrders = newWaitingOrders
}

func (o *OrdersUpdateService) start() {
	o.isRunning = true
}

func (o *OrdersUpdateService) stop() {
	o.isRunning = false
}
