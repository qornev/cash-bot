package storage

import (
	"context"
	"database/sql"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/alex1234562557/telegram-bot/internal/converter"
)

type RateDB struct {
	db *sql.DB
}

func NewRateDB(db *sql.DB) *RateDB {
	return &RateDB{db: db}
}

func (db *RateDB) Add(ctx context.Context, date int64, code string, nominal float64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "add rate to db")
	defer span.Finish()

	const query = `
		insert into rates(
			dt, 
			code, 
			nominal
		) values (
			$1, $2, $3
		);
	`
	_, err := db.db.ExecContext(ctx, query,
		date,
		code,
		nominal,
	)

	return err
}

// Get `code` rate at `date` time
func (db *RateDB) Get(ctx context.Context, date int64, code string) (*converter.Rate, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "get rate from db")
	defer span.Finish()

	const query = `
		select 
			code, 
			nominal
		from rates
		where 1 = 1
			and code = $1
			and dt <= $2
			and $2 - dt <= 24 * 60 * 60 -- difference in 1 day
		order by dt desc;
	`
	var rate converter.Rate
	err := db.db.QueryRowContext(ctx, query, code, date).Scan(&rate.Code, &rate.Nominal)

	return &rate, err
}
