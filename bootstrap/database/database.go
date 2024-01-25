package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var (
	username = "admin"
	password = "admin"
	host     = "localhost"
	port     = "5432"
	dbName   = "sales-db"
)

func NewPostgreSQLConnection(ctx context.Context) (*bun.DB, error) {
	// Prepare connection string parameterized
	connectionParams := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username,
		password,
		host,
		port,
		dbName,
	)

	// Connect to database doing a PING
	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connectionParams), pgdriver.WithTimeout(time.Second*30)))

	// Verifique se o banco de dados já existe.
	if err := db.Ping(); err != nil {
		log.Printf("erro ao conectar ao banco de dados: %v", err)
	}

	// set connection settings
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Duration(60) * time.Minute)

	bun := bun.NewDB(db, pgdialect.New())
	fmt.Println("Db connected")
	return bun, nil
}

func ChangeSchema(db *bun.DB, schemaName string) error {
	_, err := db.Exec("SET search_path=?", schemaName)
	return err
}
