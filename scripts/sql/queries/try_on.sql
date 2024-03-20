-- name: CreateTryOnResult :one
insert into try_on_results(
    clothes_id,
    user_image_id
) values ($1, $2)
returning id;

-- name: DeleteTryOnResult :exec
delete from try_on_results
where id = $1;

-- name: GetLastTryOnResult :one
select try_on_results.*
from try_on_results
join user_images u on u.id = try_on_results.user_image_id
where u.user_id = $1
order by try_on_results.created_at desc
limit 1;

-- name: GetTryOnResultsByUser :many
select try_on_results.*
from try_on_results
join user_images u on u.id = try_on_results.user_image_id
where u.user_id = $1;

-- name: GetTryOnResultsByClothes :many
select *
from try_on_results
where clothes_id = $1;

-- name: RateTryOnResult :exec
update try_on_results
set rating = sqlc.arg(rating)::int
where id = $1;
