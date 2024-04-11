-- +migrate Up
alter table clothes alter column updated_at set default now();
alter table outfits alter column updated_at set default now();
alter table user_images alter column updated_at set default now();
alter table try_on_results alter column updated_at set default now();
alter table types alter column updated_at set default now();
alter table subtypes alter column updated_at set default now();
alter table users alter column updated_at set default now();
alter table styles alter column updated_at set default now();
alter table tags alter column updated_at set default now();
