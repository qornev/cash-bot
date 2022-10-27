package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/domain"
)

type ExpenseDB struct {
	db *sql.DB
}

func NewExpenseDB(db *sql.DB) *ExpenseDB {
	return &ExpenseDB{db: db}
}

var ErrorOverBudget = errors.New("Going over budget")

// Add expense with transaction.
// If month sum of expenses will be more than a budget, so transaction will rollback.
// If not more - commit.
func (db *ExpenseDB) Add(ctx context.Context, date int64, userID int64, category string, amount float64) (err error) {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = errors.Wrap(rollbackErr, err.Error())
			}
			return
		}
		err = tx.Commit()
	}()

	const insertQuery = `
		insert into expenses(
			dt, 
			user_id, 
			amount, 
			category
		) values (
			$1, $2, $3, $4
		);
	`

	_, err = tx.ExecContext(ctx, insertQuery, date, userID, amount, category)
	if err != nil {
		return err
	}

	const monthSumQuery = `
		select sum(amount)
		from expenses
		where date_trunc('month', to_timestamp($1)) = date_trunc('month', to_timestamp(dt));
	`

	var monthSum sql.NullFloat64
	err = tx.QueryRowContext(ctx, monthSumQuery, time.Now().Unix()).Scan(&monthSum)
	if err != nil {
		return err
	}

	// Sum is null, so no operations need to limit in current month
	if !monthSum.Valid {
		return nil
	}

	const budgetQuery = `
		select budget 
		from users
		where id = $1
	`

	var budget sql.NullFloat64
	err = tx.QueryRowContext(ctx, budgetQuery, userID).Scan(&budget)
	if err != nil {
		return err
	}

	// Budget is null, so no limit on operations
	if !budget.Valid {
		return nil
	}

	if monthSum.Float64 > budget.Float64 {
		return ErrorOverBudget
	}

	return nil
}

// Return array with `userID` expenses
func (db *ExpenseDB) Get(ctx context.Context, userID int64) ([]domain.Expense, error) {
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
	expenses := make([]domain.Expense, 0)

	for rows.Next() {
		var expense domain.Expense
		if err = rows.Scan(&expense.Date, &expense.Category, &expense.Amount); err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return expenses, nil
}
