-- +migrate Up
delete from subtypes
where type_id = '0a45a9f1-992a-43c8-9ffd-5c5a00083297';

delete from types
where id = '0a45a9f1-992a-43c8-9ffd-5c5a00083297';
