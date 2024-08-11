package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func postgresDbConnection() (*sqlx.DB, error) {
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

	maxIdleConns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	maxConns, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNS"))
	maxConnLifetime, _ := strconv.Atoi(os.Getenv("DB_MAX_CONN_LIFETIME"))

	sqlxConn.SetMaxIdleConns(maxIdleConns)
	sqlxConn.SetMaxOpenConns(maxConns)
	sqlxConn.SetConnMaxLifetime(time.Duration(maxConnLifetime))

	if err = sqlxConn.Ping(); err != nil {
		defer sqlxConn.Close()
		return nil, fmt.Errorf("Failed to ping sqlxConn: %s\n", err)
	}

	return sqlxConn, nil
}
