package models

type key string

const UserContextKey key = "AuthUser"

type User struct {
	ID          int     `db:"ID"`
	Login       string  `db:"login"`
	Password    string  `db:"password"`
	Credentials string  `db:"credentials"`
	Balance     float32 `db:"balance"`
}
