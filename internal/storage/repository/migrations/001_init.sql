CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    balance NUMERIC(10, 2) DEFAULT 0 CHECK (balance >= 0)
);
CREATE UNIQUE INDEX users_login_idx ON users (login);

CREATE TABLE user_key_value (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users (id),
    p_name TEXT,
    p_value TEXT
);
CREATE UNIQUE INDEX users_param_idx ON user_key_value (
    user_id, p_name
);

CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'PROCESSED', 'INVALID');

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users (id),
    num TEXT NOT NULL UNIQUE,
    date_insert TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    accrual NUMERIC(10, 2) DEFAULT 0,
    status ORDER_STATUS DEFAULT 'NEW'
);
CREATE INDEX orders_user_id_idx ON orders (user_id);
CREATE UNIQUE INDEX orders_number_idx ON orders (num);

CREATE TABLE withdrawals (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users (id),
    num TEXT NOT NULL,
    date_withdraw TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    sum NUMERIC(10, 2) NOT NULL CHECK (sum > 0)
);
CREATE UNIQUE INDEX withdrawals_user_order_idx ON withdrawals (
    user_id, num
);
CREATE INDEX withdrawals_user_id_idx ON withdrawals (
    user_id
);
CREATE INDEX withdrawals_num_idx ON withdrawals (num);

---- create above / drop below ----

DROP TABLE users;
DROP TABLE user_key_value;
DROP TABLE orders;
DROP TABLE withdrawals;
