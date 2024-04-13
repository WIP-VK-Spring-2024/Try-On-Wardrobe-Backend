-- +migrate Up
insert into outfit_purpose(name, eng_name)
values ('На работу', 'office clothes'),
       ('Для прогулок', 'clothes for walks'),
       ('Для активного отдыха', 'clothes for outdoor activities'),
       ('В спортзал', 'clothes for gym'),
       ('Для пляжа', 'clothes for the beach'),
       ('В театр', 'clothes for theater');
