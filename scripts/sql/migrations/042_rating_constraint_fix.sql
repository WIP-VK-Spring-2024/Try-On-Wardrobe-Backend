-- +migrate Up
alter table post_ratings
    drop constraint post_ratings_value_check,
    add constraint post_ratings_value_check check (value in (-1, 0, 1));

alter table post_comment_ratings
    add constraint post_comment_ratings_value_check check (value in (-1, 0, 1));

-- +migrate Down
alter table post_comment_ratings
    drop constraint post_comment_ratings_value_check;

alter table post_ratings
    drop constraint post_ratings_value_check,
    add constraint post_ratings_value_check check (value in (-1, 1));
