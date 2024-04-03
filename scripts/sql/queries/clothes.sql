-- name: GetClothesById :one
select
    clothes.*,
    types.name as type,
    coalesce(types.tryonable, false) as tryonable,
    subtypes.name as subtype,
    styles.name as style,
    array_remove(array_agg(tags.name), null)::text[] as tags
from clothes
left join types on types.id = clothes.type_id
left join subtypes on subtypes.id = clothes.subtype_id
left join styles on styles.id = clothes.style_id
left join clothes_tags on clothes.id = clothes_tags.clothes_id
left join tags on clothes_tags.tag_id = tags.id
where clothes.id = $1
group by
    clothes.id,
    clothes.user_id,
    clothes.name,
    clothes.note,
    clothes.image,
    clothes.type_id,
    clothes.subtype_id,
    clothes.style_id,
    clothes.status,
    clothes.color,
    clothes.seasons,
    clothes.created_at,
    clothes.updated_at,
    tryonable,
    type,
    subtype,
    style;

-- name: GetClothesByUser :many
select
    clothes.*,
    types.name as type,
    coalesce(types.tryonable, false) as tryonable,
    subtypes.name as subtype,
    styles.name as style,
    array_remove(array_agg(tags.name), null)::text[] as tags
from clothes
left join types on types.id = clothes.type_id
left join subtypes on subtypes.id = clothes.subtype_id
left join styles on styles.id = clothes.style_id
left join clothes_tags on clothes.id = clothes_tags.clothes_id
left join tags on clothes_tags.tag_id = tags.id
where clothes.user_id = $1
group by
    clothes.id,
    clothes.user_id,
    clothes.name,
    clothes.note,
    clothes.image,
    clothes.type_id,
    clothes.subtype_id,
    clothes.style_id,
    clothes.status,
    clothes.color,
    clothes.seasons,
    clothes.created_at,
    clothes.updated_at,
    tryonable,
    type,
    subtype,
    style
order by clothes.created_at desc;

-- name: DeleteClothes :exec
delete from clothes
where id = $1;

-- name: CreateClothes :one
insert into clothes(
    name,
    user_id,
    type_id,
    subtype_id,
    color
)
values ($1, $2, $3, $4, $5)
returning id;

-- name: SetClothesImage :exec
update clothes
set image = $2
where id = $1;

-- name: UpdateClothes :exec
update clothes
set name = coalesce($2, name),
    note = coalesce($3, note),
    type_id = coalesce($4, type_id),
    subtype_id = coalesce($5, subtype_id),
    style_id = coalesce($6, style_id),
    color = coalesce($7, color),
    seasons = coalesce(sqlc.arg(seasons), seasons)::season[],
    updated_at = now()
where id = $1;

-- name: GetClothesIdByOutfit :many
select c.id
from clothes c
join outfits o on o.transforms ? c.id
where o.id = $1;
