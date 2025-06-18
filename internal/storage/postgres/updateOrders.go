package postgres

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/services/orders"
)

func (l *Postgres) ListWaitingOrders(ctx context.Context) ([]orders.Order, error) {

	rows, err := l.DB.QueryContext(ctx, "select * from gopher_orders where status not in ('INVALID', 'PROCESSED')")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return l.extractOrders(rows)
}
