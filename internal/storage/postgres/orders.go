package postgres

import (
	"context"
	"database/sql"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
)

func (l *Postgres) AddOrder(number string, userID int) (orders.Order, error) {

	var order orders.Order
	var isNew bool

	row := l.DB.QueryRow("WITH ins AS (\n    INSERT INTO gopher_orders (user_id, num)\n    VALUES ($1, $2)\n    ON CONFLICT (num) DO NOTHING\n    RETURNING *, 1 AS is_new\n)\nSELECT * FROM ins\nUNION  ALL\nSELECT *, 0 AS is_new FROM gopher_orders WHERE num = $2\nLIMIT 1", userID, number)
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Number,
		&order.DateInsert,
		&order.Accrual,
		&order.Status,
		&isNew,
	)
	if err != nil {
		return order, err
	}

	if !isNew {
		if order.UserID == userID {
			return order, orders.NewConflictOrderError(number)
		}
		return order, orders.NewConflictOrderUserError(order.UserID, number)
	}

	return order, nil
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
