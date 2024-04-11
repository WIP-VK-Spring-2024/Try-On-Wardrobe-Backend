-- +migrate Up
alter table subtypes add column temp_range int4range not null default '(,)'::int4range;

update subtypes
set temp_range = tmp.temp_range
from (values
    ('Кофты', '(,20]'::int4range),
    ('Майки', '[15,)'::int4range),
    ('Лонгслив', '(,20]'::int4range),
    ('Свитеры', '(,20]'::int4range),
    ('Худи', '(,20]'::int4range),
    ('Свитшоты', '(,20]'::int4range),

    ('Шорты', '[15,)'::int4range),

    ('Ветровки', '[10,)'::int4range),
    ('Пальто', '(,15]'::int4range),
    ('Шубы', '(,0]'::int4range)
) as tmp(name, temp_range)
where subtypes.name = tmp.name;

-- +migrate Down
alter table subtypes drop column temp_range;
