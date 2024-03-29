-- name: CreateUserImage :one
insert into user_images(user_id)
values ($1)
returning id;

-- name: DeleteUserImage :exec
delete from user_images
where id = $1;

-- name: GetUserImageByUser :many
select * from user_images
where user_id = $1;

-- name: GetUserImageByID :one
select * from user_images
where id = $1;

-- name: SetUserImageUrl :exec
update user_images
set image = $2
where id = $1;
