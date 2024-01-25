package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
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

	LoadModels(ctx, bun, "public")
	return bun, nil
}

func ChangeSchema(db *bun.DB, schemaName string) error {
	_, err := db.Exec("SET search_path=?", schemaName)
	return err
}

func LoadModels(ctx context.Context, db *bun.DB, schema string) {
	mu := sync.Mutex{}
	mu.Lock()

	if _, err := db.Exec("CREATE SCHEMA IF NOT EXISTS " + schema); err != nil {
		panic(err)
	}

	if err := ChangeSchema(db, schema); err != nil {
		mu.Unlock()
		panic(err)
	}

	db.RegisterModel((*entity.Entity)(nil))

	db.RegisterModel((*productentity.Size)(nil))
	db.RegisterModel((*productentity.Quantity)(nil))
	db.RegisterModel((*productentity.CategoryToAditionalCategories)(nil))
	db.RegisterModel((*productentity.Category)(nil))
	db.RegisterModel((*productentity.Process)(nil))
	db.RegisterModel((*productentity.Product)(nil))

	db.RegisterModel((*addressentity.Address)(nil))
	db.RegisterModel((*personentity.Contact)(nil))
	db.RegisterModel((*cliententity.Client)(nil))
	db.RegisterModel((*employeeentity.Employee)(nil))

	db.RegisterModel((*itementity.Item)(nil))
	db.RegisterModel((*groupitementity.GroupItem)(nil))

	db.RegisterModel((*orderentity.DeliveryOrder)(nil))
	db.RegisterModel((*orderentity.TableOrder)(nil))
	db.RegisterModel((*orderentity.PaymentOrder)(nil))
	db.RegisterModel((*orderentity.Order)(nil))

	db.RegisterModel((*tableentity.Table)(nil))
	db.RegisterModel((*shiftentity.Shift)(nil))

	if _, err := db.NewCreateTable().IfNotExists().Model((*entity.Entity)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for entity")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.Size)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for size category")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.CategoryToAditionalCategories)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for category to aditional")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.Category)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for category product")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.Quantity)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for quantity category")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.Process)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for process category")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.Product)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for product")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*addressentity.Address)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for address")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*personentity.Contact)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for address")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*personentity.Person)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for person")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*cliententity.Client)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for client")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*employeeentity.Employee)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for employee")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*itementity.Item)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for item")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*groupitementity.GroupItem)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for items")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*orderentity.DeliveryOrder)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for delivery order")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*orderentity.TableOrder)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for table order")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*orderentity.PaymentOrder)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for payment order")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*orderentity.Order)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for order")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*tableentity.Table)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for table")
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*shiftentity.Shift)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		panic("Couldn't create table for shift")
	}
}
