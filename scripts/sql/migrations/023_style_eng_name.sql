-- +migrate Up
update styles set eng_name = c.eng_name
from (values
    ('Повседневный', 'casual wear'),
    ('Официальный', 'formal wear'),
    ('Спортивный', 'sportswear')
) as c(name, eng_name) 
where c.name = styles.name;

-- +migrate Down
update styles set eng_name = c.eng_name
from (values
    ('Повседневный', 'casual clothes'),
    ('Официальный', 'formal clothes'),
    ('Спортивный', 'sportswear')
) as c(name, eng_name) 
where c.name = styles.name;
