package models

import "time"

type Order struct {
	ID         int       `db:"id"`
	UserID     int       `db:"user_id"`
	Number     string    `db:"num" json:"number"`
	DateInsert time.Time `db:"date_insert" json:"uploaded_at"`
	Accrual    float64   `db:"accrual" json:"accrual,omitempty"`
	Status     string    `db:"status" json:"status"`
}
