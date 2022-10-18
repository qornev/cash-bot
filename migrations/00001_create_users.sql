-- +goose Up
-- +goose StatementBegin
create table users
(
    id   integer generated always as identity primary key,
    code text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
