-- +migrate Up
alter table types add column tryonable boolean not null default false;

update types set tryonable = true
where name in ('Верх', 'Низ', 'Платья');

-- +migrate Down
alter table types drop column tryonable;
