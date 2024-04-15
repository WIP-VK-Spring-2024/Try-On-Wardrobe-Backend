-- +migrate Up
alter table outfits add column try_on_result_id uuid references try_on_results(id) on delete set null;

create table posts(
    id uuid primary key default gen_random_uuid(),
    outfit_id uuid references outfits(id) on delete cascade,
    rating int not null default 0,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table post_ratings(
    user_id uuid references users(id) on delete cascade,
    post_id uuid references posts(id) on delete cascade,
    value int not null,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

    primary key(outfit_id, user_id)
);

create table post_comments(
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    post_id uuid not null references posts(id) on delete cascade,
    body text not null,
    rating int not null default 0,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table post_comment_ratings(
    user_id uuid references users(id) on delete cascade,
    comment_id uuid references post_comments(id) on delete cascade,
    value int not null,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

    primary key(post_id, user_id)
);

create or replace function handle_outfit_post_link() returns trigger as $$
begin

end
$$ language plpgsql;

-- +migrate Down
alter table outfits drop column try_on_result_id;

drop table posts cascade;
