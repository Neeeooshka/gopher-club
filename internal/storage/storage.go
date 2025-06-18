package storage

import "github.com/Neeeooshka/gopher-club/internal/users"

type Storage interface {
	UserRepository
	Close() error
}

type UserRepository interface {
	AddUser(user users.User) error
	GetUserByLogin(login string) (users.User, error)
	GetUserKey(ID int) (string, error)
}
