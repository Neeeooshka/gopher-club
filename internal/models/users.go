package models

import "github.com/shopspring/decimal"

type User struct {
	ID          int             `db:"ID"`
	Login       string          `db:"login"`
	Password    string          `db:"password"`
	Credentials string          `db:"credentials"`
	Balance     decimal.Decimal `db:"balance"`
}
