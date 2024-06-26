-- name: CreateUser :one
insert into users(
    name,
    email,
    password,
    gender,
    privacy
) values ($1, $2, $3, $4, 'private')
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
left join subs on subs.subscriber_id = $1
     and subs.user_id = users.id
where lower(name) like lower(sqlc.arg(name))
      and lower(name) > sqlc.arg(since)
      and subs.user_id is null
      and users.id <> $1
order by lower(name)
limit $2;

-- name: UpdateUser :exec
update users
set name = case when sqlc.arg(name)::text = '' then name
                else sqlc.arg(name)::text end,
    gender = coalesce(sqlc.narg(gender), gender),
    privacy = coalesce(sqlc.narg(privacy), privacy),
    avatar = case when sqlc.arg(avatar)::text = '' then avatar
                  else sqlc.arg(avatar)::text end,
    updated_at = now()
where id = $1;
