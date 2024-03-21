-- +migrate Up
insert into styles(id, name) values 
    ('59e148b8-7f67-40f4-ab9a-914c3d9ac49b', 'Повседневный'),
    ('ffa16b60-8ab8-4db7-bbff-d0ae06e76d30', 'Спортивный'),
    ('e85dc243-aa53-4b7b-b840-cb41dbd8f952', 'Официальный');

-- +migrate Down
delete from styles where id in (
    '59e148b8-7f67-40f4-ab9a-914c3d9ac49b',
    'ffa16b60-8ab8-4db7-bbff-d0ae06e76d30',
    'e85dc243-aa53-4b7b-b840-cb41dbd8f952'
);
