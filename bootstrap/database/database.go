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
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

func ConnectDB(ctx context.Context) string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "admin"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "Pass!"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "sales-db"
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbName,
	)
}

var (
	db   *bun.DB
	once sync.Once
)

type tenantTxTimeoutKey struct{}

// WithTenantTransactionTimeout allows callers to customize or disable the default
// transaction timeout. Pass a duration <= 0 to disable the timeout entirely.
func WithTenantTransactionTimeout(ctx context.Context, timeout time.Duration) context.Context {
	return context.WithValue(ctx, tenantTxTimeoutKey{}, timeout)
}

func NewPostgreSQLConnection() *bun.DB {
	once.Do(func() {
		ctx := context.Background()
		connectionParams := ConnectDB(ctx)

		// Connect to database doing a PING
		sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connectionParams), pgdriver.WithTimeout(time.Second*30)))

		// Verifique se o banco de dados já existe.
		if err := sqlDB.PingContext(ctx); err != nil {
			log.Printf("erro ao conectar ao banco de dados: %v", err)
			panic(err)
		}

		// set connection settings
		sqlDB.SetMaxOpenConns(5)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(time.Duration(60) * time.Minute)

		db = bun.NewDB(sqlDB, pgdialect.New())

		if err := registerModels(db); err != nil {
			panic(err)
		}

		// if err := createAllSchemaTables(ctx, dbInstance); err != nil {
		// 	panic(err)
		// }

		if err := createPublicTables(db); err != nil {
			panic(err)
		}

		fmt.Println("Db connected")
	})
	return db
}

func NewSQLiteConnection() *bun.DB {
	// Crie um banco SQLite na memória ou em um arquivo
	db, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}

	// Configure o número máximo de conexões
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Duration(60) * time.Minute)

	// Crie a instância do Bun com o driver SQLite
	dbInstance := bun.NewDB(db, sqlitedialect.New())

	// Registre os modelos necessários
	if err := registerModels(dbInstance); err != nil {
		panic(err)
	}

	// Crie as tabelas no esquema, se necessário
	// ctx := context.Background()
	// if err := createAllSchemaTables(ctx, dbInstance); err != nil {
	// 	panic(err)
	// }

	fmt.Println("SQLite connected")
	return dbInstance
}

func GetTenantTransaction(ctx context.Context, db *bun.DB) (context.Context, *bun.Tx, context.CancelFunc, error) {
	schemaName, err := GetCurrentSchema(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	timeout := 5 * time.Second
	if v := ctx.Value(tenantTxTimeoutKey{}); v != nil {
		if custom, ok := v.(time.Duration); ok {
			timeout = custom
		}
	}

	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		cancel()
		return nil, nil, nil, err
	}

	query := fmt.Sprintf("SET search_path='%s'", schemaName)
	if _, err = tx.ExecContext(ctx, query); err != nil {
		cancel()
		return nil, nil, nil, err
	}

	return ctx, &tx, cancel, nil
}

func GetPublicTenantTransaction(ctx context.Context, db *bun.DB) (context.Context, *bun.Tx, context.CancelFunc, error) {
	ctx = context.WithValue(ctx, model.Schema("schema"), model.PUBLIC_SCHEMA)
	return GetTenantTransaction(ctx, db)
}

