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
    outfits.public;

-- name: GetOutfitsByUser :many
select
    outfits.*,
    array_remove(array_agg(tags.name), null)::text[] as tags
from outfits
left join outfits_tags ot on ot.outfit_id = outfits.id
left join tags on tags.id = ot.tag_id 
where outfits.user_id = $1
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
    outfits.public
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
    outfits.public
order by outfits.created_at desc
limit $2;

select 
    outfits.*,
    array_remove(array_agg(tags.name), null)::text[] as tags
from outfits
left join outfits_tags ot on ot.outfit_id = outfits.id
left join tags on tags.id = ot.tag_id 
where outfits.public = true
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
    outfits.public
order by outfits.created_at desc
limit 1;
-- name: DeleteOutfit :exec
delete from outfits
where id = $1;

-- name: GetOutfitClothesInfo :many
select
    clothes.id,
    case when clothes.type = 'Верх' then 'upper_body'
         when clothes.type = 'Низ' then 'lower_body'
         when clothes.type = 'Платья' then 'dresses'
         else '' end as category
from outfits
join clothes on outfit.transforms ? clothes.id
where outfits.id = $1 and category <> '';
