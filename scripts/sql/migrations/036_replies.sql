-- +migrate Up
alter table post_comments add column path uuid[];

update post_comments
    set path = array[id]::uuid[]
    where path is null;

alter table post_comments alter column path set not null;

create index post_comment_path_idx on post_comments(path);

-- +migrate StatementBegin
create or replace function handle_post_comment() returns trigger as $$
begin
    if tg_op = 'INSERT' then
        new.path = array_append(new.path, new.id);
        return new;
    elsif tg_op = 'DELETE' then
        delete from post_comments
            where old.id = any(path);
        return old;
    else
        return new;
    end if;
end
$$ language plpgsql;
-- +migrate StatementEnd

create trigger trigger_insert_comment_replies
before insert
on post_comments
for each row
execute procedure handle_post_comment();

create trigger trigger_delete_comment_replies
after delete
on post_comments
for each row
execute procedure handle_post_comment();

-- +migrate Down
drop trigger trigger_insert_comment_replies on post_comments;
drop trigger trigger_delete_comment_replies on post_comments;

drop function handle_post_comment;

alter table post_comments drop column path;
