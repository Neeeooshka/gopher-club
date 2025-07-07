-- name: AddOrder :one
WITH ins AS (
    INSERT INTO orders (user_id, num) VALUES ($1, $2)
    ON CONFLICT (num) DO NOTHING
    RETURNING *, TRUE AS is_new
)

SELECT * FROM ins
UNION ALL
SELECT *, FALSE AS is_new FROM orders WHERE num = $2
LIMIT 1;

-- name: ListUserOrders :many
SELECT * FROM orders WHERE user_id = $1 ORDER BY date_insert DESC;
