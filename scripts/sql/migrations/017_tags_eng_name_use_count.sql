-- +migrate Up
alter table tags add column eng_name text;

create or replace trigger trigger_outfit_tag_use_count
    after insert or delete
    on outfits_tags
    for each row
    execute procedure handle_tag_use_count();

-- +migrate Down
alter table tags drop column eng_name;

drop trigger trigger_outfit_tag_use_count on outfits_tags;
