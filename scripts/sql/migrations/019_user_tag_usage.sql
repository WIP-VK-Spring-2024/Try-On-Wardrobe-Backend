-- +migrate Up
create table user_tag_usage(
    user_id uuid references users(id) on delete cascade,
    tag_id uuid references tags(id) on delete cascade,
    usage int default 0,
    primary key (user_id, tag_id)
);

-- +migrate StatementBegin
create or replace function handle_user_clothes_tag_usage()
returns trigger as $$
    declare
        tag_link_user_id uuid;
    begin
        if (tg_op = 'DELETE') then
            select user_id
            into tag_link_user_id
            from clothes
            where clothes_id = old.clothes_id;

            update user_tag_usage
            set usage = usage - 1
            where tag_id = old.tag_id
                and user_id = tag_link_user_id;
            return old;
        elsif (tg_op = 'INSERT') then
            select user_id
            into tag_link_user_id
            from clothes
            where clothes_id = new.clothes_id;

            insert into user_tag_usage(user_id, tag_id)
            values (tag_link_user_id, new.tag_id)
            on conflict do update set usage = usage + 1;
            return new;
        end if;
    end
$$ language plpgsql;
-- +migrate StatementEnd

-- +migrate StatementBegin
create or replace function handle_user_outfits_tag_usage()
returns trigger as $$
    declare
        tag_link_user_id uuid;
    begin
        if (tg_op = 'DELETE') then
            select user_id
            into tag_link_user_id
            from outfits
            where outfit_id = old.outfit_id;

            update user_tag_usage
            set usage = usage - 1
            where tag_id = old.tag_id
                and user_id = tag_link_user_id;
            return old;
        elsif (tg_op = 'INSERT') then
            select user_id
            into tag_link_user_id
            from outfits
            where outfit_id = new.outfit_id;

            insert into user_tag_usage(user_id, tag_id)
            values (tag_link_user_id, new.tag_id)
            on conflict do update set usage = usage + 1;
            return new;
        end if;
    end
$$ language plpgsql;
-- +migrate StatementEnd

create trigger trigger_user_clothes_tag_usage
    after insert or delete
    on clothes_tags
    for each row
    execute procedure handle_user_clothes_tag_usage();

create trigger trigger_user_outfits_tag_usage
    after insert or delete
    on outfits_tags
    for each row
    execute procedure handle_user_outfits_tag_usage();


-- +migrate Down
drop table user_tag_usage;

drop trigger trigger_user_clothes_tag_usage;
drop function handle_user_clothes_tag_usage;

drop trigger trigger_user_outfits_tag_usage;
drop function handle_user_outfits_tag_usage;
