-- +migrate Up
alter table try_on_results drop constraint try_on_results_clothes_id_fkey;

alter table try_on_results alter column clothes_id type uuid[]
    using array[clothes_id]::uuid[];

-- +migrate Down
alter table try_on_results alter column clothes_id type uuid
    using clothes_id[1]::uuid;

alter table try_on_results add constraint try_on_results_clothes_id_fkey 
foreign key (clothes_id) REFERENCES clothes(id);
