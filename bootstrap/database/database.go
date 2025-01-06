package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Environment string

func ConnectLocalDB(ctx context.Context) string {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost" // Valor padrão para desenvolvimento local
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		"admin",
		"admin",
		dbHost,
		"5432",
		"sales-db",
	)
}

func ConnectRdsDB(ctx context.Context) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?",
		"postgres",
		"48279111",
		"sales-backend-db1.c7ou20us0aar.us-east-1.rds.amazonaws.com",
		"5432",
		"salesBackendDB",
	)
}

func NewPostgreSQLConnection() *bun.DB {
	ctx := context.Background()
	// Prepare connection string parameterized
	connectionParams := ""
	environment := ctx.Value(Environment("environment"))
	fmt.Print(environment)
	if environment == "prod" {
		connectionParams = ConnectRdsDB(ctx)
	} else {
		connectionParams = ConnectLocalDB(ctx)
	}

	// Connect to database doing a PING
	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connectionParams), pgdriver.WithTimeout(time.Second*30)))

	// Verifique se o banco de dados já existe.
	if err := db.Ping(); err != nil {
		log.Printf("erro ao conectar ao banco de dados: %v", err)
		panic(err)
	}

	// set connection settings
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Duration(60) * time.Minute)

	dbBun := bun.NewDB(db, pgdialect.New())

	if err := LoadAllSchemas(ctx, dbBun); err != nil {
		panic(err)
	}

	if err := publicTables(ctx, dbBun); err != nil {
		panic(err)
	}

	fmt.Println("Db connected")
	return dbBun
}

func ChangeSchema(ctx context.Context, db *bun.DB) error {
	schemaName, err := GetSchema(ctx)

	if err != nil {
		return err
	}

	_, err = db.Exec("SET search_path=?", schemaName)
	return err
}

func ChangeToPublicSchema(ctx context.Context, db *bun.DB) error {
	_, err := db.Exec("SET search_path=?", model.DEFAULT_SCHEMA)
	return err
}

func GetSchema(ctx context.Context) (string, error) {
	schemaName := ctx.Value(model.Schema("schema"))
	if schemaName == nil {
		return "", errors.New("schema not found")
	}
	return schemaName.(string), nil
}

func CreateSchema(ctx context.Context, db *bun.DB) error {
	schemaName, err := GetSchema(ctx)

	if err != nil {
		schemaName = model.LOST_SCHEMA
	}

	if _, err := db.Exec("CREATE SCHEMA IF NOT EXISTS " + schemaName); err != nil {
		return err
	}

	return nil
}

func publicTables(ctx context.Context, db *bun.DB) error {
	ctx = context.WithValue(ctx, model.Schema("schema"), model.DEFAULT_SCHEMA)
	if err := CreateSchema(ctx, db); err != nil {
		panic(err)
	}

	db.RegisterModel((*model.User)(nil))
	db.RegisterModel((*model.CompanyToUsers)(nil))
	db.RegisterModel((*model.CompanyWithUsers)(nil))
	db.RegisterModel((*model.Person)(nil))
	db.RegisterModel((*model.Address)(nil))
	db.RegisterModel((*model.Contact)(nil))

	if err := RegisterModels(ctx, db); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.User)(nil)).Exec(ctx); err != nil {
		return err
	}

	var _ bun.BeforeSelectHook = (*model.User)(nil)

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.CompanyToUsers)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.CompanyWithUsers)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.Person)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.Address)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.Contact)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS pgcrypto;"); err != nil {
		return err
	}

	return nil
}

