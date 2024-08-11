package database

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"my_project/app/stores"
)

var db *sqlx.DB

type Stores struct {
	*stores.UsersStore
}

func GetDbConnection() (*Stores, error) {
	if db == nil {
		var err error
		db, err = postgresDbConnection()
		if err != nil {
			return nil, err
		}
	}

	return &Stores{
		UsersStore: &stores.UsersStore{DB: db},
	}, nil
}
