package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Environment string

func ConnectDB(ctx context.Context) string {
	dbHost := os.Getenv("DB_HOST")
	if dbHost != "" {
		return dbHost
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		"admin",
		"admin",
		"localhost",
		"5432",
		"sales-db",
	)
}

func NewPostgreSQLConnection() *bun.DB {
	ctx := context.Background()
	connectionParams := ConnectDB(ctx)

	// Connect to database doing a PING
	db := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connectionParams), pgdriver.WithTimeout(time.Second*30)))

	// Verifique se o banco de dados já existe.
	if err := db.PingContext(ctx); err != nil {
		log.Printf("erro ao conectar ao banco de dados: %v", err)
		panic(err)
	}

	// set connection settings
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Duration(60) * time.Minute)

	dbBun := bun.NewDB(db, pgdialect.New())

	if err := registerModels(dbBun); err != nil {
		panic(err)
	}

	if err := createAllSchemaTables(ctx, dbBun); err != nil {
		panic(err)
	}

	if err := createPublicTables(ctx, dbBun); err != nil {
		panic(err)
	}

	fmt.Println("Db connected")
	return dbBun
}

func ChangeSchema(ctx context.Context, db *bun.DB) error {
	schemaName, err := GetCurrentSchema(ctx)

	if err != nil {
		return err
	}

	query := fmt.Sprintf("SET search_path=%s", schemaName)
	_, err = db.ExecContext(ctx, query)
	return err
}

func ChangeToPublicSchema(ctx context.Context, db *bun.DB) error {
	query := fmt.Sprintf("SET search_path=%s", model.PUBLIC_SCHEMA)
	_, err := db.ExecContext(ctx, query)
	return err
}

func GetCurrentSchema(ctx context.Context) (string, error) {
	schemaName := ctx.Value(model.Schema("schema"))
	if schemaName == nil {
		return "", errors.New("schema not found")
	}
	return schemaName.(string), nil
}

func createSchema(ctx context.Context, db *bun.DB) error {
	schemaName, err := GetCurrentSchema(ctx)
	if err != nil {
		return err
	}

	if _, err := db.Exec("CREATE SCHEMA IF NOT EXISTS " + schemaName); err != nil {
		return err
	}

	return nil
}

func createPublicTables(ctx context.Context, db *bun.DB) error {
	ctx = context.WithValue(ctx, model.Schema("schema"), model.PUBLIC_SCHEMA)
	if err := createSchema(ctx, db); err != nil {
		panic(err)
	}

	if err := createTableIfNotExists(ctx, db, (*model.CompanyToUsers)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.User)(nil)); err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS pgcrypto;"); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.Company)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.Address)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.Contact)(nil)); err != nil {
		return err
	}

	return nil
}

func createAllSchemaTables(ctx context.Context, db *bun.DB) error {
	schemasFound, err := db.QueryContext(ctx, "SELECT schema_name FROM information_schema.schemata;")

	if err != nil {
		return err
	}

	for schemasFound.Next() {
		var schemaName string
		if err := schemasFound.Scan(&schemaName); err != nil {
			return err
		}

		if !strings.Contains(schemaName, "loja_") {
			continue
		}

		ctx = context.WithValue(ctx, model.Schema("schema"), schemaName)

		if err := CreateNewCompanySchema(ctx, db); err != nil {
			return err
		}
	}

	return nil
}

func CreateNewCompanySchema(ctx context.Context, db *bun.DB) error {
	if err := createSchema(ctx, db); err != nil {
		return err
	}

	if err := ChangeSchema(ctx, db); err != nil {
		return err
	}

	if err := createTables(ctx, db); err != nil {
		return err
	}

	return nil
}

func registerModels(db *bun.DB) error {
	db.RegisterModel((*entitymodel.Entity)(nil))

	db.RegisterModel((*model.ProductCategoryToAdditional)(nil))
	db.RegisterModel((*model.ProductCategoryToComplement)(nil))
	db.RegisterModel((*model.ProductToCombo)(nil))
	db.RegisterModel((*model.Size)(nil))
	db.RegisterModel((*model.Quantity)(nil))
	db.RegisterModel((*model.ProductCategory)(nil))
	db.RegisterModel((*model.ProductCategoryWithOrderProcess)(nil))
	db.RegisterModel((*model.ProcessRule)(nil))
	db.RegisterModel((*model.ProcessRuleWithOrderProcess)(nil))
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
	db.RegisterModel((*model.CompanyToUsers)(nil))
	db.RegisterModel((*model.Company)(nil))
	db.RegisterModel((*model.User)(nil))
	var _ bun.BeforeSelectHook = (*model.User)(nil)

	return nil
}

func createTableIfNotExists(ctx context.Context, db *bun.DB, model interface{}) error {
	_, err := db.NewCreateTable().IfNotExists().Model(model).Exec(ctx)
	return err
}

func createTables(ctx context.Context, db *bun.DB) error {
	if err := createTableIfNotExists(ctx, db, (*model.Size)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.ProductCategoryToAdditional)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.ProductCategoryToComplement)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.ProductToCombo)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.ProductCategory)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.Quantity)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.ProcessRule)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.Product)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.Address)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.Contact)(nil)); err != nil {
		return err
	}

	index := "CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_contact ON contacts (ddd, number, type);"

	if _, err := db.ExecContext(ctx, index); err != nil {
		return err
	}

	if err := setupContactFtSearch(ctx, db); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.Client)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.Employee)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.OrderProcess)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.OrderQueue)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.GroupItem)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.OrderProcessToProductToGroupItem)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.Item)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.ItemToAdditional)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.OrderPickup)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.OrderDelivery)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.DeliveryDriver)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.OrderTable)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.PaymentOrder)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.Order)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.Table)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.Place)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.PlaceToTables)(nil)); err != nil {
		return err
	}

	if _, err := db.NewCreateIndex().IfNotExists().Model((*model.PlaceToTables)(nil)).
		Unique().
		Index("idx_place_id_row_and_column").
		Column("place_id", "row", "column").
		Exec(ctx); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.Shift)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.CompanyToUsers)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, db, (*model.Company)(nil)); err != nil {
		return err
	}

	return nil
}
