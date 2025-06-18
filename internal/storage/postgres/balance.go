package postgres

import (
	"context"
	"github.com/Neeeooshka/gopher-club/internal/services/balance"
	"github.com/Neeeooshka/gopher-club/internal/services/models"
)

func (l *Postgres) GetWithdrawals(ctx context.Context, user models.User) ([]balance.Withdraw, error) {

	var withdrawals []balance.Withdraw

	rows, err := l.DB.QueryContext(ctx, "select * from gopher_withdrawals where user_id = $1 order by date_withdraw desc", user.ID)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var withdraw balance.Withdraw
		err := rows.Scan(&withdraw)
		if err != nil {
			return nil, err
		}
	}

	return withdrawals, nil
}

func (l *Postgres) WithdrawBalance(ctx context.Context, w balance.Withdraw) error {

	tx, err := l.DB.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "insert into gopher_withdrawals (user_id, order_id, sum) values ($1, $2, $3)", w.UserID, w.OrderID, w.Sum)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "update gopher_users set balance = balance + $1 where user_id = $2", w.Sum, w.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (l *Postgres) GetWithdrawn(ctx context.Context, user models.User) (float64, error) {

	var withdrawn float64

	row := l.DB.QueryRowContext(ctx, "select sum(sum) as withdrawn from gopher_withdrawals where user_id = $1 group by user_id", user.ID)
	if err := row.Scan(&withdrawn); err != nil {
		return 0, err
	}

	return withdrawn, nil
}
