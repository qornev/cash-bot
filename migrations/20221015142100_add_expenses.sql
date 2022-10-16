-- +goose Up
-- +goose StatementBegin
create table expenses
(
    id         integer generated always as identity primary key,
    created_at timestamp not null default now(),
    unix_ts    integer   not null,
    user_id    integer   not null,
    amount     real      not null,
    category   text      not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table rates;
-- +goose StatementEnd
