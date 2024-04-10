-- +migrate Up
create type status as enum ('active', 'wishlist', 'repair', 'give_away');

alter table clothes add column status status default 'active';

-- +migrate Down
alter table clothes drop column status;

drop type status;
