-- name: GetUserByLogin :one
SELECT * FROM users u WHERE u.login = $1 LIMIT 1;

-- name: AuthUser :one
SELECT * FROM users u WHERE u.login = $1 AND password = crypt($2, password) LIMIT 1;

-- name: AddUser :one
INSERT INTO users (login, password) VALUES ($1, crypt($2, gen_salt('des')))
ON CONFLICT (login) DO UPDATE SET login = EXCLUDED.login
RETURNING id, (xmax = 0) AS is_new;

-- name: UpdateBalance :exec
UPDATE users SET balance = balance + $1 WHERE id = $2;
