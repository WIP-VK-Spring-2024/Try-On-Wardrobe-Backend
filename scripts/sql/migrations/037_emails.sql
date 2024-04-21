-- +migrate Up
update users
set email = 'anastasia@test.ru'
where name = 'Anastasia';

update users
set email = 'nikita@test.ru'
where name = 'Nikita';

alter table users alter column email set not null;

-- +migrate Down
alter table users alter column email set null;
