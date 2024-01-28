package companyrepositorybun

import (
	"context"
	"database/sql"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

type CompanyRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewCompanyRepositoryBun(db *bun.DB) *CompanyRepositoryBun {
	return &CompanyRepositoryBun{db: db}
}

func (r *CompanyRepositoryBun) NewCompany(ctx context.Context, company *companyentity.Company) error {
	r.mu.Lock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		r.mu.Unlock()
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		r.mu.Unlock()
		return err
	}
	_, err = tx.NewInsert().Model(company).Exec(ctx)
	if err != nil {
		tx.Rollback()
		r.mu.Unlock()
		return err
	}

	_, err = tx.NewInsert().Model(&company.Address).Exec(ctx)
	if err != nil {
		tx.Rollback()
		r.mu.Unlock()
		return err
	}

	if err = tx.Commit(); err != nil {
		r.mu.Unlock()
		return err
	}

	r.mu.Unlock()

	return nil
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
