-- +migrate Up
alter table outfits_tags
drop constraint outfits_tags_outfit_id_fkey,
add constraint outfits_tags_outfit_id_fkey
   foreign key (outfit_id)
   references outfits(id)
   on delete cascade;

alter table outfits_tags
drop constraint outfits_tags_tag_id_fkey,
add constraint outfits_tags_tag_id_fkey
   foreign key (tag_id)
   references tags(id)
   on delete cascade;

-- +migrate Down
alter table outfits_tags
drop constraint outfits_tags_outfit_id_fkey,
add constraint outfits_tags_outfit_id_fkey
   foreign key (outfit_id)
   references outfits(id);

alter table outfits_tags
drop constraint outfits_tags_tag_id_fkey,
add constraint outfits_tags_tag_id_fkey
   foreign key (tag_id)
   references tags(id);
