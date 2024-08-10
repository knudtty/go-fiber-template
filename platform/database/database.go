package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func PostgresDbConnection() (*sqlx.DB, error) {
	if db != nil {
		return db, nil
	}

	dsn := "postgres://" +
		os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASSWORD") + "@" +
		os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + "/" +
		os.Getenv("DB_NAME") + "?sslmode=" +
		os.Getenv("DB_SSL_MODE")

	sqlxConn, err := sqlx.Connect("pgx", dsn)
	if err != nil {
        return nil, fmt.Errorf("Error connecting to db: %s", err)
	}
	db = sqlxConn

	maxIdleConns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	maxConns, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNS"))
	maxConnLifetime, _ := strconv.Atoi(os.Getenv("DB_MAX_CONN_LIFETIME"))

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxConns)
	db.SetConnMaxLifetime(time.Duration(maxConnLifetime))

	if err = db.Ping(); err != nil {
		defer db.Close()
		return nil, fmt.Errorf("Failed to ping db: %s\n", err)
	}

	return db, nil
}
