-- +migrate Up
create unique index if not exists users_name_idx on users (name varchar_pattern_ops);

create unique index if not exists users_email_idx on users (lower(email));

create index if not exists clothes_name_idx on clothes (name varchar_pattern_ops);

create index if not exists clothes_style_id_idx on clothes (style_id);

create index if not exists clothes_user_id_idx on clothes (user_id);

create index if not exists clothes_type_id_idx on clothes (type_id);

create unique index if not exists tags_name_idx on tags (name varchar_pattern_ops);

create unique index if not exists styles_name_idx on styles (name varchar_pattern_ops);

create unique index if not exists types_name_idx on types (name varchar_pattern_ops);

create index if not exists subtypes_name_idx on subtypes (name varchar_pattern_ops);

-- +migrate Down
drop index users_name_idx;

drop index users_email_idx;

drop index clothes_name_idx;

drop index clothes_style_id_idx;

drop index clothes_user_id_idx;

drop index clothes_type_id_idx;

drop index tags_name_idx;

drop index styles_name_idx;

drop index types_name_idx;

drop index subtypes_name_idx;
