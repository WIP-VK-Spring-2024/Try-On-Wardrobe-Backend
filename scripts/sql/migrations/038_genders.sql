-- +migrate Up
update users
set gender = 'female'
where name = 'Anastasia';

update users
set gender = 'male'
where name = 'Nikita';

alter table users alter column gender set not null;

-- +migrate Down
alter table users alter column gender set null;
