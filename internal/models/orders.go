package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type Order struct {
	ID         int             `db:"id"`
	UserID     int             `db:"user_id"`
	Number     string          `db:"num" json:"number"`
	DateInsert time.Time       `db:"date_insert" json:"uploaded_at"`
	Accrual    decimal.Decimal `db:"accrual" json:"accrual,omitempty"`
	Status     string          `db:"status" json:"status"`
	mementos   map[string]*orderMemento
}

func (o *Order) CreateMemento(state string) {
	o.mementos[state] = &orderMemento{accrual: o.Accrual, status: o.Status}
}

func (o *Order) GetMemento(state string) *orderMemento {
	memento, ok := o.mementos[state]
	if !ok {
		return nil
	}
	return memento
}

type orderMemento struct {
	accrual decimal.Decimal
	status  string
}

func (m *orderMemento) GetAccrual() decimal.Decimal {
	return m.accrual
}
