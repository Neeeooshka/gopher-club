-- name: GetUserByLogin :one
select u.*, up.p_value as credentials from gopher_users u
join gopher_user_params up on up.user_id = u.id and p_name = 'credentials'
where u.login = $1
limit 1;

-- name: AddUser :one
with ins as (
    insert into gopher_users (login, password) values ($1, $2)
    on conflict (login) do nothing
    returning id
)
select id, true as is_new from ins
union all
select id, false as is_new from gopher_users where login = $1
limit 1;

-- name: AddCredentials :exec
insert into gopher_user_params (user_id, p_name, p_value) values ($1, 'credentials', $2);

-- name: UpdateBalance :exec
update gopher_users set balance = balance + $1 where id = $2;
