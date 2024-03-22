-- name: CreateTags :exec
insert into tags (name) values (  
  unnest(sqlc.arg(names)::varchar[])
) on conflict do nothing;

-- name: CreateClothesTagLinks :exec
insert into clothes_tags (clothes_id, tag_id)
    select sqlc.arg(clothes_id), id
    from tags where name = any(sqlc.arg(tags)::text[]);

-- name: GetTags :many
select
    id, 
    name,
    use_count
from tags
order by use_count desc
limit $1 offset $2;
