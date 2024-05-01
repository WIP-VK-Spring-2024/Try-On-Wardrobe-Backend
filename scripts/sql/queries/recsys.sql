-- name: GetClothesTensors :many
select
    outfits.id,
    outfits.user_id,
    array_agg(cv.tensor) as clothes
from outfits
join clothes_vector cv on outfits.transforms ? cv.clothes_id::text
where outfits.privacy = 'public'
group by
    outfits.id,
    outfits.user_id;
