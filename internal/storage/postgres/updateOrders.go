package postgres

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
)

func (l *Postgres) ListWaitingOrders(ctx context.Context) ([]orders.Order, error) {

	finishedStates := []string{orders.StatusInvalid, orders.StatusProcessed}

	rows, err := l.DB.QueryContext(ctx, "select * from gopher_orders where status not in ($1)", finishedStates)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return l.extractOrders(rows)
}

func (l *Postgres) UpdateOrders(ctx context.Context, orders []orders.Order) error {

	tx, err := l.DB.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "update gopher_orders set status = $1, accrual = $2 where order_id = $3")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, order := range orders {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_, err := stmt.ExecContext(ctx, order.Status, order.Accrual, order.ID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
