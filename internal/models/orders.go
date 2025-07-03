package models

import (
	"time"
)

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID         int         `db:"id" json:"-"`
	UserID     int         `db:"user_id" json:"-"`
	Number     string      `db:"num" json:"number"`
	DateInsert time.Time   `db:"date_insert" json:"uploaded_at"`
	Accrual    float32     `db:"accrual" json:"accrual,omitempty"`
	Status     OrderStatus `db:"status" json:"status"`
	mementos   map[string]orderMemento
}

func (o *Order) CreateMemento(state string) {
	if o.mementos == nil {
		o.mementos = make(map[string]orderMemento)
	}
	o.mementos[state] = orderMemento{accrual: o.Accrual, status: o.Status}
}

func (o *Order) GetMemento(state string) (orderMemento, bool) {
	memento, ok := o.mementos[state]
	return memento, ok
}

type orderMemento struct {
	accrual float32
	status  OrderStatus
}

func (m *orderMemento) GetAccrual() float32 {
	return m.accrual
}
