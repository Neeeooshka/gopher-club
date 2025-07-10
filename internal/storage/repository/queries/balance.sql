-- name: GetWithdrawals :many
SELECT * FROM withdrawals WHERE user_id = $1 ORDER BY date_withdraw DESC;

-- name: WithdrawBalance :exec
INSERT INTO withdrawals (user_id, num, sum) VALUES ($1, $2, $3);

-- name: GetWithdrawn :one
SELECT SUM(sum) AS withdrawn FROM withdrawals WHERE user_id = $1 GROUP BY user_id;
