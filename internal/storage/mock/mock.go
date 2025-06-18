package mock

import "github.com/Neeeooshka/gopher-club/internal/services/users"

type MockRepository struct {
	Users map[string]users.User
}

func (m *MockRepository) Close() error {
	return nil
}
