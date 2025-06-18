package storage

import (
	"github.com/Neeeooshka/gopher-club/internal/services/users"
)

type Storage interface {
	Close() error
}

type UserRepository interface {
	AddUser(users.User, string) error
	GetUserByLogin(string) (users.User, error)
}
