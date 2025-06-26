package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type Withdraw struct {
	ID           int             `db:"id"`
	UserID       int             `db:"user_id"`
	OrderNum     string          `db:"number" json:"order"`
	DateWithdraw time.Time       `db:"date_withdraw" json:"processed_at"`
	Sum          decimal.Decimal `db:"sum" json:"sum"`
}
