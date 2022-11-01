-- +goose Up
-- +goose StatementBegin
insert into users(id, code, updated_at)
values (1234, 'RUB', 1666186334);
insert into users(id, code, updated_at)
values (4321, 'RUB', 1666186334);
insert into users(id, code, updated_at)
values (102938, 'RUB', 1666186334);
insert into users(id, code, updated_at)
values (47389, 'RUB', 1666186334);
insert into users(id, code, updated_at)
values (213900, 'RUB', 1666186334);
insert into users(id, code, updated_at)
values (3452, 'RUB', 1666186334);

insert into expenses(dt, user_id, amount, category)
values (1666186334, 1234, 123.45, 'еда');
insert into expenses(dt, user_id, amount, category)
values (1665186334, 1234, 800.00, 'интернет');
insert into expenses(dt, user_id, amount, category)
values (1664186334, 4321, 1500.31, 'жкх');
insert into expenses(dt, user_id, amount, category)
values (1663186334, 4321, 700.11, 'кино');
insert into expenses(dt, user_id, amount, category)
values (1662186334, 102938, 432.45, 'шампунь');
insert into expenses(dt, user_id, amount, category)
values (1661186334, 47389, 500.00, 'потерял');
insert into expenses(dt, user_id, amount, category)
values (1659186334, 213900, 780.33, 'лекарства');
insert into expenses(dt, user_id, amount, category)
values (1658186334, 3452, 1500.00, 'врач');
insert into expenses(dt, user_id, amount, category)
values (1659186334, 1234, 345.45, 'еда');
insert into expenses(dt, user_id, amount, category)
values (1662186334, 1234, 1.45, 'пакет');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate table expenses;
-- +goose StatementEnd
