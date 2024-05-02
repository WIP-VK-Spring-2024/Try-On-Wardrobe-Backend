-- name: CreateOutfit :one
insert into outfits(
    user_id,
    transforms,
    privacy
)
select $1, $2, users.privacy
from users where users.id = $1
returning id, created_at, updated_at;

-- name: UpdateOutfit :one
update outfits
set name = coalesce($2, name),
    note = coalesce($3, note),
    style_id = coalesce($4, style_id),
    transforms = coalesce($5, transforms),
    seasons = coalesce(sqlc.arg(seasons), seasons)::season[],
    privacy = coalesce(sqlc.narg(privacy)::privacy, privacy),
    updated_at = now()
where id = $1
returning created_at, updated_at;

-- name: SetOutfitTryOnResult :exec
update outfits
set try_on_result_id = $2
where id = $1;

-- name: SetOutfitImage :exec
update outfits
set image = $2::text,
    updated_at = now()
where id = $1;

-- name: GetOutfit :one
select
    outfits.*,
    array_remove(array_agg(tags.name), null)::text[] as tags
from outfits
left join outfits_tags ot on ot.outfit_id = outfits.id
left join tags on tags.id = ot.tag_id
where outfits.id = $1
group by
    outfits.id,
    outfits.user_id,
    outfits.style_id,
    outfits.created_at,
    outfits.updated_at,
    outfits.name,
    outfits.note,
    outfits.image,
    outfits.transforms,
    outfits.seasons,
    outfits.privacy;

-- name: GetOutfitsByUser :many
select
    outfits.*,
    array_remove(array_agg(tags.name), null)::text[] as tags
from outfits
left join outfits_tags ot on ot.outfit_id = outfits.id
left join tags on tags.id = ot.tag_id 
where outfits.user_id = $1
    and (sqlc.arg(public_only)::boolean = false
        or outfits.privacy = 'public')
group by
    outfits.id,
    outfits.user_id,
    outfits.style_id,
    outfits.created_at,
    outfits.updated_at,
    outfits.name,
    outfits.note,
    outfits.image,
    outfits.transforms,
    outfits.seasons,
    outfits.privacy
order by outfits.created_at desc;

-- name: GetOutfits :many
select 
    outfits.*,
    array_remove(array_agg(tags.name), null)::text[] as tags
from outfits
left join outfits_tags ot on ot.outfit_id = outfits.id
left join tags on tags.id = ot.tag_id
where outfits.public = true and outfits.created_at < $1
group by
    outfits.id,
    outfits.user_id,
    outfits.style_id,
    outfits.created_at,
    outfits.updated_at,
    outfits.name,
    outfits.note,
    outfits.image,
    outfits.transforms,
    outfits.seasons,
    outfits.privacy
order by outfits.created_at desc
limit $2;

-- name: DeleteOutfit :exec
delete from outfits
where id = $1;

-- name: GetOutfitClothesInfo :many
select
    clothes.id,
    try_on_type(types.name) as category
from outfits
join clothes on outfits.transforms ? clothes.id
join types on types.id = clothes.type_id
where outfits.id = $1 and try_on_type(types.name) <> '';
