package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type ConfigGetter interface {
	HostDB() string
	PortDB() int
	UsernameDB() string
	PasswordDB() string
}

func Connect(configGetter ConfigGetter) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		configGetter.HostDB(),
		configGetter.PortDB(),
		configGetter.UsernameDB(),
		configGetter.PasswordDB(),
	)
	db, err := sql.Open("postgres", dataSourceName)
	return db, err
}
