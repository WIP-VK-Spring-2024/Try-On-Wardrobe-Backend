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
    types.updated_at
order by types.created_at, types.name;

-- name: GetSubtypes :many
select * from subtypes;

-- -- name: GetTypeIdsByEngName :many
-- select id from types
-- where eng_name in ($1::text[]);

-- -- name: GetSubtypeIdsByEngName :many
-- select id from types
-- where eng_name in ($1::text[]);
