-- +migrate Up
create table clothes_vector(
    id uuid primary key default gen_random_uuid(),
    clothes_id uuid not null references clothes(id) on delete cascade,
    tensor bytea not null
);

create unique index clothes_vector_clothes_id_idx on clothes_vector(clothes_id);

-- +migrate Down
drop table clothes_vector;
