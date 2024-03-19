-- +migrate Up
alter table subtypes add column type_id uuid references types(id) not null;

-- +migrate Down
alter table subtypes drop column type_id uuid;
