-- +migrate Up
alter table types add column eng_name text;

alter table subtypes add column eng_name text;

alter table styles add column eng_name text;

create index subtypes_eng_name_idx on subtypes(eng_name);

update types t set
    eng_name = c.eng_name
from (values
    ('Верх', 'upper garment'),
    ('Низ', 'lower garment'),
    ('Обувь', 'shoe'),
    ('Нижнее бельё', 'underwear'),
    ('Верзняя одежда', 'outerwear'),
    ('Обувь', 'shoe'),
    ('Платья', 'dress'),
    ('Головные уборы', 'hat'),
    ('Аксессуары', 'accessory')
) as c(name, eng_name) 
where c.name = t.name;

update subtypes
set name = 'Бюстгальтеры'
where name = 'Бюстгалтеры'; -- lmao

update subtypes s set
    eng_name = c.eng_name
from (values
    ('Летние', 'summer dress'),
    ('Вечерние', 'evening dress'),
    ('Деловые', 'business dress'),
    ('Кепки', 'baseball cap'),
    ('Панамы', 'panama'),
    ('Зимние шапки', 'winter cap'),
    ('Осенние шапки', 'autumn cap'),
    ('Украшения', 'jewelry'),
    ('Часы', 'watch'),
    ('Шарфы', 'scarf'),
    ('Перчатки', 'gloves'),
    ('Ремни', 'belt'),
    ('Галстуки', 'tie'),
    ('Сумки', 'bag'),
    ('Майки', 'tank top'),
    ('Бюстгальтеры', 'bra'),
    ('Трусы', 'pants'),
    ('Носки', 'socks'),
    ('Колготки', 'stockings'),
    ('Кальсоны', 'underpants'),
    ('Термобельё', 'thermal underwear'),
    ('Футболки', 't-shirt'),
    ('Поло', 'polo'),
    ('Лонгслив', 'longsleeve'),
    ('Свитеры', 'sweater'),
    ('Худи', 'hoodie'),
    ('Свитшоты', 'sweatshirt'),
    ('Топы', 'top'),
    ('Рубашки', 'shirt'),
    ('Кофты', 'blouse'),
    ('Джинсы', 'jeans'),
    ('Брюки', 'slacks'),
    ('Шорты', 'shorts'),
    ('Юбки', 'skirt'),
    ('Штаны', 'trousers'),
    ('Леггинсы', 'leggins'),
    ('Бриджи', 'breeches'),
    ('Полусапоги', 'half boots'),
    ('Туфли', 'shoes'),
    ('Сапоги', 'boots'),
    ('Кеды', 'sneakers'),
    ('Кроссовки', 'sneakers'),
    ('Босоножки', 'sandals'),
    ('Ботинки', 'boots'),
    ('Полуботинки', 'low shoes')
) as c(name, eng_name) 
where c.name = s.name;

update styles s set
    eng_name = c.eng_name
from (values
    ('Повседневный', 'casual clothes'),
    ('Официальный', 'formal clothes'),
    ('Спортивный', 'sportswear')
) as c(name, eng_name) 
where c.name = s.name;

alter table types alter column eng_name set not null;

alter table subtypes alter column eng_name set not null;

alter table styles alter column eng_name set not null;

-- +migrate Down
alter table types drop column eng_name;

alter table subtypes drop column eng_name;

alter table styles drop column eng_name;
