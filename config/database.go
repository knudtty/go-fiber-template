package config

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitDb() {
    var err error
    dsn := "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@db:5432/" + os.Getenv("POSTGRES_DB") + "?sslmode=disable"
    DB, err = sqlx.Connect("postgres", dsn)
    if err != nil {
        log.Fatalln(err)
    }
}
