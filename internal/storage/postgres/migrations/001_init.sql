CREATE TABLE IF NOT EXISTS gopher_users (
    id SERIAL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    balance NUMERIC(10, 2) DEFAULT 0 CHECK (balance >= 0)
);
CREATE UNIQUE INDEX IF NOT EXISTS users_login_idx ON gopher_users (login);

CREATE TABLE IF NOT EXISTS gopher_user_params (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES gopher_users(id),
    p_name TEXT,
    p_value TEXT
);
CREATE UNIQUE INDEX IF NOT EXISTS users_param_idx ON gopher_user_params (user_id, p_name);

CREATE TABLE IF NOT EXISTS gopher_orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES gopher_users(id),
    num TEXT NOT NULL UNIQUE,
    date_insert TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    accrual NUMERIC(10, 2) DEFAULT 0,
    status TEXT DEFAULT 'NEW'
);
CREATE INDEX IF NOT EXISTS orders_user_id_idx ON gopher_orders (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS orders_number_idx ON gopher_orders (num);

CREATE TABLE IF NOT EXISTS gopher_withdrawals (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES gopher_users(id),
    num TEXT NOT NULL,
    date_withdraw TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    sum NUMERIC(10, 2) NOT NULL CHECK (sum > 0)
);
CREATE UNIQUE INDEX IF NOT EXISTS withdrawals_user_order_idx ON gopher_withdrawals (user_id, num);
CREATE INDEX IF NOT EXISTS withdrawals_user_id_idx ON gopher_withdrawals (user_id);
CREATE INDEX IF NOT EXISTS withdrawals_num_idx ON gopher_withdrawals (num);

---- create above / drop below ----

DROP TABLE IF EXISTS gopher_users;
DROP TABLE IF EXISTS gopher_user_params;
DROP TABLE IF EXISTS gopher_orders;
DROP TABLE IF EXISTS gopher_withdrawals;
