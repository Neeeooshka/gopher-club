package models

import "time"

type Withdraw struct {
	ID           int       `db:"id"`
	UserID       int       `db:"user_id"`
	OrderID      string    `db:"order_id" json:"order"`
	DateWithdraw time.Time `db:"date_withdraw" json:"processed_at"`
	Sum          float64   `db:"sum" json:"sum"`
}
