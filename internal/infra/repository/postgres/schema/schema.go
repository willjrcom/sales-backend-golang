package schemarepositorybun

import (
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
	"golang.org/x/net/context"
)

type SchemaRepositoryBun struct {
	db *bun.DB
}

func NewSchemaRepositoryBun(db *bun.DB) *SchemaRepositoryBun {
	return &SchemaRepositoryBun{db: db}
}

func (r *SchemaRepositoryBun) NewSchema(ctx context.Context) error {
	if err := loadCompanyModels(ctx, r.db); err != nil {
		return err
	}

	if err := setupFtSearch(ctx, r.db); err != nil {
		return err
	}

	return nil
}

func loadCompanyModels(ctx context.Context, db *bun.DB) error {
	mu := sync.Mutex{}

	mu.Lock()
	if err := database.CreateSchema(ctx, db); err != nil {
		mu.Unlock()
		return err
	}

	if err := database.ChangeSchema(ctx, db); err != nil {
		mu.Unlock()
		return err
	}

	db.RegisterModel((*entity.Entity)(nil))

	db.RegisterModel((*productentity.Size)(nil))
	db.RegisterModel((*productentity.Quantity)(nil))
	db.RegisterModel((*productentity.CategoryToAdditionalCategories)(nil))
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
	db.RegisterModel((*companyentity.Company)(nil))

	if _, err := db.NewCreateTable().IfNotExists().Model((*entity.Entity)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.Size)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.CategoryToAdditionalCategories)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.Category)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.Quantity)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.Process)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*productentity.Product)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*addressentity.Address)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*personentity.Contact)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*personentity.Person)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*cliententity.Client)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*employeeentity.Employee)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*itementity.Item)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*groupitementity.GroupItem)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*orderentity.DeliveryOrder)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*orderentity.TableOrder)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*orderentity.PaymentOrder)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*orderentity.Order)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*tableentity.Table)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*shiftentity.Shift)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*companyentity.Company)(nil)).Exec(ctx); err != nil {
		mu.Unlock()
		return err
	}

	return nil
}
