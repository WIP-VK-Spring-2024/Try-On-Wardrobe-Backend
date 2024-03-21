-- +migrate Up

insert into users(id, name, password) values
    ('2a78df8a-0277-4c72-a2d9-43fb8fef1d2c', 'Nikita', 'JggwXh8/045R53WX9SFqOQdLp5119faUJQEhbfEmNxs=:66c9a22dfb')
    on conflict do nothing;

insert into types(id, name) values
    ('a4981358-9e59-45db-8ff4-ea8c11dfee66', 'Верх'),
    ('99bfce00-b014-4458-9e26-a4675f72e352', 'Низ'),
    ('b2e705a6-5e35-4957-93d4-fd801beac077', 'Обувь')
    on conflict do nothing;

insert into subtypes(id, type_id, name) values
    ('37770f2e-b58c-42d1-a114-3b989ab96b3c', 'b2e705a6-5e35-4957-93d4-fd801beac077', 'Кроссовки'),
    ('de84c658-b79e-4ace-a025-1c53786a3e1f', 'a4981358-9e59-45db-8ff4-ea8c11dfee66', 'Рубашки'),
    ('a9c3e01d-5ef4-46ea-b874-b9ce0d778ebb', 'a4981358-9e59-45db-8ff4-ea8c11dfee66', 'Кофты'),
    ('497dd9e7-6f3f-43e8-aa44-41716cb6d39c', '99bfce00-b014-4458-9e26-a4675f72e352', 'Штаны'),
    ('d7f4684f-de41-475a-862b-4148c3849f41', 'b2e705a6-5e35-4957-93d4-fd801beac077', 'Полусапоги')
    on conflict do nothing;

insert into clothes(id, name, type_id, subtype_id, user_id) values
    ('0208edd3-5dcc-4543-993c-f8da2764bb03', 'Полусапоги', 'b2e705a6-5e35-4957-93d4-fd801beac077', 'd7f4684f-de41-475a-862b-4148c3849f41', '2a78df8a-0277-4c72-a2d9-43fb8fef1d2c'),
    ('7711b98c-1d8e-4720-b7b1-515e7147703f', 'Рубашка', 'a4981358-9e59-45db-8ff4-ea8c11dfee66', 'de84c658-b79e-4ace-a025-1c53786a3e1f', '2a78df8a-0277-4c72-a2d9-43fb8fef1d2c'),
    ('af0a63e5-9a32-4026-9648-a9f02640aa48', 'Кроссовки', 'b2e705a6-5e35-4957-93d4-fd801beac077', '37770f2e-b58c-42d1-a114-3b989ab96b3c', '2a78df8a-0277-4c72-a2d9-43fb8fef1d2c'),
    ('62e29ffe-b3dd-4652-bc18-d4aebb76068f', 'Рубашка', 'a4981358-9e59-45db-8ff4-ea8c11dfee66', 'de84c658-b79e-4ace-a025-1c53786a3e1f', '2a78df8a-0277-4c72-a2d9-43fb8fef1d2c'),
    ('d7f12bdb-83ca-45c6-adc3-9142099f2816', 'Штаны', '99bfce00-b014-4458-9e26-a4675f72e352', '497dd9e7-6f3f-43e8-aa44-41716cb6d39c', '2a78df8a-0277-4c72-a2d9-43fb8fef1d2c'),
    ('a5359bd0-7e92-4ce8-99b5-e4035b7881e2', 'Штаны', '99bfce00-b014-4458-9e26-a4675f72e352', '497dd9e7-6f3f-43e8-aa44-41716cb6d39c', '2a78df8a-0277-4c72-a2d9-43fb8fef1d2c')
    on conflict do nothing;

insert into user_images(id, user_id) values
    ('add893fb-d24d-4ee0-9976-03a6875cde8d', '2a78df8a-0277-4c72-a2d9-43fb8fef1d2c')
    on conflict do nothing;

-- +migrate Down
delete from users where id in ('2a78df8a-0277-4c72-a2d9-43fb8fef1d2c');

delete from types where id in ('a4981358-9e59-45db-8ff4-ea8c11dfee66',
                            '99bfce00-b014-4458-9e26-a4675f72e352',
                            'b2e705a6-5e35-4957-93d4-fd801beac077');

delete from subtypes where id in ('37770f2e-b58c-42d1-a114-3b989ab96b3c', 
                                'de84c658-b79e-4ace-a025-1c53786a3e1f', 
                                'a9c3e01d-5ef4-46ea-b874-b9ce0d778ebb', 
                                '497dd9e7-6f3f-43e8-aa44-41716cb6d39c', 
                                'd7f4684f-de41-475a-862b-4148c3849f41');

delete from clothes where id in  ('0208edd3-5dcc-4543-993c-f8da2764bb03',
                                'af0a63e5-9a32-4026-9648-a9f02640aa48',
                                '62e29ffe-b3dd-4652-bc18-d4aebb76068f', 
                                'd7f12bdb-83ca-45c6-adc3-9142099f2816',
                                'a5359bd0-7e92-4ce8-99b5-e4035b7881e2')

delete from user_images where id in ('add893fb-d24d-4ee0-9976-03a6875cde8d');
