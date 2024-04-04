-- +migrate Up
alter table user_tag_usage alter column usage set default 1;

-- +migrate Down
alter table user_tag_usage alter column usage set default 0;
