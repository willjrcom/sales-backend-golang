package companyrepositorybun

import (
	"context"
	"database/sql"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type CompanyRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewCompanyRepositoryBun(db *bun.DB) *CompanyRepositoryBun {
	return &CompanyRepositoryBun{db: db}
}

func (r *CompanyRepositoryBun) NewCompany(ctx context.Context, company *model.Company) (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	// Insert to private schema
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

	if err = database.ChangeToPublicSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err = r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	// Insert on public schema
	if _, err = tx.NewInsert().Model(company).Exec(ctx); err != nil {
		return err
	}

	if _, err = tx.NewInsert().Model(company.Address).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CompanyRepositoryBun) GetCompany(ctx context.Context) (*model.Company, error) {
	company := &model.Company{}
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

func (r *CompanyRepositoryBun) ValidateUserToPublicCompany(ctx context.Context, userID uuid.UUID) (bool, error) {
	schema := ctx.Value(model.Schema("schema")).(string)

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		return false, err
	}

	company := &model.Company{}
	if err := r.db.NewSelect().Model(company).Where("schema_name = ?", schema).Scan(ctx); err != nil {
		return false, err
	}

	companyToUsers := &model.CompanyToUsers{}
	if err := r.db.NewSelect().Model(companyToUsers).Where("company_with_users_id = ? AND user_id = ?", company.ID, userID).Scan(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (r *CompanyRepositoryBun) AddUserToPublicCompany(ctx context.Context, userID uuid.UUID) error {
	schema := ctx.Value(model.Schema("schema")).(string)

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		return err
	}

	company := &model.Company{}
	if err := r.db.NewSelect().Model(company).Where("schema_name = ?", schema).Scan(ctx); err != nil {
		return err
	}

	companyToUsers := &model.CompanyToUsers{CompanyID: company.ID, UserID: userID}
	_, err := r.db.NewInsert().Model(companyToUsers).Exec(ctx)

	return err

}

func (r *CompanyRepositoryBun) RemoveUserFromPublicCompany(ctx context.Context, userID uuid.UUID) error {
	schema := ctx.Value(model.Schema("schema")).(string)

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeToPublicSchema(ctx, r.db); err != nil {
		return err
	}

	company := &model.Company{}
	if err := r.db.NewSelect().Model(company).Where("schema_name = ?", schema).Scan(ctx); err != nil {
		return err
	}

	_, err := r.db.NewDelete().Model(&model.CompanyToUsers{}).Where("company_with_users_id = ? AND user_id = ?", company.ID, userID).Exec(ctx)

	return err

}
