-- name: GetPosts :many
select
    posts.id,
    posts.created_at,
    posts.updated_at,
    posts.outfit_id,
    outfits.user_id,
    outfits.image as outfit_image,
    users.avatar as user_image,
    posts.rating,
    case when post_ratings.user_id is not null then true
        else false end as liked
from posts
join outfits on outfits.id = posts.outfit_id
join users on users.id = outfits.user_id
left join post_ratings on post_ratings.user_id = $1
where posts.created_at < sqlc.arg(since)::timestamp
order by posts.created_at
limit $2;

-- name: GetComments :many
select
    post_comments.id,
    post_comments.created_at,
    post_comments.updated_at,
    post_comments.user_id,
    post_comments.body,
    post_comments.rating,
    users.avatar as user_image,
    case when post_comment_ratings.user_id is not null then true
        else false end as liked
from post_comments
join users on users.id = post_comments.user_id
left join post_comment_ratings on post_comment_ratings.user_id = $1
where post_comments.post_id = $2
  and post_comments.created_at < sqlc.arg(since)::timestamp
order by post_comments.created_at
limit $3;

-- name: RatePost :exec
insert into post_ratings(post_id, user_id, value)
    values($1, $2, $3)
    on conflict (post_id, user_id) do update
    set value = excluded.value;

-- name: RateComment :exec
insert into post_comment_ratings(comment_id, user_id, value)
    values($1, $2, $3)
    on conflict (comment_id, user_id) do update
    set value = excluded.value;

-- name: CreateComment :one
insert into post_comments(post_id, user_id, body)
    values($1, $2, $3)
    returning id;
