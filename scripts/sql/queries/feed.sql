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
    post_ratings.user_id is not null as liked
from posts
join outfits on outfits.id = posts.outfit_id
join users on users.id = outfits.user_id
left join post_ratings on post_ratings.user_id = $1
where posts.created_at < $2
order by posts.created_at
limit $3;

-- name: GetComments :many
select
    post_comments.id,
    post_comments.created_at,
    post_comments.updated_at,
    post_comments.user_id,
    post_comments.body,
    post_comments.rating,
    users.avatar as user_image,
    post_ratings.user_id is not null as liked
from post_comments
join users on users.id = post_comments.user_id
left join post_comment_ratings on post_comment_ratings.user_id = $1
where post_comments.post_id = $2
  and post_comments.created_at < $3
order by post_comments.created_at
limit $4;
