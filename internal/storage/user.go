package storage

import (
	"context"
	"database/sql"
	"time"

	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/domain"
)

type UserDB struct {
	db *sql.DB
}

func NewUserDB(db *sql.DB) *UserDB {
	return &UserDB{db: db}
}

func (db *UserDB) UserExist(ctx context.Context, userID int64) (bool, error) {
	const query = `
		select 
			id 
		from users
		where id = $1;
	`
	var id int64
	err := db.db.QueryRowContext(ctx, query, userID).Scan(&id)

	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (db *UserDB) SetCode(ctx context.Context, userID int64, code string) error {
	isUserExist, err := db.UserExist(ctx, userID)
	if err != nil {
		return err
	}

	if isUserExist {
		return db.UpdateCode(ctx, userID, code)
	}
	return db.AddCode(ctx, userID, code)
}

func (db *UserDB) UpdateCode(ctx context.Context, userID int64, code string) error {
	const query = `
		update users 
		set 
			code = $1,
			updated_at = $2
		where id = $3;
	`

	_, err := db.db.ExecContext(ctx, query,
		code,
		time.Now().Unix(),
		userID,
	)

	return err
}

func (db *UserDB) AddCode(ctx context.Context, userID int64, code string) error {
	const query = `
		insert into users(
			code,
			id,
			updated_at
		) values (
			$1, $2, $3
		);
	`

	_, err := db.db.ExecContext(ctx, query,
		code,
		userID,
		time.Now().Unix(),
	)

	return err
}

func (db *UserDB) GetCode(ctx context.Context, userID int64) (string, error) {
	const query = `
		select
			code
		from users
		where id = $1;
	`

	var code string
	err := db.db.QueryRowContext(ctx, query, userID).Scan(&code)

	switch {
	case err == sql.ErrNoRows:
		err := db.AddCode(ctx, userID, converter.RUB)
		return converter.RUB, err
	case err != nil:
		return "", err
	default:
		return code, nil
	}
}

func (db *UserDB) SetBudget(ctx context.Context, userID int64, budget float64) error {
	isUserExist, err := db.UserExist(ctx, userID)
	if err != nil {
		return err
	}

	if isUserExist {
		return db.UpdateBudget(ctx, userID, budget)
	}
	return db.AddBudget(ctx, userID, budget)
}

func (db *UserDB) UpdateBudget(ctx context.Context, userID int64, budget float64) error {
	const query = `
		update users 
		set 
			budget = $1, 
			updated_at = $2
		where id = $3;
	`

	_, err := db.db.ExecContext(ctx, query,
		budget,
		time.Now().Unix(),
		userID,
	)

	return err
}

func (db *UserDB) AddBudget(ctx context.Context, userID int64, budget float64) error {
	const query = `
		insert into users(
			id, 
			code,
			budget,
			updated_at
		) values (
			$1, $2, $3, $4
		);
	`

	_, err := db.db.ExecContext(ctx, query,
		userID,
		converter.RUB,
		budget,
		time.Now().Unix(),
	)

	return err
}

func (db *UserDB) GetBudget(ctx context.Context, userID int64) (*float64, string, int64, error) {
	const query = `
		select
			budget,
			code,
			updated_at
		from users
		where id = $1;
	`

	var budget *float64
	var code string
	var date int64
	err := db.db.QueryRowContext(ctx, query, userID).Scan(&budget, &code, &date)

	switch {
	case err == sql.ErrNoRows:
		err = db.AddCode(ctx, userID, converter.RUB)
		return nil, code, date, err
	case err != nil:
		return nil, "", 0, err
	default:
		return budget, code, date, nil
	}
}

func (db *UserDB) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	const query = `
		select 
			id,
			code,
			budget,
			updated_at
		from users;
	`

	rows, err := db.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]domain.User, 0)

	for rows.Next() {
		var user domain.User
		if err = rows.Scan(&user.ID, &user.Code, &user.Budget, &user.Updated); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}
