-- name: RateComment :exec
insert into post_comment_ratings(comment_id, user_id, value)
    values($1, $2, $3)
    on conflict (comment_id, user_id) do update
    set value = excluded.value;

-- name: CreateComment :one
insert into post_comments(post_id, user_id, body, path)
    values($1, $2, $3, (select path from post_comments p where p.id = sqlc.arg(parent_id)))
    returning id;

-- name: GetComment :one
select *
from post_comments
where id = $1;

-- name: DeleteComment :exec
delete from post_comments
where id = $1;

-- name: UpdateComment :exec
update post_comments
set body = $2,
    updated_at = now()
where id = $1;

-- name: GetComments :many
select
    post_comments.id,
    post_comments.created_at,
    post_comments.updated_at,
    post_comments.user_id,
    post_comments.body,
    post_comments.rating,
    users.avatar as user_image,
    users.name as user_name,
    array_length(path, 1) as level,
    coalesce(post_comment_ratings.value, 0) as user_rating,
    case when post_comments.path[1] = post_comments.id then uuid_nil()
         else post_comments.path[1]::uuid end as parent_id
from post_comments
join users on users.id = post_comments.user_id
left join post_comment_ratings on post_comment_ratings.user_id = $1
            and post_comment_ratings.comment_id = post_comments.id
where post_comments.post_id = $2
  and post_comments.created_at < sqlc.arg(since)::timestamp
order by post_comments.created_at desc
limit $3;

-- name: GetCommentsTree :many
with parents as (
    select
        p.id,
        p.created_at,
        p.updated_at,
        p.created_at as sort_key,
        p.user_id,
        p.body,
        p.rating,
        u.avatar as user_image,
        u.name as user_name,
        array_length(p.path::uuid[], 1) as level,
        coalesce(r.value, 0) as user_rating,
        p.path
    from post_comments p
    join users u on u.id = p.user_id
    left join post_comment_ratings r on r.user_id = $1
        and post_comment_ratings.comment_id = post_comments.id
    where p.post_id = $2
      and p.id = p.path[1]
      and p.created_at < sqlc.arg(since)::timestamp
    order by p.created_at desc
    limit $3
), final as (
    select
        p.id,
        p.created_at,
        p.updated_at,
        parents.created_at as sort_key,
        p.user_id,
        p.body,
        p.rating,
        u.avatar as user_image,
        u.name as user_name,
        array_length(p.path, 1) as level,
        coalesce(r.value, 0) as user_rating,
        p.path
    from post_comments p
    join users u on u.id = p.user_id
    left join post_comment_ratings r on r.user_id = $1
        and post_comment_ratings.comment_id = post_comments.id
    join parents on parents.id = p.path[1]
    where p.id != p.path[1]
    union all
    select * from parents
) select
    id,
    created_at,
    updated_at,
    user_id,
    body,
    rating,
    user_image,
    user_name,
    level,
    user_rating,
    case when path[1] = id then uuid_nil()
         else path[1]::uuid end as parent_id
  from final
  order by sort_key desc, path;
