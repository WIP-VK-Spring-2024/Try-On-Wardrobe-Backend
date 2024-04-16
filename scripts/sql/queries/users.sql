-- name: CreateUser :one
insert into users(
    name,
    email,
    password
) values (sqlc.arg(name), sqlc.arg(email), sqlc.arg(password))
returning id;

-- name: GetUserByID :one
select * from users
where id = $1;

-- name: GetUserByName :one
select * from users
where name = $1;

-- name: GetSubscribedToUsers :many
select users.*
from users
join subs on subs.subscriber_id = $1
     and subs.user_id = users.id;

-- name: SearchUsers :many
select users.*
from users
where lower(name) like lower($1);
