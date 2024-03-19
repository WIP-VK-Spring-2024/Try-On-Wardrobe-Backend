-- +migrate Up

create table users (
    id uuid default gen_random_uuid() primary key,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    name varchar(256),
    email varchar(512),
    password varchar(256),
    gender gender default gender('unknown')
);

create table types (
    id uuid default gen_random_uuid() primary key,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    name varchar(64) not null
);

create table subtypes (
    id uuid default gen_random_uuid() primary key,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    name varchar(64) not null
);

create table styles (
    id uuid default gen_random_uuid() primary key,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    name varchar(64) not null
);

create table clothes (
    id uuid default gen_random_uuid() primary key,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    name varchar(128) not null,
    note varchar(512),
    image varchar(256),
    user_id uuid not null references users(id),
    style_id uuid default null references styles(id),
    type_id uuid not null references types(id),
    subtype_id uuid not null  references subtypes(id),
    color char(7),
    seasons season[]
);

create table tags (
    id uuid default gen_random_uuid() primary key,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    name varchar(64) not null
);

create table clothes_tags (
    clothes_id uuid references clothes(id),
    tag_id uuid references tags(id),
    primary key(clothes_id, tag_id)
);

create table user_images (
    id uuid default gen_random_uuid() primary key,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    user_id uuid not null,
    image text
);

create table try_on_results (
    id uuid default gen_random_uuid() primary key,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    image text not null,
    rating bigint,
    user_id uuid not null references users(id),
    clothes_id uuid not null references clothes(id)
);


-- +migrate Down
DROP table users;

DROP table types;

DROP table subtypes;

DROP table styles;

DROP table clothes;

DROP table tags;

DROP table clothes_tags;

DROP table user_images;

DROP table try_on_results;
