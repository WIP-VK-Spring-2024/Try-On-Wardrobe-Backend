-- name: GetTypes :many
select
    types.*,
    array_agg(subtypes.id)::uuid[] as subtype_ids,
    array_agg(subtypes.name)::text[] as subtype_names
from types
left join subtypes on types.id = subtypes.type_id
group by
    types.id,
    types.name,
    types.created_at,
    types.updated_at;

-- name: GetSubtypes :many
select * from subtypes;