func LoadAllSchemas(ctx context.Context, db *bun.DB) error {
	results, err := db.QueryContext(ctx, "SELECT schema_name FROM information_schema.schemata;")

	if err != nil {
		return err
	}

	for results.Next() {
		var schemaName string
		if err := results.Scan(&schemaName); err != nil {
			return err
		}

		if !strings.Contains(schemaName, "loja_") {
			continue
		}

		ctx = context.WithValue(ctx, model.Schema("schema"), schemaName)

		if err := RegisterModels(ctx, db); err != nil {
			return err
		}

		if err := LoadCompanyModels(ctx, db); err != nil {
			return err
		}

		if _, err := db.QueryContext(ctx, "CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_contact ON contacts (ddd, number, type);"); err != nil {
			return err
		}

		if err := SetupContactFtSearch(ctx, db); err != nil {
			return err
		}
	}

	return nil
}

func RegisterModels(ctx context.Context, db *bun.DB) error {
	mu := sync.Mutex{}

	mu.Lock()
	defer mu.Unlock()

	if err := CreateSchema(ctx, db); err != nil {
		return err
	}

	if err := ChangeSchema(ctx, db); err != nil {
		return err
	}

	db.RegisterModel((*entitymodel.Entity)(nil))

	db.RegisterModel((*model.ProductCategoryToAdditional)(nil))
	db.RegisterModel((*model.ProductCategoryToComplement)(nil))
	db.RegisterModel((*model.ProductToCombo)(nil))
	db.RegisterModel((*model.Size)(nil))
	db.RegisterModel((*model.Quantity)(nil))
	db.RegisterModel((*model.ProductCategory)(nil))
	db.RegisterModel((*model.ProcessRule)(nil))
	db.RegisterModel((*model.Product)(nil))

	db.RegisterModel((*model.Address)(nil))
	db.RegisterModel((*model.Contact)(nil))
	db.RegisterModel((*model.Client)(nil))
	db.RegisterModel((*model.Employee)(nil))

	db.RegisterModel((*model.OrderProcessToProductToGroupItem)(nil))
	db.RegisterModel((*model.OrderProcess)(nil))
	db.RegisterModel((*model.OrderQueue)(nil))
	db.RegisterModel((*model.ItemToAdditional)(nil))
	db.RegisterModel((*model.Item)(nil))
	db.RegisterModel((*model.GroupItem)(nil))

	var _ bun.BeforeUpdateHook = (*model.GroupItem)(nil)

	db.RegisterModel((*model.OrderPickup)(nil))
	db.RegisterModel((*model.OrderDelivery)(nil))
	db.RegisterModel((*model.DeliveryDriver)(nil))
	db.RegisterModel((*model.OrderTable)(nil))
	db.RegisterModel((*model.PaymentOrder)(nil))
	db.RegisterModel((*model.Order)(nil))
	// var _ bun.AfterScanRowHook = (*model.Order)(nil)

	db.RegisterModel((*model.Table)(nil))
	db.RegisterModel((*model.PlaceToTables)(nil))
	db.RegisterModel((*model.Place)(nil))

	db.RegisterModel((*model.Shift)(nil))
	db.RegisterModel((*model.Company)(nil))

	return nil
}

func LoadCompanyModels(ctx context.Context, db *bun.DB) error {
	mu := sync.Mutex{}

	mu.Lock()
	defer mu.Unlock()

	if err := ChangeSchema(ctx, db); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.Size)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.ProductCategoryToAdditional)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.ProductCategoryToComplement)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.ProductToCombo)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.ProductCategory)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.Quantity)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.ProcessRule)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.Product)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.Address)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.Contact)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.Client)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.Employee)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.OrderProcess)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.OrderQueue)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.GroupItem)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.OrderProcessToProductToGroupItem)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.Item)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.ItemToAdditional)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.OrderPickup)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.OrderDelivery)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.DeliveryDriver)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.OrderTable)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.PaymentOrder)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.Order)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.Table)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.Place)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.PlaceToTables)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateIndex().IfNotExists().Model((*model.PlaceToTables)(nil)).Unique().Index("idx_place_id_row_and_column").Column("place_id", "row", "column").Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.Shift)(nil)).Exec(ctx); err != nil {
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*model.Company)(nil)).Exec(ctx); err != nil {
		return err
	}

	return nil
}
