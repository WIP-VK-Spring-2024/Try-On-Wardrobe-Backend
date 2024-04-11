-- name: GetTypes :many
select
    types.*,
    array_agg(subtypes.id order by subtypes.name)::uuid[] as subtype_ids,
    array_agg(subtypes.name order by subtypes.name)::text[] as subtype_names,
    array_agg(subtypes.created_at order by subtypes.name)::timestamp[] as subtypes_created_at
from types
left join subtypes on types.id = subtypes.type_id
group by
    types.id,
    types.name,
    types.created_at,
    types.updated_at,
    types.tryonable
order by types.created_at, types.name;

-- name: GetSubtypes :many
select * from subtypes;

-- name: GetTypeIdByEngName :one
select id from types
where eng_name = $1
limit 1;

-- name: GetSubtypeIdsByEngName :one
select id from subtypes
where eng_name = $1
limit 1;

-- name: GetTypeEngNames :many
select eng_name
from types;

-- name: GetSubtypeEngNames :many
select eng_name
from subtypes;

-- name: GetTypeBySubtype :one
select t.id, t.tryonable
from types t
join subtypes s on s.type_id = t.id
where s.id = $1;
