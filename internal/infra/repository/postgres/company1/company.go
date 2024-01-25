package companyrepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
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
	_, err := r.db.NewInsert().Model(company).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *CompanyRepositoryBun) GetCompanyById(ctx context.Context, id string) (*companyentity.Company, error) {
	company := &companyentity.Company{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(company).Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return company, err
}

func (r *CompanyRepositoryBun) GetAllCompaniesBySchemaName(ctx context.Context, schemaName string) ([]companyentity.Company, error) {
	companies := []companyentity.Company{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(&companies).Where("schema_name = ?", schemaName).Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return companies, nil
}
