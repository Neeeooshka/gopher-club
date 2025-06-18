package storage

type Storage interface {
	Close() error
}

type ConflictUserError struct {
	ID    int
	login string
}

func NewConflictUserError(ID int, login string) *ConflictUserError {
	return &ConflictUserError{ID, login}
}

func (e *ConflictUserError) Error() string {
	return "User with login " + e.login + " already exsists"
}
