-- +goose Up
-- +goose StatementBegin
create table expenses
(
    id         integer generated always as identity primary key,
    dt         integer   not null,
    user_id    integer   not null,
    amount     real      not null,
    category   text      not null
);
alter table expenses
    add constraint fk_expenses_users foreign key (user_id) references users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table expenses 
    drop constraint fk_expenses_users;
drop table expenses;
-- +goose StatementEnd
