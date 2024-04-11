-- +migrate Up
alter table outfits add column generated boolean not null default false;

-- +migrate StatementBegin
create function is_valid_for_generation(category text) returns boolean as $$
begin
    return category in ('upper garment', 'lower garment', 'shoe', 'outerwear'); 
end
$$ language plpgsql;
-- +migrate StatementEnd

-- +migrate Down
alter table outfits drop column generated;

drop function is_valid_for_generation;
