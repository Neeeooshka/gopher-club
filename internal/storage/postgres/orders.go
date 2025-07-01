package postgres

import (
	"context"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/models"
	"github.com/Neeeooshka/gopher-club/internal/storage"
	"github.com/Neeeooshka/gopher-club/internal/storage/postgres/sqlc"
)

func (s *Postgres) AddOrder(ctx context.Context, number string, userID int) (models.Order, error) {

	var order models.Order

	result, err := s.sqlc.AddOrder(ctx, sqlc.AddOrderParams{UserID: userID, Num: number})
	if err != nil {
		return order, fmt.Errorf("error adding order: %w", err)
	}

	if !result.IsNew {
		if result.UserID == userID {
			return order, storage.NewConflictOrderError(number)
		}
		return order, storage.NewConflictOrderUserError(result.UserID, number)
	}

	order.ID = result.ID
	order.UserID = result.UserID
	order.Number = result.Num
	order.Status = result.Status
	order.Accrual = result.Accrual
	order.DateInsert = result.DateInsert

	return order, nil
}

func (s *Postgres) ListUserOrders(ctx context.Context, user models.User) ([]models.Order, error) {

	results, err := s.sqlc.ListUserOrders(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("error listing user orders: %w", err)
	}

	return s.extractOrders(results), nil
}

func (s *Postgres) extractOrders(results []sqlc.GopherOrder) []models.Order {

	orders := make([]models.Order, len(results))

	for i, result := range results {
		orders[i] = models.Order{
			ID:         result.ID,
			UserID:     result.UserID,
			Number:     result.Num,
			Status:     result.Status,
			Accrual:    result.Accrual,
			DateInsert: result.DateInsert,
		}
	}

	return orders
}
