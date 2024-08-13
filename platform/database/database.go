package database

import (
	_ "github.com/jackc/pgx/v5/stdlib"

	"my_project/app/stores"
)

type Stores struct {
	*stores.UsersStore
}

func GetDbConnection() (*Stores, error) {
	var err error
	db, err := postgresDbConnection()
	if err != nil {
		return nil, err
	}

	return &Stores{
		UsersStore: &stores.UsersStore{DB: db},
	}, nil
}
