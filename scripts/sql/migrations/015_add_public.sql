-- +migrate Up
alter table outfits add column public boolean not null default true;

-- +migrate Down
alter table outfits drop column public;
