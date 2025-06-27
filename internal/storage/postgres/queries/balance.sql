-- name: GetWithdrawals :many
select * from gopher_withdrawals where user_id = $1 order by date_withdraw desc;

-- name: WithdrawBalance :exec
insert into gopher_withdrawals (user_id, num, sum) values ($1, $2, $3);

-- name: GetWithdrawn :one
select sum(sum) as withdrawn from gopher_withdrawals where user_id = $1 group by user_id;
