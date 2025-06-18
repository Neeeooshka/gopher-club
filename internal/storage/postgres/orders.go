package postgres

import (
	"context"
	"database/sql"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
)

func (l *Postgres) AddOrder(number string, userId int) error {

	var id int
	var isNew bool
	var uID int

	row := l.DB.QueryRow("WITH ins AS (\n    INSERT INTO gopher_orders (user_id, num)\n    VALUES ($1, $2)\n    ON CONFLICT (number) DO NOTHING\n        RETURNING user_id\n)\nSELECT id, 1 as is_new, $1 as user_id FROM ins\nUNION  ALL\nSELECT id, 0 as is_new, user_id FROM gopher_users WHERE number = $2\nLIMIT 1", userId, number)
	err := row.Scan(&id, &isNew, &uID)
	if err != nil {
		return err
	}

	if !isNew {
		if uID == userId {
			return orders.NewConflictOrderError(number)
		}
		return orders.NewConflictOrderUserError(uID, number)
	}

	return nil
}

func (l *Postgres) UpdateOrders(ctx context.Context, orders []orders.Order) error {

	stmt, err := l.DB.BeginTx(ctx, nil)
	defer stmt.Rollback()

	for _, order := range orders {
		_, err = stmt.ExecContext(ctx, "update gopher_orders set status = $1, accrual = $2 where order_id = $3", order.Status, order.Accrual, order.UserID)
		if err != nil {
			return err
		}
	}

	return stmt.Commit()
}

func (l *Postgres) ListOrders(ctx context.Context) ([]orders.Order, error) {

	rows, err := l.DB.QueryContext(ctx, "select * from gopher_orders where status not in ('INVALID', 'PROCESSED')")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return l.extractOrders(rows)
}

func (l *Postgres) ListUserOrders(ctx context.Context, user users.User) ([]orders.Order, error) {

	rows, err := l.DB.QueryContext(ctx, "select * from gopher_orders where user_id = $1 order by date_insert desc", user.ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return l.extractOrders(rows)
}

func (l *Postgres) extractOrders(rows *sql.Rows) ([]orders.Order, error) {

	var result []orders.Order

	for rows.Next() {
		var o orders.Order
		if err := rows.Scan(&o); err != nil {
			return nil, err
		}
		result = append(result, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
