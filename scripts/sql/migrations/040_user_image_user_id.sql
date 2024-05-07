-- +migrate Up
alter table user_images
    add constraint user_images_user_id_fkey
    foreign key (user_id)
    references users(id)
    on delete cascade;

create index user_images_user_id_idx on user_images(user_id);

-- +migrate Down
drop index user_images_user_id_idx;
alter table user_images drop constraint user_images_user_id_fkey;
