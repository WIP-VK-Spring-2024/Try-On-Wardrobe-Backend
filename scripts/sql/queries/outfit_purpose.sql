-- name: GetOutfitPurposeEngNames :many
select eng_name
from outfit_purpose
where name = any(sqlc.arg(eng_names)::text[]);

-- name: GetOutfitPurposeByEngName :many
select *
from outfit_purpose
where eng_name = any(sqlc.arg(eng_names)::text[]);

-- name: GetOutfitPurposes :many
select *
from outfit_purpose;
