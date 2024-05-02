-- name: GetClothesTensors :many
select
    outfits.id as outfit_id,
    outfits.user_id,
    array_agg(cv.clothes_id)::uuid[] as clothes_id,
    array_agg(cv.tensor)::bytea[] as clothes_tensor
from outfits
join clothes_vector cv on outfits.transforms ? cv.clothes_id::text
where outfits.privacy = 'public'
group by
    outfits.id,
    outfits.user_id;
