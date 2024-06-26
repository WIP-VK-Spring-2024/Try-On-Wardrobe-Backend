-- +migrate Up
alter table tags add column use_count int not null default 0;

create index tags_use_count_idx on tags(use_count) include (id, name);

-- +migrate StatementBegin
create or replace function handle_tag_use_count() returns trigger as $$
    begin
        if (tg_op = 'DELETE') then
            update tags set use_count = use_count - 1 where id = old.tag_id;
            return old;
        elsif (tg_op = 'INSERT') then
            update tags set use_count = use_count + 1 where id = new.tag_id;
            return new;
        end if;
        return null;
    end
$$ language plpgsql;
-- +migrate StatementEnd

create or replace trigger trigger_tag_use_count
    after insert or delete
    on clothes_tags
    for each row
    execute procedure handle_tag_use_count();

-- +migrate Down
alter table tags drop column use_count;

drop index tags_use_count_idx;

drop trigger trigger_tag_use_count; on clothes_tags;

drop function handle_tag_use_count;
