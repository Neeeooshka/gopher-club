package postgres

import (
	"context"
	"fmt"
	"github.com/Neeeooshka/gopher-club/internal/models"
)

func (l *Postgres) ListWaitingOrders(ctx context.Context) ([]models.Order, error) {

	rows, err := l.DB.QueryContext(ctx, "select * from gopher_orders where status not in ('INVALID', 'PROCESSED')")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return l.extractOrders(rows)
}

func (l *Postgres) UpdateOrders(ctx context.Context, orders []models.Order) error {

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
			// update order
			_, err := stmt.ExecContext(ctx, order.Status, order.Accrual, order.ID)
			if err != nil {
				return err
			}

			// get accrual before update
			orderMementoBefore := order.GetMemento("beforeUpdate")
			if orderMementoBefore != nil {
				tx.Rollback()
				return fmt.Errorf("cannot get accrual before update")
			}

			// add balance to user
			addBalance := order.Accrual - orderMementoBefore.GetAccrual()
			if addBalance > 0 {
				_, err = tx.ExecContext(ctx, "update gopher_users set balance = balance + $1 where user_id = $2", addBalance, order.ID)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	return tx.Commit()
}
