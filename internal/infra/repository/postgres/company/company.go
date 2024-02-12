package companyrepositorybun

import (
	"context"
	"database/sql"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
)

type CompanyRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewCompanyRepositoryBun(db *bun.DB) *CompanyRepositoryBun {
	return &CompanyRepositoryBun{db: db}
}

func (r *CompanyRepositoryBun) NewCompany(ctx context.Context, company *companyentity.Company) (publicCompanyID uuid.UUID, err error) {
	r.mu.Lock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		r.mu.Unlock()
		return uuid.Nil, err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		r.mu.Unlock()
		return uuid.Nil, err
	}

	if _, err = tx.NewInsert().Model(company).Exec(ctx); err != nil {
		tx.Rollback()
		r.mu.Unlock()
		return uuid.Nil, err
	}

	if _, err = tx.NewInsert().Model(&company.Address).Exec(ctx); err != nil {
		tx.Rollback()
		r.mu.Unlock()
		return uuid.Nil, err
	}

	if err = tx.Commit(); err != nil {
		r.mu.Unlock()
		return uuid.Nil, err
	}

	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), schemaentity.DEFAULT_SCHEMA)
	companyWithUsers := &companyentity.CompanyWithUsers{
		Entity:                  entity.NewEntity(),
		CompanyCommonAttributes: company.CompanyCommonAttributes,
	}

	if err = database.ChangeSchema(ctx, r.db); err != nil {
		r.mu.Unlock()
		return uuid.Nil, err
	}

	if _, err = r.db.NewInsert().Model(companyWithUsers).Exec(ctx); err != nil {
		r.mu.Unlock()
		return uuid.Nil, err
	}

	r.mu.Unlock()

	return companyWithUsers.ID, nil
}

func (r *CompanyRepositoryBun) GetCompany(ctx context.Context) (*companyentity.Company, error) {
	company := &companyentity.Company{}
	r.mu.Lock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		r.mu.Unlock()
		return nil, err
	}

	err := r.db.NewSelect().Model(company).Relation("Address").Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return company, err
}

func (r *CompanyRepositoryBun) AddUser(ctx context.Context, companyID uuid.UUID, userID uuid.UUID) error {
	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), schemaentity.DEFAULT_SCHEMA)

	r.mu.Lock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		r.mu.Unlock()
		return err
	}

	_, err := r.db.NewInsert().Model(&companyentity.CompanyToUsers{CompanyWithUsersID: companyID, UserID: userID}).Exec(ctx)

	r.mu.Unlock()
	return err

}
