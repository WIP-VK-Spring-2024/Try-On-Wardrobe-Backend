-- name: CreateTryOnResult :one
insert into try_on_results(
    clothes_id,
    user_image_id,
    image
) values ($1, $2, $3)
returning id;

-- name: DeleteTryOnResult :exec
delete from try_on_results
where id = $1;

-- name: GetTryOnResult :one
select try_on_results.*
from try_on_results
where id = $1;

-- name: GetTryOnResultsByUser :many
select try_on_results.*
from try_on_results
join user_images u on u.id = try_on_results.user_image_id
where u.user_id = $1
order by try_on_results.created_at desc;;

-- name: GetTryOnResultByClothes :one
select *
from try_on_results
where $2 <@ clothes_id
    and user_image_id = $1
limit 1;

-- name: GetTryOnResultByOutfit :one
select try_on_results.*
from try_on_results
join outfits on outfits.try_on_result_id = try_on_results.id
where outfits.id = $2
     and user_image_id = $1
limit 1;

-- name: RateTryOnResult :exec
update try_on_results 
set rating = sqlc.arg(rating)::int,
    updated_at = now()
where id = $1;
