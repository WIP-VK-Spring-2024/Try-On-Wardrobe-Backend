// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: styles.sql

package sqlc

import (
	"context"

	"try-on/internal/pkg/utils"
)

const createStyle = `-- name: CreateStyle :one
insert into styles(name, created_at)
values ($1, now())
on conflict(name) do update
set name = excluded.name
returning id
`

func (q *Queries) CreateStyle(ctx context.Context, name string) (utils.UUID, error) {
	row := q.db.QueryRow(ctx, createStyle, name)
	var id utils.UUID
	err := row.Scan(&id)
	return id, err
}

const getStyleEngNames = `-- name: GetStyleEngNames :many
select eng_name
from styles
`

func (q *Queries) GetStyleEngNames(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, getStyleEngNames)
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

const getStyleIdByEngName = `-- name: GetStyleIdByEngName :one
select id from styles
where eng_name = $1
limit 1
`

func (q *Queries) GetStyleIdByEngName(ctx context.Context, engName string) (utils.UUID, error) {
	row := q.db.QueryRow(ctx, getStyleIdByEngName, engName)
	var id utils.UUID
	err := row.Scan(&id)
	return id, err
}

const getStyles = `-- name: GetStyles :many
select id, created_at, updated_at, name, eng_name from styles
`

func (q *Queries) GetStyles(ctx context.Context) ([]Style, error) {
	rows, err := q.db.Query(ctx, getStyles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Style
	for rows.Next() {
		var i Style
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
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
