package mock

import (
	"context"
	"errors"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
)

func (m *MockRepository) AddUser(_ context.Context, user users.User, salt string) error {
	if _, exists := m.Users[user.Login]; exists {
		return users.NewConflictUserError(1, user.Login)
	}
	user.Credentials = salt
	m.Users[user.Login] = user
	return nil
}

func (m *MockRepository) GetUserByLogin(login string) (users.User, error) {
	user, exists := m.Users[login]
	if !exists {
		return users.User{}, errors.New("user not found")
	}
	return user, nil
}
