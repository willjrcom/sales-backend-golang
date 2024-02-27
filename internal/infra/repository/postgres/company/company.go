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

func (r *CompanyRepositoryBun) NewCompany(ctx context.Context, company *companyentity.Company) (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	if _, err = tx.NewInsert().Model(company).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.NewInsert().Model(company.Address).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	companyWithUsers := &companyentity.CompanyWithUsers{
		Entity:                  entity.NewEntity(),
		CompanyCommonAttributes: company.CompanyCommonAttributes,
	}

	if err = database.ChangeToPublicSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err = r.db.NewInsert().Model(companyWithUsers).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *CompanyRepositoryBun) GetCompany(ctx context.Context) (*companyentity.Company, error) {
	company := &companyentity.Company{}
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	err := r.db.NewSelect().Model(company).Relation("Address").Scan(ctx)

	if err != nil {
		return nil, err
	}

	return company, err
}

func (r *CompanyRepositoryBun) AddUserToPublicCompany(ctx context.Context, userID uuid.UUID) error {
	schema := ctx.Value(schemaentity.Schema("schema")).(string)

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		return err
	}

	companyWithUsers := &companyentity.CompanyWithUsers{}
	if err := r.db.NewSelect().Model(companyWithUsers).Where("schema_name = ?", schema).Scan(ctx); err != nil {
		return err
	}

	_, err := r.db.NewInsert().Model(&companyentity.CompanyToUsers{CompanyWithUsersID: companyWithUsers.ID, UserID: userID}).Exec(ctx)

	return err

}

func (r *CompanyRepositoryBun) RemoveUserFromPublicCompany(ctx context.Context, userID uuid.UUID) error {
	schema := ctx.Value(schemaentity.Schema("schema")).(string)

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		return err
	}

	companyWithUsers := &companyentity.CompanyWithUsers{}
	if err := r.db.NewSelect().Model(companyWithUsers).Where("schema_name = ?", schema).Scan(ctx); err != nil {
		return err
	}

	_, err := r.db.NewDelete().Model(&companyentity.CompanyToUsers{}).Where("company_with_users_id = ? AND user_id = ?", companyWithUsers.ID, userID).Exec(ctx)

	return err

}
