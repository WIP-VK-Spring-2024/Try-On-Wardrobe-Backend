-- +migrate Up
create table outfits(
    id uuid default gen_random_uuid() primary key,
    user_id uuid references users(id),
    style_id uuid references styles(id),
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    name varchar(256),
    note text,
    image varchar,
    transforms jsonb,
    seasons season[]
);

create table outfits_tags (
    outfit_id uuid references outfits(id),
    tag_id uuid references tags(id),
    primary key(outfit_id, tag_id)
);

-- +migrate Down
drop table outfits;
drop table outfits_tags;
