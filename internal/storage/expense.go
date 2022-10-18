package storage

import (
	"context"
	"database/sql"

	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/domain"
)

type ExpenseDB struct {
	db *sql.DB
}

func NewExpenseDB(db *sql.DB) *ExpenseDB {
	return &ExpenseDB{db: db}
}

func (db *ExpenseDB) Add(ctx context.Context, date int64, userID int64, category string, amount float64) error {
	const query = `
		insert into expenses(
			dt, 
			user_id, 
			amount, 
			category
		) values (
			$1, $2, $3, $4
		);
	`

	_, err := db.db.ExecContext(ctx, query,
		date,
		userID,
		amount,
		category,
	)

	return err
}

// Return array with `userID` expenses
func (db *ExpenseDB) Get(ctx context.Context, userID int64) ([]*domain.Expense, error) {
	const query = `
		select
			dt, 
			category, 
			amount
		from expenses
		where user_id = $1
	`

	rows, err := db.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	expenses := make([]*domain.Expense, 0)

	for rows.Next() {
		var expense domain.Expense
		if err = rows.Scan(&expense); err != nil {
			return nil, err
		}
		expenses = append(expenses, &expense)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return expenses, nil
}