func GetCurrentSchema(ctx context.Context) (string, error) {
	schemaName := ctx.Value(model.Schema("schema"))
	if schemaName == nil || schemaName == "" {
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

func createPublicTables(db *bun.DB) error {
	ctx := context.WithValue(context.Background(), model.Schema("schema"), model.PUBLIC_SCHEMA)
	if err := createSchema(ctx, db); err != nil {
		panic(err)
	}

	ctx, tx, cancel, err := GetTenantTransaction(ctx, db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if err := createTableIfNotExists(ctx, tx, (*model.CompanyToUsers)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.CompanyToCategory)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.User)(nil)); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS pgcrypto;"); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.CompanyCategory)(nil)); err != nil {
		return err
	}
	if err := createTableIfNotExists(ctx, tx, (*model.Company)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.Address)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.CompanyPayment)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.CompanyUsageCost)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.CompanySubscription)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.FiscalInvoice)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.FiscalSettings)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.Contact)(nil)); err != nil {
		return err
	}

	// if err := setupPublicMigrations(ctx, tx); err != nil {
	// 	return err
	// }

	if err := tx.Commit(); err != nil {
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

		if !strings.Contains(schemaName, "company_") {
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

	ctx, tx, cancel, err := GetTenantTransaction(ctx, db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if err := createTables(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func registerModels(db *bun.DB) error {
	db.RegisterModel((*entitymodel.Entity)(nil))

	db.RegisterModel((*model.ProductCategoryToAdditional)(nil))
	db.RegisterModel((*model.ProductCategoryToComplement)(nil))
	db.RegisterModel((*model.Size)(nil))
	db.RegisterModel((*model.ProductCategory)(nil))
	db.RegisterModel((*model.ProductCategoryWithOrderProcess)(nil))
	db.RegisterModel((*model.ProcessRule)(nil))
	db.RegisterModel((*model.ProcessRuleWithOrderProcess)(nil))
	db.RegisterModel((*model.Product)(nil))

	db.RegisterModel((*model.Stock)(nil))
	db.RegisterModel((*model.StockMovement)(nil))
	db.RegisterModel((*model.StockAlert)(nil))

	db.RegisterModel((*model.Address)(nil))
	db.RegisterModel((*model.Contact)(nil))
	db.RegisterModel((*model.Client)(nil))
	db.RegisterModel((*model.Employee)(nil))
	db.RegisterModel((*model.EmployeeSalaryHistory)(nil))
	db.RegisterModel((*model.PaymentEmployee)(nil))

	db.RegisterModel((*model.OrderProcessToProductToGroupItem)(nil))
	db.RegisterModel((*model.OrderProcess)(nil))
	db.RegisterModel((*model.OrderQueue)(nil))
	db.RegisterModel((*model.ItemToAdditional)(nil))
	db.RegisterModel((*model.Item)(nil))
	db.RegisterModel((*model.GroupItem)(nil))

	// var _ bun.BeforeUpdateHook = (*model.GroupItem)(nil)

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
	db.RegisterModel((*model.CompanyToCategory)(nil))
	db.RegisterModel((*model.CompanyCategory)(nil))
	db.RegisterModel((*model.Company)(nil))
	db.RegisterModel((*model.CompanyPayment)(nil))
	db.RegisterModel((*model.User)(nil))
	var _ bun.BeforeSelectHook = (*model.User)(nil)

	// Fiscal invoice models
	db.RegisterModel((*model.CompanyUsageCost)(nil))
	db.RegisterModel((*model.CompanySubscription)(nil))
	db.RegisterModel((*model.FiscalInvoice)(nil))
	db.RegisterModel((*model.FiscalSettings)(nil))

	return nil
}

func createTableIfNotExists(ctx context.Context, tx *bun.Tx, model interface{}) error {
	_, err := tx.NewCreateTable().IfNotExists().Model(model).Exec(ctx)
	return err
}

func createTables(ctx context.Context, tx *bun.Tx) error {
	if err := createTableIfNotExists(ctx, tx, (*model.Size)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.ProductCategoryToAdditional)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.ProductCategoryToComplement)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.ProductCategory)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.ProcessRule)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.Product)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.Stock)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.StockMovement)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.StockAlert)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.Address)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.Contact)(nil)); err != nil {
		return err
	}

	index := "CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_contact ON contacts (number, type);"

	if _, err := tx.ExecContext(ctx, index); err != nil {
		return err
	}

	if err := setupContactFtSearch(ctx, tx); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.Client)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.Employee)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.EmployeeSalaryHistory)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.PaymentEmployee)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.OrderProcess)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.OrderQueue)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.GroupItem)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.OrderProcessToProductToGroupItem)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.Item)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.ItemToAdditional)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.OrderPickup)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.OrderDelivery)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.DeliveryDriver)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.OrderTable)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.PaymentOrder)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.Order)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.Table)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.Place)(nil)); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.PlaceToTables)(nil)); err != nil {
		return err
	}

	if _, err := tx.NewCreateIndex().IfNotExists().Model((*model.PlaceToTables)(nil)).
		Unique().
		Index("idx_place_id_row_and_column").
		Column("place_id", "row", "column").
		Exec(ctx); err != nil {
		return err
	}

	if err := createTableIfNotExists(ctx, tx, (*model.Shift)(nil)); err != nil {
		return err
	}

	return nil
}
