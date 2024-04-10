-- +migrate Up
create type privacy as enum ('public', 'private', 'friends');

alter table clothes add column privacy privacy not null default 'public';
alter table users add column privacy privacy not null default 'public';

alter table outfits alter column public drop default;

alter table outfits alter column public type privacy
    using case when public = true then 'public'::privacy
         else 'private'::privacy end;  

alter table outfits alter column public set default 'public';

-- +migrate Down
alter table clothes drop column privacy;
alter table user drop column privacy;

alter table outfits alter column public type boolean default true
    using case public = 'public' when true
         else false end;  

drop type privacy;
