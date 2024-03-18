<<<<<<< Updated upstream:scripts/sql/schema.sql
CREATE TABLE "users" (
    "id" uuid DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "name" varchar(256),
    "email" varchar(512),
    "password" varchar(256),
    "gender" gender DEFAULT gender('unknown'),
    PRIMARY KEY ("id")
=======
-- +migrate Up

create table users (
    id uuid default gen_random_uuid() primary key,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    name varchar(256),
    email varchar(512),
    password varchar(256),
    gender gender default gender('unknown')
>>>>>>> Stashed changes:scripts/sql/migrations/001_schema.sql
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

<<<<<<< Updated upstream:scripts/sql/schema.sql
CREATE TABLE "clothes" (
    "id" uuid DEFAULT gen_random_uuid(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "name" varchar(128),
    "note" varchar(512),
    "image" varchar(256),
    "user_id" uuid,
    "style_id" uuid DEFAULT null,
    "type_id" uuid DEFAULT null,
    "subtype_id" uuid DEFAULT null,
    "color" char(7),
    "seasons" season [],
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_clothes_user" FOREIGN KEY ("user_id") REFERENCES "users"("id"),
    CONSTRAINT "fk_clothes_style" FOREIGN KEY ("style_id") REFERENCES "styles"("id"),
    CONSTRAINT "fk_clothes_type" FOREIGN KEY ("type_id") REFERENCES "types"("id"),
    CONSTRAINT "fk_clothes_subtype" FOREIGN KEY ("subtype_id") REFERENCES "subtypes"("id")
=======
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
>>>>>>> Stashed changes:scripts/sql/migrations/001_schema.sql
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
<<<<<<< Updated upstream:scripts/sql/schema.sql
=======

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
>>>>>>> Stashed changes:scripts/sql/migrations/001_schema.sql
