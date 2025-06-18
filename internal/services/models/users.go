package models

type User struct {
	ID          int     `db:"ID"`
	Login       string  `db:"login"`
	Password    string  `db:"password"`
	Credentials string  `db:"credentials"`
	Balance     float64 `db:"balance"`
}
