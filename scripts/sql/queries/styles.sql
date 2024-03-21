-- name: CreateStyle :one
insert into styles(name, created_at)
values ($1, now())
on conflict(name) do update
set name = excluded.name
returning id;

-- name: GetStyles :many
select * from styles;
