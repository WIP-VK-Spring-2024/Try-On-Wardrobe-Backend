-- name: CreateUser :one
insert into users(
    name,
    email,
    password
) values ($1, $2, $3)
returning id;

-- name: GetUserByID :one
select * from users
where id = $1;

-- name: GetUserByName :one
select * from users
where name = $1;
