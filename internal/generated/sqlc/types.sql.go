// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: types.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"try-on/internal/pkg/utils"
)

const getSubtypeEngNames = `-- name: GetSubtypeEngNames :many
select eng_name
from subtypes
`

func (q *Queries) GetSubtypeEngNames(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, getSubtypeEngNames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var eng_name string
		if err := rows.Scan(&eng_name); err != nil {
			return nil, err
		}
		items = append(items, eng_name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSubtypeIdsByEngName = `-- name: GetSubtypeIdsByEngName :one
select id from subtypes
where eng_name = $1
limit 1
`

func (q *Queries) GetSubtypeIdsByEngName(ctx context.Context, engName string) (utils.UUID, error) {
	row := q.db.QueryRow(ctx, getSubtypeIdsByEngName, engName)
	var id utils.UUID
	err := row.Scan(&id)
	return id, err
}

const getSubtypes = `-- name: GetSubtypes :many
select id, created_at, updated_at, name, type_id, eng_name, layer, temp_range from subtypes
`

func (q *Queries) GetSubtypes(ctx context.Context) ([]Subtype, error) {
	rows, err := q.db.Query(ctx, getSubtypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Subtype
	for rows.Next() {
		var i Subtype
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.TypeID,
			&i.EngName,
			&i.Layer,
			&i.TempRange,
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

const getTypeBySubtype = `-- name: GetTypeBySubtype :one
select t.id, t.tryonable
from types t
join subtypes s on s.type_id = t.id
where s.id = $1
`

type GetTypeBySubtypeRow struct {
	ID        utils.UUID
	Tryonable bool
}

func (q *Queries) GetTypeBySubtype(ctx context.Context, id utils.UUID) (GetTypeBySubtypeRow, error) {
	row := q.db.QueryRow(ctx, getTypeBySubtype, id)
	var i GetTypeBySubtypeRow
	err := row.Scan(&i.ID, &i.Tryonable)
	return i, err
}

const getTypeEngNames = `-- name: GetTypeEngNames :many
select eng_name
from types
`

func (q *Queries) GetTypeEngNames(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, getTypeEngNames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var eng_name string
		if err := rows.Scan(&eng_name); err != nil {
			return nil, err
		}
		items = append(items, eng_name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTypeIdByEngName = `-- name: GetTypeIdByEngName :one
select id from types
where eng_name = $1
limit 1
`

func (q *Queries) GetTypeIdByEngName(ctx context.Context, engName string) (utils.UUID, error) {
	row := q.db.QueryRow(ctx, getTypeIdByEngName, engName)
	var id utils.UUID
	err := row.Scan(&id)
	return id, err
}

const getTypes = `-- name: GetTypes :many
select
    types.id, types.created_at, types.updated_at, types.name, types.tryonable, types.eng_name,
    array_agg(subtypes.id order by subtypes.name)::uuid[] as subtype_ids,
    array_agg(subtypes.name order by subtypes.name)::text[] as subtype_names,
    array_agg(subtypes.created_at order by subtypes.name)::timestamp[] as subtypes_created_at
from types
left join subtypes on types.id = subtypes.type_id
group by
    types.id,
    types.name,
    types.created_at,
    types.updated_at,
    types.tryonable
order by types.created_at, types.name
`

type GetTypesRow struct {
	ID                utils.UUID
	CreatedAt         pgtype.Timestamptz
	UpdatedAt         pgtype.Timestamptz
	Name              string
	Tryonable         bool
	EngName           string
	SubtypeIds        []utils.UUID
	SubtypeNames      []string
	SubtypesCreatedAt []pgtype.Timestamp
}

func (q *Queries) GetTypes(ctx context.Context) ([]GetTypesRow, error) {
	rows, err := q.db.Query(ctx, getTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTypesRow
	for rows.Next() {
		var i GetTypesRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Tryonable,
			&i.EngName,
			&i.SubtypeIds,
			&i.SubtypeNames,
			&i.SubtypesCreatedAt,
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
