-- name: ListWaitingOrders :many
SELECT * FROM orders WHERE status NOT IN ('INVALID', 'PROCESSED');

-- name: UpdateOrders :batchexec
UPDATE orders SET status = $1, accrual = $2 WHERE id = $3;
