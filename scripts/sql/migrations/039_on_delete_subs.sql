-- +migrate Up
alter table subs
drop constraint subs_user_id_fkey,
add constraint subs_user_id_fkey
   foreign key (user_id)
   references users(id)
   on delete cascade;

alter table subs
drop constraint subs_subscriber_id_fkey,
add constraint subs_subscriber_id_fkey
   foreign key (subscriber_id)
   references users(id)
   on delete cascade;

-- +migrate Down
alter table subs
drop constraint subs_user_id_fkey,
add constraint subs_user_id_fkey
   foreign key (user_id)
   references users(id);

alter table subs
drop constraint subs_subscriber_id_fkey,
add constraint subs_subscriber_id_fkey
   foreign key (subscriber_id)
   references users(id);
