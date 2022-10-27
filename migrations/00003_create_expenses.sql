-- +goose Up
-- +goose StatementBegin
create table expenses
(
    id       integer generated always as identity primary key,
    dt       integer not null,
    user_id  integer not null,
    amount   real    not null,
    category text    not null
);
alter table expenses
    add constraint fk_expenses_users foreign key (user_id) references users(id);
create unique index expenses_idx on expenses (dt, user_id);
-- Создаю индекс по колонкам dt и user_id, потому что в запросах по ним происходит основное сравнение строк друг с другом
-- B-tree, потому что оно дает быстрый поиск в глубину, т.к. оно смотрит есть ли в текущем предке нужное значение
-- Если значение есть, то спускаемся вниз, если нет - смотрим по другим потомкам
-- Так же можно итерироваться между листьев
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index expenses_idx;
alter table expenses 
    drop constraint fk_expenses_users;
drop table expenses;
-- +goose StatementEnd
