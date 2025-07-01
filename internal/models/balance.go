package models

import "time"

type Withdraw struct {
	ID           int       `db:"id" json:"-"`
	UserID       int       `db:"user_id" json:"-"`
	OrderNum     string    `db:"number" json:"order"`
	DateWithdraw time.Time `db:"date_withdraw" json:"processed_at"`
	Sum          float32   `db:"sum" json:"sum"`
}
