// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: comments.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"try-on/internal/pkg/utils"
)

const createComment = `-- name: CreateComment :one
insert into post_comments(post_id, user_id, body, path)
    values($1, $2, $3, (select path from post_comments p where p.id = $4))
    returning id
`

type CreateCommentParams struct {
	PostID   utils.UUID
	UserID   utils.UUID
	Body     string
	ParentID utils.UUID
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (utils.UUID, error) {
	row := q.db.QueryRow(ctx, createComment,
		arg.PostID,
		arg.UserID,
		arg.Body,
		arg.ParentID,
	)
	var id utils.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteComment = `-- name: DeleteComment :exec
delete from post_comments
where id = $1
`

func (q *Queries) DeleteComment(ctx context.Context, id utils.UUID) error {
	_, err := q.db.Exec(ctx, deleteComment, id)
	return err
}

const getComment = `-- name: GetComment :one
select id, user_id, post_id, body, rating, created_at, updated_at, path
from post_comments
where id = $1
`

func (q *Queries) GetComment(ctx context.Context, id utils.UUID) (PostComment, error) {
	row := q.db.QueryRow(ctx, getComment, id)
	var i PostComment
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PostID,
		&i.Body,
		&i.Rating,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Path,
	)
	return i, err
}

const getComments = `-- name: GetComments :many
select
    post_comments.id,
    post_comments.created_at,
    post_comments.updated_at,
    post_comments.user_id,
    post_comments.body,
    post_comments.rating,
    users.avatar as user_image,
    users.name as user_name,
    array_length(path, 1) as level,
    coalesce(post_comment_ratings.value, 0) as user_rating,
    case when post_comments.path[1] = post_comments.id then uuid_nil()
         else post_comments.path[1]::uuid end as parent_id
from post_comments
join users on users.id = post_comments.user_id
left join post_comment_ratings on post_comment_ratings.user_id = $1
where post_comments.post_id = $2
  and post_comments.created_at < $4::timestamp
order by post_comments.created_at desc
limit $3
`

type GetCommentsParams struct {
	UserID utils.UUID
	PostID utils.UUID
	Limit  int32
	Since  pgtype.Timestamp
}

type GetCommentsRow struct {
	ID         utils.UUID
	CreatedAt  pgtype.Timestamp
	UpdatedAt  pgtype.Timestamp
	UserID     utils.UUID
	Body       string
	Rating     int32
	UserImage  string
	UserName   string
	Level      int32
	UserRating int32
	ParentID   utils.UUID
}

func (q *Queries) GetComments(ctx context.Context, arg GetCommentsParams) ([]GetCommentsRow, error) {
	rows, err := q.db.Query(ctx, getComments,
		arg.UserID,
		arg.PostID,
		arg.Limit,
		arg.Since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCommentsRow
	for rows.Next() {
		var i GetCommentsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.Body,
			&i.Rating,
			&i.UserImage,
			&i.UserName,
			&i.Level,
			&i.UserRating,
			&i.ParentID,
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

const getCommentsTree = `-- name: GetCommentsTree :many
with parents as (
    select
        p.id,
        p.created_at,
        p.updated_at,
        p.created_at as sort_key,
        p.user_id,
        p.body,
        p.rating,
        u.avatar as user_image,
        u.name as user_name,
        array_length(p.path::uuid[], 1) as level,
        coalesce(r.value, 0) as user_rating,
        p.path
    from post_comments p
    join users u on u.id = p.user_id
    left join post_comment_ratings r on r.user_id = $1
    where p.post_id = $2
      and p.id = p.path[1]
      and p.created_at < $4::timestamp
    order by p.created_at desc
    limit $3
), final as (
    select
        p.id,
        p.created_at,
        p.updated_at,
        parents.created_at as sort_key,
        p.user_id,
        p.body,
        p.rating,
        u.avatar as user_image,
        u.name as user_name,
        array_length(p.path, 1) as level,
        coalesce(r.value, 0) as user_rating,
        p.path
    from post_comments p
    join users u on u.id = p.user_id
    left join post_comment_ratings r on r.user_id = $1
    join parents on parents.id = p.path[1]
    where p.id != p.path[1]
    union all
    select id, created_at, updated_at, sort_key, user_id, body, rating, user_image, user_name, level, user_rating, path from parents
) select
    id,
    created_at,
    updated_at,
    user_id,
    body,
    rating,
    user_image,
    user_name,
    level,
    user_rating,
    case when path[1] = id then uuid_nil()
         else path[1]::uuid end as parent_id
  from final
  order by sort_key desc, path
`

type GetCommentsTreeParams struct {
	UserID utils.UUID
	PostID utils.UUID
	Limit  int32
	Since  pgtype.Timestamp
}

type GetCommentsTreeRow struct {
	ID         utils.UUID
	CreatedAt  pgtype.Timestamp
	UpdatedAt  pgtype.Timestamp
	UserID     utils.UUID
	Body       string
	Rating     int32
	UserImage  string
	UserName   string
	Level      int32
	UserRating int32
	ParentID   utils.UUID
}

func (q *Queries) GetCommentsTree(ctx context.Context, arg GetCommentsTreeParams) ([]GetCommentsTreeRow, error) {
	rows, err := q.db.Query(ctx, getCommentsTree,
		arg.UserID,
		arg.PostID,
		arg.Limit,
		arg.Since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCommentsTreeRow
	for rows.Next() {
		var i GetCommentsTreeRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.Body,
			&i.Rating,
			&i.UserImage,
			&i.UserName,
			&i.Level,
			&i.UserRating,
			&i.ParentID,
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

const rateComment = `-- name: RateComment :exec
insert into post_comment_ratings(comment_id, user_id, value)
    values($1, $2, $3)
    on conflict (comment_id, user_id) do update
    set value = excluded.value
`

type RateCommentParams struct {
	CommentID utils.UUID
	UserID    utils.UUID
	Value     int32
}

func (q *Queries) RateComment(ctx context.Context, arg RateCommentParams) error {
	_, err := q.db.Exec(ctx, rateComment, arg.CommentID, arg.UserID, arg.Value)
	return err
}

const updateComment = `-- name: UpdateComment :exec
update post_comments
set body = $2,
    updated_at = now()
where id = $1
`

func (q *Queries) UpdateComment(ctx context.Context, iD utils.UUID, body string) error {
	_, err := q.db.Exec(ctx, updateComment, iD, body)
	return err
}