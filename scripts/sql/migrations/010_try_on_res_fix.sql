-- +migrate Up
alter table try_on_results
drop constraint try_on_results_clothes_id_fkey,
add constraint try_on_results_clothes_id_fkey
   foreign key (clothes_id)
   references clothes(id)
   on delete cascade;

-- +migrate Down
alter table try_on_results
drop constraint try_on_results_clothes_id_fkey,
add constraint try_on_results_clothes_id_fkey
   foreign key (clothes_id)
   references clothes(id);
