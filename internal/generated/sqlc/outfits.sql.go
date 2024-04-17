// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: outfits.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"
)

const createOutfit = `-- name: CreateOutfit :one
insert into outfits(
    user_id,
    transforms,
    privacy
)
select $1, $2, users.privacy
from users where users.id = $1
returning id, created_at, updated_at
`

type CreateOutfitRow struct {
	ID        utils.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

func (q *Queries) CreateOutfit(ctx context.Context, userID utils.UUID, transforms []byte) (CreateOutfitRow, error) {
	row := q.db.QueryRow(ctx, createOutfit, userID, transforms)
	var i CreateOutfitRow
	err := row.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const deleteOutfit = `-- name: DeleteOutfit :exec
delete from outfits
where id = $1
`

func (q *Queries) DeleteOutfit(ctx context.Context, id utils.UUID) error {
	_, err := q.db.Exec(ctx, deleteOutfit, id)
	return err
}

const getNotViewedOutfits = `-- name: GetNotViewedOutfits :many
select
    outfits.id, outfits.user_id, outfits.style_id, outfits.created_at, outfits.updated_at, outfits.name, outfits.note, outfits.image, outfits.transforms, outfits.seasons, outfits.privacy, outfits.generated, outfits.purpose_ids, outfits.try_on_result_id
from outfits
join users on users.id = outfits.user_id
where users.id = $1
    and outfits.viewed = false
`

func (q *Queries) GetNotViewedOutfits(ctx context.Context, id utils.UUID) ([]Outfit, error) {
	rows, err := q.db.Query(ctx, getNotViewedOutfits, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Outfit
	for rows.Next() {
		var i Outfit
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.StyleID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Note,
			&i.Image,
			&i.Transforms,
			&i.Seasons,
			&i.Privacy,
			&i.Generated,
			&i.PurposeIds,
			&i.TryOnResultID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOutfit = `-- name: GetOutfit :one
select
    outfits.id, outfits.user_id, outfits.style_id, outfits.created_at, outfits.updated_at, outfits.name, outfits.note, outfits.image, outfits.transforms, outfits.seasons, outfits.privacy, outfits.generated, outfits.purpose_ids, outfits.try_on_result_id,
    array_remove(array_agg(tags.name), null)::text[] as tags
from outfits
left join outfits_tags ot on ot.outfit_id = outfits.id
left join tags on tags.id = ot.tag_id
where outfits.id = $1
group by
    outfits.id,
    outfits.user_id,
    outfits.style_id,
    outfits.created_at,
    outfits.updated_at,
    outfits.name,
    outfits.note,
    outfits.image,
    outfits.transforms,
    outfits.seasons,
    outfits.privacy
`

type GetOutfitRow struct {
	ID            utils.UUID
	UserID        utils.UUID
	StyleID       utils.UUID
	CreatedAt     pgtype.Timestamptz
	UpdatedAt     pgtype.Timestamptz
	Name          pgtype.Text
	Note          pgtype.Text
	Image         pgtype.Text
	Transforms    []byte
	Seasons       []domain.Season
	Privacy       domain.Privacy
	Generated     bool
	PurposeIds    []utils.UUID
	TryOnResultID utils.UUID
	Tags          []string
}

func (q *Queries) GetOutfit(ctx context.Context, id utils.UUID) (GetOutfitRow, error) {
	row := q.db.QueryRow(ctx, getOutfit, id)
	var i GetOutfitRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.StyleID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Note,
		&i.Image,
		&i.Transforms,
		&i.Seasons,
		&i.Privacy,
		&i.Generated,
		&i.PurposeIds,
		&i.TryOnResultID,
		&i.Tags,
	)
	return i, err
}

const getOutfitClothesInfo = `-- name: GetOutfitClothesInfo :many
select
    clothes.id,
    try_on_type(types.name) as category
from outfits
join clothes on outfit.transforms ? clothes.id
join types on types.id = clothes.type_id
where outfits.id = $1 and try_on_type(types.name) <> ''
`

type GetOutfitClothesInfoRow struct {
	ID       utils.UUID
	Category string
}

func (q *Queries) GetOutfitClothesInfo(ctx context.Context, id utils.UUID) ([]GetOutfitClothesInfoRow, error) {
	rows, err := q.db.Query(ctx, getOutfitClothesInfo, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOutfitClothesInfoRow
	for rows.Next() {
		var i GetOutfitClothesInfoRow
		if err := rows.Scan(&i.ID, &i.Category); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOutfits = `-- name: GetOutfits :many
select 
    outfits.id, outfits.user_id, outfits.style_id, outfits.created_at, outfits.updated_at, outfits.name, outfits.note, outfits.image, outfits.transforms, outfits.seasons, outfits.privacy, outfits.generated, outfits.purpose_ids, outfits.try_on_result_id,
    array_remove(array_agg(tags.name), null)::text[] as tags
from outfits
left join outfits_tags ot on ot.outfit_id = outfits.id
left join tags on tags.id = ot.tag_id 
where outfits.public = true and outfits.created_at < $1
group by
    outfits.id,
    outfits.user_id,
    outfits.style_id,
    outfits.created_at,
    outfits.updated_at,
    outfits.name,
    outfits.note,
    outfits.image,
    outfits.transforms,
    outfits.seasons,
    outfits.privacy
order by outfits.created_at desc
limit $2
`

type GetOutfitsRow struct {
	ID            utils.UUID
	UserID        utils.UUID
	StyleID       utils.UUID
	CreatedAt     pgtype.Timestamptz
	UpdatedAt     pgtype.Timestamptz
	Name          pgtype.Text
	Note          pgtype.Text
	Image         pgtype.Text
	Transforms    []byte
	Seasons       []domain.Season
	Privacy       domain.Privacy
	Generated     bool
	PurposeIds    []utils.UUID
	TryOnResultID utils.UUID
	Tags          []string
}

func (q *Queries) GetOutfits(ctx context.Context, createdAt pgtype.Timestamptz, limit int32) ([]GetOutfitsRow, error) {
	rows, err := q.db.Query(ctx, getOutfits, createdAt, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOutfitsRow
	for rows.Next() {
		var i GetOutfitsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.StyleID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Note,
			&i.Image,
			&i.Transforms,
			&i.Seasons,
			&i.Privacy,
			&i.Generated,
			&i.PurposeIds,
			&i.TryOnResultID,
			&i.Tags,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOutfitsByUser = `-- name: GetOutfitsByUser :many
select
    outfits.id, outfits.user_id, outfits.style_id, outfits.created_at, outfits.updated_at, outfits.name, outfits.note, outfits.image, outfits.transforms, outfits.seasons, outfits.privacy, outfits.generated, outfits.purpose_ids, outfits.try_on_result_id,
    array_remove(array_agg(tags.name), null)::text[] as tags
from outfits
left join outfits_tags ot on ot.outfit_id = outfits.id
left join tags on tags.id = ot.tag_id 
where outfits.user_id = $1
group by
    outfits.id,
    outfits.user_id,
    outfits.style_id,
    outfits.created_at,
    outfits.updated_at,
    outfits.name,
    outfits.note,
    outfits.image,
    outfits.transforms,
    outfits.seasons,
    outfits.privacy
order by outfits.created_at desc
`

type GetOutfitsByUserRow struct {
	ID            utils.UUID
	UserID        utils.UUID
	StyleID       utils.UUID
	CreatedAt     pgtype.Timestamptz
	UpdatedAt     pgtype.Timestamptz
	Name          pgtype.Text
	Note          pgtype.Text
	Image         pgtype.Text
	Transforms    []byte
	Seasons       []domain.Season
	Privacy       domain.Privacy
	Generated     bool
	PurposeIds    []utils.UUID
	TryOnResultID utils.UUID
	Tags          []string
}

func (q *Queries) GetOutfitsByUser(ctx context.Context, userID utils.UUID) ([]GetOutfitsByUserRow, error) {
	rows, err := q.db.Query(ctx, getOutfitsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOutfitsByUserRow
	for rows.Next() {
		var i GetOutfitsByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.StyleID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Note,
			&i.Image,
			&i.Transforms,
			&i.Seasons,
			&i.Privacy,
			&i.Generated,
			&i.PurposeIds,
			&i.TryOnResultID,
			&i.Tags,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setOutfitImage = `-- name: SetOutfitImage :exec
update outfits
set image = $2::text,
    updated_at = now()
where id = $1
`

func (q *Queries) SetOutfitImage(ctx context.Context, iD utils.UUID, column2 string) error {
	_, err := q.db.Exec(ctx, setOutfitImage, iD, column2)
	return err
}

const updateOutfit = `-- name: UpdateOutfit :one
update outfits
set name = coalesce($2, name),
    note = coalesce($3, note),
    style_id = coalesce($4, style_id),
    transforms = coalesce($5, transforms),
    seasons = coalesce($6, seasons)::season[],
    privacy = coalesce($7::privacy, privacy),
    updated_at = now()
where id = $1
returning created_at, updated_at
`

type UpdateOutfitParams struct {
	ID         utils.UUID
	Name       pgtype.Text
	Note       pgtype.Text
	StyleID    utils.UUID
	Transforms []byte
	Seasons    []domain.Season
	Privacy    NullPrivacy
}

type UpdateOutfitRow struct {
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

func (q *Queries) UpdateOutfit(ctx context.Context, arg UpdateOutfitParams) (UpdateOutfitRow, error) {
	row := q.db.QueryRow(ctx, updateOutfit,
		arg.ID,
		arg.Name,
		arg.Note,
		arg.StyleID,
		arg.Transforms,
		arg.Seasons,
		arg.Privacy,
	)
	var i UpdateOutfitRow
	err := row.Scan(&i.CreatedAt, &i.UpdatedAt)
	return i, err
}
