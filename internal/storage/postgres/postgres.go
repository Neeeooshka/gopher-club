package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

type Postgres struct {
	DB *sql.DB
}

func (l *Postgres) Close() error {
	return l.DB.Close()
}

func NewPostgresStorage(conn string) (pgx *Postgres, err error) {

	pgx = &Postgres{}

	pgx.DB, err = sql.Open("pgx", conn)
	if err != nil {
		return nil, err
	}

	return pgx, nil
}

func (l *Postgres) PingHandler(w http.ResponseWriter, _ *http.Request) {
	err := l.DB.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Bootstrap execute first initialization DB, creating tables and indexes
func (l *Postgres) Bootstrap(ctx context.Context) error {

	tx, err := l.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS gopher_users (
			id SERIAL PRIMARY KEY,
			login TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			balance NUMERIC(10, 2) DEFAULT 0 CHECK (balance >= 0)
		);
		CREATE UNIQUE INDEX IF NOT EXISTS users_login_idx ON gopher_users (login);
	`)
	if err != nil {
		return fmt.Errorf("failed to create gopher_users table: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS gopher_user_params (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES gopher_users(id),
			p_name TEXT,
			p_value TEXT,
		);
		CREATE UNIQUE INDEX IF NOT EXISTS users_param_idx ON gopher_user_params (user_id, p_name);
	`)
	if err != nil {
		return fmt.Errorf("failed to create gopher_user_params table: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS gopher_orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES gopher_users(id),
			num TEXT NOT NULL UNIQUE,
			date_insert TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			accrual NUMERIC(10, 2) DEFAULT 0,
			status TEXT DEFAULT 'NEW',
		);
		CREATE INDEX IF NOT EXISTS orders_user_id_idx ON gopher_orders (user_id);
		CREATE UNIQUE INDEX IF NOT EXISTS orders_number_idx ON gopher_orders (num);
	`)
	if err != nil {
		return fmt.Errorf("failed to create gopher_orders table: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS gopher_withdrawals (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES gopher_users(id),
			order_id TEXT NOT NULL REFERENCES gopher_orders(id),
			date_withdraw TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			sum NUMERIC(10, 2) NOT NULL CHECK (sum > 0)
		);
		CREATE UNIQUE INDEX IF NOT EXISTS withdrawals_user_order_idx ON gopher_withdrawals (user_id, order_id);
		CREATE INDEX IF NOT EXISTS withdrawals_user_id_idx ON gopher_withdrawals (user_id);
		CREATE INDEX IF NOT EXISTS withdrawals_order_id_idx ON gopher_withdrawals (order_id);
	`)
	if err != nil {
		return fmt.Errorf("failed to create gopher_withdrawals table: %w", err)
	}

	return tx.Commit()
}
