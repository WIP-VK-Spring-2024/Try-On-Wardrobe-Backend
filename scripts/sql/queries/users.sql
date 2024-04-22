-- name: CreateUser :one
insert into users(
    name,
    email,
    password,
    gender
) values ($1, $2, $3, $4)
returning id;

-- name: GetUserByID :one
select * from users
where id = $1;

-- name: GetUserByName :one
select * from users
where name = $1;

-- name: GetUserByEmail :one
select * from users
where lower(email) = lower($1);

-- name: GetSubscribedToUsers :many
select users.*
from users
join subs on subs.subscriber_id = $1
     and subs.user_id = users.id;

-- name: SearchUsers :many
select users.*
from users
where lower(name) like lower($1);

-- name: UpdateUser :exec
update users
set name = case when sqlc.arg(name)::text = '' then name
                else sqlc.arg(name)::text end,
    gender = coalesce($1, gender),
    privacy = coalesce($2, privacy),
    avatar = case when sqlc.arg(avatar)::text = '' then name
                  else sqlc.arg(avatar)::text end,
    updated_at = now()
where id = $1;
