CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(25) NOT NULL UNIQUE,
    password VARCHAR(64) NOT NULL,
    balance NUMERIC(10, 2) DEFAULT 0 CHECK (balance >= 0)
);

CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'PROCESSED', 'INVALID');

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users (id),
    num VARCHAR(100) NOT NULL UNIQUE,
    date_insert TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    accrual NUMERIC(10, 2) DEFAULT 0,
    status ORDER_STATUS DEFAULT 'NEW'
);
CREATE INDEX order_user_id_idx ON orders (user_id);

CREATE TABLE withdrawals (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users (id),
    num VARCHAR(100) NOT NULL,
    date_withdraw TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    sum NUMERIC(10, 2) NOT NULL CHECK (sum > 0)
);
CREATE UNIQUE INDEX withdrawal_user_order_idx ON withdrawals (
    user_id, num
);
CREATE INDEX withdrawal_user_id_idx ON withdrawals (
    user_id
);
CREATE INDEX withdrawal_num_idx ON withdrawals (num);

CREATE EXTENSION pgcrypto;

---- create above / drop below ----

DROP TABLE withdrawals CASCADE;
DROP TABLE orders CASCADE;
DROP TABLE users CASCADE;
DROP TYPE order_status;
DROP EXTENSION pgcrypto;
