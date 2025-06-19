package storage

import "fmt"

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

type ConflictOrderError struct {
	orderID string
}

func NewConflictOrderError(orderID string) *ConflictOrderError {
	return &ConflictOrderError{orderID}
}

func (e *ConflictOrderError) Error() string {
	return fmt.Sprintf("Order %s already exsists", e.orderID)
}

type ConflictOrderUserError struct {
	userID  int
	orderID string
}

func NewConflictOrderUserError(userID int, orderID string) *ConflictOrderUserError {
	return &ConflictOrderUserError{userID, orderID}
}

func (e *ConflictOrderUserError) Error() string {
	return fmt.Sprintf("Order %s already exsists belongs to another user", e.orderID)
}
