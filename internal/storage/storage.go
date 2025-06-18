package storage

import (
	"github.com/Neeeooshka/gopher-club.git/internal/auth"
)

type Membership interface {
	Close() error
	AddUser(user auth.User) error
}
