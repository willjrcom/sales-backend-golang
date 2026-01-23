package companyrepositorybun

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type CompanyRepositoryBun struct {
	db *bun.DB
}

func NewCompanyRepositoryBun(db *bun.DB) model.CompanyRepository {
	return &CompanyRepositoryBun{db: db}
}

func (r *CompanyRepositoryBun) NewCompany(ctx context.Context, company *model.Company) (err error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	// Insert on public schema
	if _, err = tx.NewInsert().Model(company).Exec(ctx); err != nil {
		return err
	}

	if _, err = tx.NewInsert().Model(company.Address).Exec(ctx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CompanyRepositoryBun) UpdateCompany(ctx context.Context, company *model.Company) (err error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	// Insert on public schema
	if _, err = tx.NewUpdate().Model(company).WherePK().Exec(ctx); err != nil {
		return err
	}

	if _, err := tx.NewUpdate().Model(company.Address).WherePK().Exec(ctx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CompanyRepositoryBun) GetCompany(ctx context.Context) (*model.Company, error) {
	company := &model.Company{}
	schema, ok := ctx.Value(model.Schema("schema")).(string)
	if !ok {
		return nil, fmt.Errorf("company schema not found")
	}

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(company).Where("schema_name = ?", schema).Relation("Address").Relation("Users").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return company, err
}

func (r *CompanyRepositoryBun) GetCompanyOnlyByID(ctx context.Context, id uuid.UUID) (*model.Company, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	company := &model.Company{}
	if err := tx.NewSelect().
		Model(company).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return company, nil
}

func (r *CompanyRepositoryBun) ValidateUserToPublicCompany(ctx context.Context, userID uuid.UUID) (bool, error) {
	schema, ok := ctx.Value(model.Schema("schema")).(string)
	if !ok {
		return false, fmt.Errorf("company schema not found")
	}

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return false, err
	}

	defer cancel()
	defer tx.Rollback()

	company := &model.Company{}
	if err := tx.NewSelect().Model(company).Where("schema_name = ?", schema).Scan(ctx); err != nil {
		return false, err
	}

	companyToUsers := &model.CompanyToUsers{}
	if err := tx.NewSelect().Model(companyToUsers).Where("company_id = ? AND user_id = ?", company.ID, userID).Scan(ctx); err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

func (r *CompanyRepositoryBun) AddUserToPublicCompany(ctx context.Context, userID uuid.UUID) error {
	schema, ok := ctx.Value(model.Schema("schema")).(string)
	if !ok {
		return fmt.Errorf("company schema not found")
	}

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	company := &model.Company{}
	if err := tx.NewSelect().Model(company).Where("schema_name = ?", schema).Scan(ctx); err != nil {
		return err
	}

	companyToUsers := &model.CompanyToUsers{CompanyID: company.ID, UserID: userID}
	if _, err := tx.NewInsert().Model(companyToUsers).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil

}

func (r *CompanyRepositoryBun) RemoveUserFromPublicCompany(ctx context.Context, userID uuid.UUID) error {
	schema, ok := ctx.Value(model.Schema("schema")).(string)
	if !ok {
		return fmt.Errorf("company schema not found")
	}

	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	company := &model.Company{}
	if err := tx.NewSelect().Model(company).Where("schema_name = ?", schema).Scan(ctx); err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.CompanyToUsers{}).Where("company_id = ? AND user_id = ?", company.ID, userID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// GetCompanyUsers retrieves a paginated list of users for the public company and the total count.
func (r *CompanyRepositoryBun) GetCompanyUsers(ctx context.Context, page, perPage int) ([]model.User, int, error) {
	// switch to public schema
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	defer cancel()
	defer tx.Rollback()

	// find company by schema name
	schema, ok := ctx.Value(model.Schema("schema")).(string)
	if !ok {
		return nil, 0, fmt.Errorf("company schema not found")
	}

	company := &model.Company{}
	if err := tx.NewSelect().Model(company).Where("schema_name = ?", schema).Scan(ctx); err != nil {
		return nil, 0, err
	}
	// count total users
	totalCount, err := tx.NewSelect().Model((*model.CompanyToUsers)(nil)).Where("company_id = ?", company.ID).Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	// fetch paginated user records via join
	users := []model.User{}
	// select users and join with company_to_users to filter by company
	if err := tx.NewSelect().
		Model(&users).
		// FROM public.users AS u defined by Model; join membership
		Join("INNER JOIN public.company_to_users ctu ON ctu.user_id = u.id").
		Where("ctu.company_id = ?", company.ID).
		Relation("Contact").
		Relation("Address").
		Limit(perPage).
		Offset(page * perPage).
		Scan(ctx); err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return users, int(totalCount), nil
}

// ListPublicCompanies returns all companies stored in the public schema with basic fields only.
func (r *CompanyRepositoryBun) ListPublicCompanies(ctx context.Context) ([]model.Company, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	companies := []model.Company{}
	if err := tx.NewSelect().
		Model(&companies).
		Column("id", "business_name", "trade_name", "email", "cnpj", "schema_name").
		Order("business_name ASC").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return companies, nil
}

func (r *CompanyRepositoryBun) ListCompaniesForBilling(ctx context.Context) ([]model.Company, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	companies := []model.Company{}
	if err := tx.NewSelect().
		Model(&companies).
		Column("id", "schema_name", "business_name", "trade_name", "subscription_expires_at", "is_blocked").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return companies, nil
}

func (r *CompanyRepositoryBun) UpdateCompanySubscription(ctx context.Context, companyID uuid.UUID, schema string, expiresAt *time.Time, isBlocked bool) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	update := tx.NewUpdate().
		Model((*model.Company)(nil)).
		Set("is_blocked = ?", isBlocked).
		Where("id = ?", companyID)

	if expiresAt == nil {
		update = update.Set("subscription_expires_at = NULL")
	} else {
		update = update.Set("subscription_expires_at = ?", expiresAt)
	}

	if _, err := update.Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *CompanyRepositoryBun) CreateCompanyPayment(ctx context.Context, payment *model.CompanyPayment) error {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(payment).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *CompanyRepositoryBun) GetCompanyPaymentByProviderID(ctx context.Context, provider string, paymentID string) (*model.CompanyPayment, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	payment := &model.CompanyPayment{}
	if err := tx.NewSelect().
		Model(payment).
		Where("provider = ?", provider).
		Where("provider_payment_id = ?", paymentID).
		Limit(1).
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *CompanyRepositoryBun) ListCompanyPayments(ctx context.Context, companyID uuid.UUID, page, perPage int) ([]model.CompanyPayment, int, error) {
	ctx, tx, cancel, err := database.GetPublicTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	defer cancel()
	defer tx.Rollback()

	total, err := tx.NewSelect().
		Model((*model.CompanyPayment)(nil)).
		Where("company_id = ?", companyID).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		if err := tx.Commit(); err != nil {
			return nil, 0, err
		}
		return []model.CompanyPayment{}, 0, nil
	}

	payments := make([]model.CompanyPayment, 0, perPage)

	if err := tx.NewSelect().
		Model(&payments).
		Where("company_id = ?", companyID).
		Order("paid_at DESC").
		Order("created_at DESC").
		Limit(perPage).
		Offset(page * perPage).
		Scan(ctx); err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}
