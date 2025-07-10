package models

type key string

const UserContextKey key = "AuthUser"

type User struct {
	ID       int     `db:"ID"`
	Login    string  `db:"login"`
	Password string  `db:"password"`
	Balance  float32 `db:"balance"`
}

func (u *User) CanWithdraw(sum float32) bool {
	return u.Balance >= sum
}
