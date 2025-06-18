package postgres

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
	"github.com/Neeeooshka/gopher-club/internal/services/users"
)

func (l *Postgres) AddOrder(number string, userId int) error {

	var id int
	var isNew bool
	var uID int

	row := l.DB.QueryRow("WITH ins AS (\n    INSERT INTO gopher_orders (user_id, number)\n    VALUES ($1, $2)\n    ON CONFLICT (number) DO NOTHING\n        RETURNING user_id\n)\nSELECT id, 1 as is_new, $1 as user_id FROM ins\nUNION  ALL\nSELECT id, 0 as is_new, user_id FROM gopher_users WHERE number = $2\nLIMIT 1", userId, number)
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

func (l *Postgres) ListOrders(ctx context.Context, user users.User) ([]orders.Order, error) {

	var result []orders.Order

	rows, err := l.DB.QueryContext(ctx, "select * from gopher_orders where user_id = $1 order by date_insert desc", user.ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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
