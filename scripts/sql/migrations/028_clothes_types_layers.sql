-- +migrate Up
insert into
types(id, name, eng_name, tryonable)
values ('8b31dc49-5f9b-451f-85c6-ecf17a3a024c', 'Верхняя одежда', 'outerwear', false);

insert into
subtypes(type_id, id, name, eng_name)
values
    ('8b31dc49-5f9b-451f-85c6-ecf17a3a024c', 'a541205c-2edb-4276-8df5-1dd0373cf29e', 'Ветровки', 'windbreaker'),
    ('8b31dc49-5f9b-451f-85c6-ecf17a3a024c', '3954df5e-b37e-4afe-9974-0d3ca94b0f9d', 'Куртки', 'jacket'),
    ('8b31dc49-5f9b-451f-85c6-ecf17a3a024c', '9803c786-755f-458b-bc0e-dde321f657a1', 'Пальто', 'coat'),
    ('8b31dc49-5f9b-451f-85c6-ecf17a3a024c', '59bac08e-efe2-4e9f-ae6b-4eac73d4bd27', 'Шубы', 'fur coat'),
    ('8b31dc49-5f9b-451f-85c6-ecf17a3a024c', 'f885aa1f-74ca-430e-8755-33be9cac8787', 'Пуховики', 'down jacket');

alter table subtypes add column layer smallint not null default 0;

update subtypes
set layer = tmp.layer
from (values
    ('Кофты', 1),
    ('Рубашки', 1),
    ('Лонгслив', 1),
    ('Свитеры', 2),
    ('Худи', 2),
    ('Свитшоты', 2)
) as tmp(name, layer) 
where tmp.name = subtypes.name;

update subtypes
set layer = 3
where type_id = '8b31dc49-5f9b-451f-85c6-ecf17a3a024c';

-- +migrate Down
alter table subtypes drop column layer; 
