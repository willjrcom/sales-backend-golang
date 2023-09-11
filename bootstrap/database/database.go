package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// NewPostgreSQLConnection creates a connection with PostgreSQL using sqlx
func NewPostgreSQLConnection() (*bun.DB, error) {
	// Prepare connection string parameterized
	connectionParams := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		"name",
		"password",
		"host",
		"5432",
		"dbname",
	)

	// Connect to database doing a PING
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connectionParams), pgdriver.WithTimeout(time.Second*30)))

	// set connection settings
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Duration(60) * time.Minute)
	db := bun.NewDB(sqlDB, pgdialect.New())

	return db, nil
}
