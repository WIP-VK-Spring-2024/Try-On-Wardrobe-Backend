-- +migrate Up

drop index subtypes_name_idx;
create index subtypes_name_idx on subtypes (name varchar_pattern_ops);

alter table clothes_tags
drop constraint clothes_tags_clothes_id_fkey,
add constraint clothes_tags_clothes_id_fkey
   foreign key (clothes_id)
   references clothes(id)
   on delete cascade;

alter table clothes_tags
drop constraint clothes_tags_tag_id_fkey,
add constraint clothes_tags_tag_id_fkey
   foreign key (tag_id)
   references tags(id)
   on delete cascade;


-- +migrate Down
drop index subtypes_name_idx;
create unique index subtypes_name_idx on subtypes (name varchar_pattern_ops);

alter table clothes_tags
drop constraint clothes_tags_clothes_id_fkey,
add constraint clothes_tags_clothes_id_fkey
   foreign key (clothes_id)
   references clothes(id);

alter table clothes_tags
drop constraint clothes_tags_tag_id_fkey,
add constraint clothes_tags_tag_id_fkey
   foreign key (tag_id)
   references tags(id);
