package config

import (
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitDb() {
	var err error
	dsn := "postgres://" +
		os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASSWORD") + "@" +
		os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + "/" +
		os.Getenv("DB_NAME") + "?sslmode=" +
		os.Getenv("DB_SSL_MODE")

	DB, err = sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	maxIdleConns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	maxConns, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNS"))
	maxConnLifetime, _ := strconv.Atoi(os.Getenv("DB_MAX_CONN_LIFETIME"))

	DB.SetMaxIdleConns(maxIdleConns)
	DB.SetMaxOpenConns(maxConns)
	DB.SetConnMaxLifetime(time.Duration(maxConnLifetime))

	if err = DB.Ping(); err != nil {
		DB.Close()
		log.Fatal("Failed to ping db", err)
	}
}
