package orders

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/config"
	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/pkg/httputil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"net/http"
	"testing"
	"time"
)

const accrualSystemAddress = "localhost"
const accrualSystemPort = "8899"

type accrualSystemOrder struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

// accrualSystem emulates server operation
type accrualSystem struct {
	URI    string
	Port   string
	orders map[string]accrualSystemOrder
}

// newAccrualSystem fills some orders
func newAccrualSystem(u string, p string) *accrualSystem {

	orders := make(map[string]accrualSystemOrder, 3)
	orders["4532015112830366"] = accrualSystemOrder{"4532015112830366", "PROCESS", 0}
	orders["6011324432123452"] = accrualSystemOrder{"6011324432123452", "PROCESSED", 729.98}
	orders["378282246310005"] = accrualSystemOrder{"378282246310005", "INVALID", 0}

	return &accrualSystem{u, p, orders}
}

// accrualOrdersHandler gets order info from map
func (a *accrualSystem) accrualOrdersHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	o, ok := a.orders[id]
	if !ok {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	httputil.WriteJSON(w, o)
}

func (a *accrualSystem) start(stopCh chan struct{}) {

	mux := http.NewServeMux()
	mux.HandleFunc("/api/orders/{id}", a.accrualOrdersHandler)

	server := &http.Server{Addr: a.URI + ":" + a.Port, Handler: mux}

	go func() {
		<-stopCh
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()

	log.Fatal(server.ListenAndServe())
}

type MockOrdersUpdateRepository struct {
	mock.Mock
}

func (m *MockOrdersUpdateRepository) UpdateOrders(ctx context.Context, orders []models.Order) error {
	args := m.Called(ctx, orders)
	return args.Error(0)
}

func (m *MockOrdersUpdateRepository) ListWaitingOrders(ctx context.Context) ([]models.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Order), args.Error(1)
}

func TestUpdateOrders(t *testing.T) {
	mockRepo := new(MockOrdersUpdateRepository)
	mockRepo.On("UpdateOrders", mock.Anything, mock.Anything).Return(nil)

	as := newAccrualSystem(accrualSystemAddress, accrualSystemPort)
	stopCh := make(chan struct{})
	defer close(stopCh)

	go as.start(stopCh)

	// waiting for the Accrual system to start
	time.Sleep(time.Second)

	var opt config.Options
	err := opt.AccrualAddress.Set("http://" + accrualSystemAddress + ":" + accrualSystemPort)
	assert.NoError(t, err)

	service := OrdersUpdateService{
		storage:        mockRepo,
		waitingOrders:  []models.Order{{Number: "4532015112830366", Status: StatusNew}},
		opt:            opt,
		updateInterval: time.Second,
	}

	service.updateOrders()

	// waiting for all orders to be written
	time.Sleep(time.Second * 1)

	mockRepo.AssertExpectations(t)
}

func TestUpdateOrdersProcessor(t *testing.T) {
	mockRepo := new(MockOrdersUpdateRepository)
	mockRepo.On("UpdateOrders", mock.Anything, mock.Anything).Return(nil)

	order := models.Order{Number: "4532015112830366", Status: StatusProcessed}

	service := OrdersUpdateService{
		storage:       mockRepo,
		waitingOrders: []models.Order{{Number: order.Number, Status: StatusNew}},
	}

	dataCh := make(chan models.Order)
	go service.updateOrdersProcessor(dataCh)

	dataCh <- order
	close(dataCh)

	// waiting for all gorutines
	time.Sleep(time.Second * 1)

	mockRepo.AssertCalled(t, "UpdateOrders", mock.Anything, []models.Order{order})
}

func TestApplyUpdates(t *testing.T) {
	mockRepo := new(MockOrdersUpdateRepository)
	mockRepo.On("UpdateOrders", mock.Anything, mock.Anything).Return(nil)

	order := models.Order{Number: "4532015112830366", Status: StatusProcessing}

	service := OrdersUpdateService{
		storage:       mockRepo,
		waitingOrders: []models.Order{order},
	}

	updates := map[string]models.Order{
		order.Number: {Number: order.Number, Status: StatusProcessed},
	}

	service.applyUpdates(updates)
	assert.Len(t, service.waitingOrders, 0)
	mockRepo.AssertExpectations(t)
}
