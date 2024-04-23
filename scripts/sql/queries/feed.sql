-- name: GetPosts :many
select
    posts.id,
    posts.created_at,
    posts.updated_at,
    posts.outfit_id,
    outfits.user_id,
    outfits.image as outfit_image,
    users.avatar as user_image,
    users.name as user_name,
    posts.rating,
    case when subs.created_at is not null then true
         else false end is_subbed,
    coalesce(post_ratings.value, 0) as user_rating,
    coalesce(try_on_results.image, '') as try_on_image,
    coalesce(try_on_results.id, uuid_nil()) as try_on_id
from posts
join outfits on outfits.id = posts.outfit_id
join users on users.id = outfits.user_id
left join post_ratings on post_ratings.user_id = $1
        and post_ratings.post_id = posts.id
left join subs on subs.subscriber_id = $1 and subs.user_id = outfits.user_id
left join try_on_results on try_on_results.id = outfits.try_on_result_id
where posts.created_at < sqlc.arg(since)::timestamp
order by posts.created_at desc
limit $2;

-- name: GetPostsByUser :many
select
    posts.id,
    posts.created_at,
    posts.updated_at,
    posts.outfit_id,
    outfits.user_id,
    outfits.image as outfit_image,
    users.avatar as user_image,
    users.name as user_name,
    posts.rating,
    case when subs.created_at is not null then true
         else false end is_subbed,
    coalesce(post_ratings.value, 0) as user_rating,
    coalesce(try_on_results.image, '') as try_on_image,
    coalesce(try_on_results.id, uuid_nil()) as try_on_id
from posts
join outfits on outfits.id = posts.outfit_id
join users on users.id = outfits.user_id
left join post_ratings on post_ratings.user_id = $1
        and post_ratings.post_id = posts.id
left join subs on subs.subscriber_id = $1 and subs.user_id = outfits.user_id
left join try_on_results on try_on_results.id = outfits.try_on_result_id
where posts.created_at < sqlc.arg(since)::timestamp
      and outfits.user_id = sqlc.arg(author_id)
order by posts.created_at desc
limit $2;

-- name: GetLikedPosts :many
select
    posts.id,
    posts.created_at,
    posts.updated_at,
    posts.outfit_id,
    outfits.user_id,
    outfits.image as outfit_image,
    users.avatar as user_image,
    users.name as user_name,
    posts.rating,
    case when subs.created_at is not null then true
         else false end is_subbed,
    post_ratings.value as user_rating,
    coalesce(try_on_results.image, '') as try_on_image,
    coalesce(try_on_results.id, uuid_nil()) as try_on_id
from posts
join outfits on outfits.id = posts.outfit_id
join users on users.id = outfits.user_id
join post_ratings on post_ratings.user_id = $1
        and post_ratings.post_id = posts.id
left join subs on subs.subscriber_id = $1 and subs.user_id = outfits.user_id
left join try_on_results on try_on_results.id = outfits.try_on_result_id
where posts.created_at < sqlc.arg(since)::timestamp
    and post_ratings.value = 1
order by posts.created_at desc
limit $2;

-- name: GetSubscriptionPosts :many
select
    posts.id,
    posts.created_at,
    posts.updated_at,
    posts.outfit_id,
    outfits.user_id,
    outfits.image as outfit_image,
    users.avatar as user_image,
    users.name as user_name,
    posts.rating,
    true as is_subbed,
    coalesce(post_ratings.value, 0) as user_rating,
    coalesce(try_on_results.image, '') as try_on_image,
    coalesce(try_on_results.id, uuid_nil()) as try_on_id
from posts
join outfits on outfits.id = posts.outfit_id
join users on users.id = outfits.user_id
join subs on subs.user_id = users.id and subs.subscriber_id = $1
left join post_ratings on post_ratings.user_id = $1
        and post_ratings.post_id = posts.id
left join try_on_results on try_on_results.id = outfits.try_on_result_id
where posts.created_at < sqlc.arg(since)::timestamp
order by posts.created_at
limit $2;

-- name: RatePost :exec
insert into post_ratings(post_id, user_id, value)
    values($1, $2, $3)
    on conflict (post_id, user_id) do update
    set value = excluded.value;

-- name: Subscribe :exec
insert into subs(subscriber_id, user_id)
    values($1, $2);

-- name: Unsubscribe :exec
delete from subs
where subscriber_id = $1 and user_id = $2;
