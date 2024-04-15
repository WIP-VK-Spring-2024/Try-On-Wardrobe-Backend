-- +migrate Up
alter table outfits add column try_on_result_id uuid references try_on_results(id) on delete set null;

alter table outfits rename column public to privacy;

alter table users add column avatar text not null default '';

create table posts(
    id uuid primary key default gen_random_uuid(),
    outfit_id uuid unique references outfits(id) on delete cascade,
    rating int not null default 0,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table post_ratings(
    user_id uuid references users(id) on delete cascade,
    post_id uuid references posts(id) on delete cascade,
    value int not null check (value = 1 or value = -1),

    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

    primary key(post_id, user_id)
);

create table post_comments(
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    post_id uuid not null references posts(id) on delete cascade,
    body text not null,
    rating int not null default 0,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

create table post_comment_ratings(
    user_id uuid references users(id) on delete cascade,
    comment_id uuid references post_comments(id) on delete cascade,
    value int not null,

    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

    primary key(comment_id, user_id)
);

-- +migrate StatementBegin
create or replace function handle_outfit_post_link() returns trigger as $$
begin
    if tg_op = 'CREATE' and new.privacy <> 'private' then
        insert into posts(outfit_id) values (new.id);
        return new;
    elsif tg_op = 'UPDATE' and old.privacy <> new.privacy then
        if new.privacy = 'private' then 
            delete from posts where outfit_id = new.id;
        elsif old.privacy = 'private' then
            insert into posts(outfit_id) values (new.id);
        end if;
        return new;
    elsif tg_op = 'DELETE' and old.privacy <> 'private' then
        delete from posts where outfit_id = old.id;
        return old;
    end if;
end
$$ language plpgsql;
-- +migrate StatementEnd

create trigger trigger_outfit_post_link
    after delete or insert or update of privacy
    on outfits
    for each row
    execute procedure handle_outfit_post_link();

-- +migrate StatementBegin
create or replace function post_rating_count() returns trigger as $$
begin
    if tg_op = 'CREATE' then
        update posts
            set rating = rating + new.value
            where id = new.post_id;
        return new;
    elsif tg_op = 'UPDATE' then
        update posts
            set rating = rating + new.value - old.value
            where id = new.post_id;
        return new;
    elsif tg_op = 'DELETE' then
        update posts
            set rating = rating - old.value
            where id = new.post_id;
        return old;
    end if;
end
$$ language plpgsql;
-- +migrate StatementEnd

create trigger trigger_post_rating_count
    after delete or insert or update of value
    on post_ratings
    for each row
    execute procedure post_rating_count();

-- +migrate StatementBegin
create or replace function post_comment_rating_count() returns trigger as $$
begin
    if tg_op = 'CREATE' then
        update post_comments
            set rating = rating + new.value
            where id = new.comment_id;
        return new;
    elsif tg_op = 'UPDATE' then
        update post_comments
            set rating = rating + new.value - old.value
            where id = new.comment_id;
        return new;
    elsif tg_op = 'DELETE' then
        update post_comments
            set rating = rating - old.value
            where id = new.comment_id;
        return old;
    end if;
end
$$ language plpgsql;
-- +migrate StatementEnd

create trigger trigger_post_comment_rating_count
    after delete or insert or update of value
    on post_comment_ratings
    for each row
    execute procedure post_comment_rating_count();

-- +migrate Down
alter table outfits drop column try_on_result_id;
alter table outfits rename column privacy to public;
alter table users drop column avatar;

drop table posts cascade;
drop table post_ratings;
drop table post_comments cascade;
drop table post_comment_ratings;

drop function handle_outfit_post_link;
drop function post_rating_count;
drop function post_comment_rating_count;

drop trigger trigger_outfit_post_link on outfits;
