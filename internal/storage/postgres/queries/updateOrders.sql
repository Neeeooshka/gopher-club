-- name: ListWaitingOrders :many
select * from gopher_orders where status not in ('INVALID', 'PROCESSED');

-- name: UpdateOrders :batchexec
update gopher_orders set status = $1, accrual = $2 where id = $3;
