-- name: AddOrder :one
with ins as (
    insert into gopher_orders (user_id, num) values ($1, $2)
    on conflict (num) do nothing
    returning *, true as is_new
)
select * from ins
union all
select *, false as is_new from gopher_orders where num = $2
limit 1;

-- name: ListUserOrders :many
select * from gopher_orders where user_id = $1 order by date_insert desc;
