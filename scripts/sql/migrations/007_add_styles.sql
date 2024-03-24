-- +migrate Up
update clothes set style_id = '59e148b8-7f67-40f4-ab9a-914c3d9ac49b'
where id in ('62e29ffe-b3dd-4652-bc18-d4aebb76068f', '0208edd3-5dcc-4543-993c-f8da2764bb03', 'd7f12bdb-83ca-45c6-adc3-9142099f2816');

update clothes set style_id = 'ffa16b60-8ab8-4db7-bbff-d0ae06e76d30'
where id in ('a5359bd0-7e92-4ce8-99b5-e4035b7881e2', 'af0a63e5-9a32-4026-9648-a9f02640aa48');

update clothes set style_id = 'e85dc243-aa53-4b7b-b840-cb41dbd8f952'
where id in ('7711b98c-1d8e-4720-b7b1-515e7147703f');

-- +migrate Down
update clothes set style_id = null
where id in ('62e29ffe-b3dd-4652-bc18-d4aebb76068f', '0208edd3-5dcc-4543-993c-f8da2764bb03', 'd7f12bdb-83ca-45c6-adc3-9142099f2816',
             'a5359bd0-7e92-4ce8-99b5-e4035b7881e2', 'af0a63e5-9a32-4026-9648-a9f02640aa48', '7711b98c-1d8e-4720-b7b1-515e7147703f'
);
