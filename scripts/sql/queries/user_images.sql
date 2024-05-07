-- name: CreateUserImage :one
insert into user_images(user_id, image)
values ($1, $2)
returning id;

-- name: DeleteUserImage :exec
delete from user_images
where id = $1;

-- name: GetUserImageByUser :many
select * from user_images
where user_id = $1
order by created_at desc;

-- name: GetUserImageByID :one
select * from user_images
where id = $1;

-- name: SetUserImageUrl :exec
update user_images
set image = $2,
    updated_at = now()
where id = $1;
