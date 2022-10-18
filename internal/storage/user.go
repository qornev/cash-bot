package storage

import (
	"context"
	"database/sql"

	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
)

type UserDB struct {
	db *sql.DB
}

func NewUserDB(db *sql.DB) *UserDB {
	return &UserDB{db: db}
}

func (db *UserDB) Set(ctx context.Context, userID int64, code string) error {
	isUserExist, err := db.UserExist(ctx, userID)
	if err != nil {
		return err
	}

	var query string
	if isUserExist {
		query = `
			update users 
			set code = $1 
			where id = $2;
		`
	} else {
		query = `
			insert into users(
				id, 
				code
			) values (
				$2, $1
			);
		`
	}

	_, err = db.db.ExecContext(ctx, query,
		code,
		userID,
	)

	return err
}

func (db *UserDB) UserExist(ctx context.Context, userID int64) (bool, error) {
	const query = `
		select 
			id 
		from users
		where user_id = $1;
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

func (db *UserDB) Get(ctx context.Context, userID int64) (string, error) {
	const query = `
		select
			code
		from users
		where user_id = $1;
	`

	var code string
	err := db.db.QueryRowContext(ctx, query, userID).Scan(&code)

	switch {
	case err == sql.ErrNoRows:
		err := db.Set(ctx, userID, converter.RUB)
		return converter.RUB, err
	case err != nil:
		return "", err
	default:
		return code, nil
	}

}
