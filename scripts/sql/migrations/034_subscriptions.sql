-- +migrate Up
create table subs(
    subscriber_id uuid not null references users(id),
    user_id uuid not null references users(id),
    created_at timestamp not null default now()
);

-- +migrate Down
drop table subs;
