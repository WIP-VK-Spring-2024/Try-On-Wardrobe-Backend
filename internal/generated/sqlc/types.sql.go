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

const getSubtypes = `-- name: GetSubtypes :many
select id, created_at, updated_at, name, type_id, eng_name from subtypes
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
    types.updated_at
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
