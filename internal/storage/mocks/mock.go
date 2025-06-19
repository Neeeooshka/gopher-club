package mocks

import (
	"github.com/Neeeooshka/gopher-club/internal/models"
)

type MockRepository struct {
	Users map[string]models.User
}

func (m *MockRepository) Close() error {
	return nil
}
