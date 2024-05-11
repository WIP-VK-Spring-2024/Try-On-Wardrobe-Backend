-- +migrate Up
alter table subs add constraint subs_primary_key primary key(subscriber_id, user_id);

-- +migrate Down
alter table subs drop constraint subs_primary_key;
