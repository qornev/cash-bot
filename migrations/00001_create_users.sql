-- +goose Up
-- +goose StatementBegin
create table users
(
    id         integer primary key,
    code       text    not null,
    budget     real,
    updated_at integer not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
