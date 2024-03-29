-- +migrate Up
alter table clothes add column image varchar;

update clothes
set image = 'cut/' || id::text
where image is null;

alter table clothes alter column image set not null;

alter table user_images add column image varchar;

update user_images
set image = 'people/' || id::text
where image is null;

alter table user_images alter column image set not null;

-- +migrate Down
alter table clothes drop column image;
alter table user_images drop column image;
