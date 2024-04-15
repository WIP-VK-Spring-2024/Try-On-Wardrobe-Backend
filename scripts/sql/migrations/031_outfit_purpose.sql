-- +migrate Up
create table outfit_purpose(
    id uuid primary key default gen_random_uuid(),
    name text not null,
    eng_name text not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

alter table outfits add column purpose_ids uuid[] not null default '{}'::uuid[];

-- +migrate StatementBegin
create function delete_purpose_from_outfit() returns trigger as $$
begin
    update outfits
    set purpose_ids = array_remove(purpose_ids, old.id)
    where old.id = any(purpose_ids);
end
$$ language plpgsql;
-- +migrate StatementEnd

create trigger remove_purpose_from_outfit_on_deletion
    after delete
    on outfit_purpose
    for each row
    execute procedure delete_purpose_from_outfit();

-- +migrate Down
drop trigger remove_purpose_from_outfit_on_deletion on outfit_purpose;

drop function delete_purpose_from_outfit;

alter table outfits drop column purpose_ids;

drop table outfit_purpose;
