package mocks

import (
	"context"
	"errors"
	"github.com/Neeeooshka/gopher-club/internal/services/models"
	"github.com/Neeeooshka/gopher-club/internal/storage"
)

func (m *MockRepository) AddUser(_ context.Context, user models.User, salt string) error {
	if _, exists := m.Users[user.Login]; exists {
		return storage.NewConflictUserError(1, user.Login)
	}
	user.Credentials = salt
	m.Users[user.Login] = user
	return nil
}

func (m *MockRepository) GetUserByLogin(login string) (models.User, error) {
	user, exists := m.Users[login]
	if !exists {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}
