-- name: GetClothesById :one
select
    clothes.*,
    types.name as type,
    subtypes.name as subtype,
    styles.name as style,
    array_agg(coalesce(tags.name, ''))::text[] as tags
from clothes
join types on types.id = clothes.type_id
join subtypes on subtypes.id = clothes.subtype_id
left join styles on styles.id = clothes.style_id
left join clothes_tags on clothes.id = clothes_tags.clothes_id
left join tags on clothes_tags.tag_id = tags.id
where clothes.id = $1
group by
    clothes.id,
    clothes.name,
    clothes.note,
    clothes.type_id,
    clothes.subtype_id,
    clothes.style_id,
    clothes.color,
    clothes.seasons,
    clothes.created_at,
    clothes.updated_at,
    type,
    subtype,
    style;

-- name: GetClothesByUser :many
select
    clothes.*,
    types.name as type,
    subtypes.name as subtype,
    styles.name as style,
    array_agg(coalesce(tags.name, ''))::text[] as tags
from clothes
join types on types.id = clothes.type_id
join subtypes on subtypes.id = clothes.subtype_id
left join styles on styles.id = clothes.style_id
left join clothes_tags on clothes.id = clothes_tags.clothes_id
left join tags on clothes_tags.tag_id = tags.id
where clothes.user_id = $1
group by
    clothes.id,
    clothes.name,
    clothes.note,
    clothes.type_id,
    clothes.subtype_id,
    clothes.style_id,
    clothes.color,
    clothes.seasons,
    clothes.created_at,
    clothes.updated_at,
    type,
    subtype,
    style;

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

-- name: UpdateClothes :exec
update clothes
set name = coalesce($2, name),
    note = coalesce($3, note),
    type_id = coalesce($4, type_id),
    subtype_id = coalesce($5, subtype_id),
    style_id = coalesce($6, style_id),
    color = coalesce($7, color),
    seasons = coalesce($8, seasons)::season[],
    updated_at = now()
where id = $1;
