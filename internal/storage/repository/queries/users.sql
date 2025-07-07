-- name: GetUserByLogin :one
SELECT u.*, up.p_value AS credentials FROM users u
    JOIN user_key_value up ON up.user_id = u.id AND p_name = 'credentials'
WHERE u.login = $1
LIMIT 1;

-- name: AddUser :one
WITH ins AS (
    INSERT INTO users (login, password) VALUES ($1, $2)
    ON CONFLICT (login) DO NOTHING
    RETURNING id
)

SELECT id, TRUE AS is_new FROM ins
UNION ALL
SELECT id, FALSE AS is_new FROM users WHERE login = $1
LIMIT 1;

-- name: AddCredentials :exec
INSERT INTO user_key_value (user_id, p_name, p_value) VALUES (
    $1, 'credentials', $2
);

-- name: UpdateBalance :exec
UPDATE users SET balance = balance + $1 WHERE id = $2;
