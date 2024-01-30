package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
	userentity "github.com/willjrcom/sales-backend-go/internal/domain/user"
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

	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), schemaentity.LOST_SCHEMA)
	if err := CreateSchema(ctx, bun); err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), schemaentity.DEFAULT_SCHEMA)
	if err := CreateSchema(ctx, bun); err != nil {
		return nil, err
	}

	if err := defaultTables(ctx, bun); err != nil {
		return nil, err
	}

	fmt.Println("Db connected")
	return bun, nil
}

func ChangeSchema(ctx context.Context, db *bun.DB) error {
	schemaName, err := GetSchema(ctx)

	if err != nil {
		return err
	}

	_, err = db.Exec("SET search_path=?", schemaName)
	return err
}

func GetSchema(ctx context.Context) (string, error) {
	schemaName := ctx.Value(schemaentity.Schema("schema"))
	if schemaName == nil {
		return "", errors.New("schema not found")
	}
	return schemaName.(string), nil
}

func CreateSchema(ctx context.Context, db *bun.DB) error {
	schemaName, err := GetSchema(ctx)

	if err != nil {
		schemaName = schemaentity.LOST_SCHEMA
	}

	if _, err := db.Exec("CREATE SCHEMA IF NOT EXISTS " + schemaName); err != nil {
		return err
	}

	return nil
}

func defaultTables(ctx context.Context, db *bun.DB) error {
	db.RegisterModel((*userentity.User)(nil))

	if _, err := db.NewCreateTable().IfNotExists().Model((*userentity.User)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS pgcrypto;"); err != nil {
		return err
	}

	return nil
}
