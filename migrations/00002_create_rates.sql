-- +goose Up
-- +goose StatementBegin
create table rates 
(
    id         integer generated always as identity primary key,
    dt         integer   not null,
    code       text      not null,
    nominal    real      not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists rates;
-- +goose StatementEnd
