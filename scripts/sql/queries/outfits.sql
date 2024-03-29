-- name: CreateOutfit :one
insert into outfits(
    user_id,
    transforms
) values ($1, $2)
returning id;

-- name: UpdateOutfit :exec
update outfits
set name = coalesce($2, name),
    note = coalesce($3, note),
    style_id = coalesce($4, style_id),
    transforms = coalesce($5, transforms),
    seasons = coalesce(sqlc.arg(seasons), seasons)::season[],
    updated_at = now()
where id = $1;

-- name: SetOutfitImage :exec
update outfits
set image = $2
where id = $1;

-- name: GetOutfit :one
select *
from outfits
where id = $1;

-- name: GetOutfitsByUser :many
select *
from outfits
where user_id = $1;

-- name: DeleteOutfit :exec
delete from outfits
where id = $1;
