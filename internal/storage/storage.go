package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type ConfigGetter interface {
	HostDB() string
	Port() int
	Username() string
	Password() string
}

func Connect(configGetter ConfigGetter) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		configGetter.HostDB(),
		configGetter.Port(),
		configGetter.Username(),
		configGetter.Password(),
	)
	db, err := sql.Open("postgres", dataSourceName)
	return db, err
}
