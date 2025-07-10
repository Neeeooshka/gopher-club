-- name: AddOrder :one
INSERT INTO orders (user_id, num) VALUES ($1, $2)
ON CONFLICT (num) DO UPDATE SET num = EXCLUDED.num
RETURNING *, (xmax = 0) AS is_new;

-- name: ListUserOrders :many
SELECT * FROM orders WHERE user_id = $1 ORDER BY date_insert DESC;
