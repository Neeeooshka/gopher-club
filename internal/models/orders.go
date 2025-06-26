package models

import "time"

type Order struct {
	ID         int       `db:"id"`
	UserID     int       `db:"user_id"`
	Number     string    `db:"num" json:"number"`
	DateInsert time.Time `db:"date_insert" json:"uploaded_at"`
	Accrual    float32   `db:"accrual" json:"accrual,omitempty"`
	Status     string    `db:"status" json:"status"`
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
	accrual float32
	status  string
}

func (m *orderMemento) GetAccrual() float32 {
	return m.accrual
}
