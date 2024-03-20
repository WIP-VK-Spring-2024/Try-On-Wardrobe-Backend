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
