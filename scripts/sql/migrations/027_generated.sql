-- +migrate Up
alter table outfits add column generated boolean not null default false;
alter table outfits add column viewed boolean;

-- +migrate Down
alter table outfits drop column generated;
alter table outfits drop column viewed;
