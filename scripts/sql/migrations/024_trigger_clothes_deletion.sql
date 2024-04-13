-- +migrate Up
alter table try_on_results
drop constraint try_on_results_user_image_id_fkey,
add constraint try_on_results_user_image_id_fkey
   foreign key (user_image_id)
   references user_images(id)
   on delete cascade;

-- +migrate StatementBegin
create function delete_try_on_results() returns trigger as $$
begin
    delete from try_on_results tr
    where old.id = any(tr.clothes_id);
end
$$ language plpgsql;
-- +migrate StatementEnd

create trigger trigger_del_try_on_with_clothes
    after delete
    on clothes
    for each row
    execute procedure delete_try_on_results();

-- +migrate StatementBegin
create function delete_outfits() returns trigger as $$
begin
    delete from outfits o
    where o.transforms ? cast(old.id as text);
end
$$ language plpgsql;
-- +migrate StatementEnd

create trigger trigger_del_outfit_with_clothes
    after delete
    on clothes
    for each row
    execute procedure delete_outfits();

-- +migrate Down
alter table try_on_results
drop constraint try_on_results_user_image_id_fkey,
add constraint try_on_results_user_image_id_fkey
   foreign key (user_image_id)
   references user_images(id);

drop trigger trigger_del_try_on_with_clothes on clothes;

drop function delete_try_on_results;

drop trigger trigger_del_outfit_with_clothes on clothes;

drop function delete_outfits;
