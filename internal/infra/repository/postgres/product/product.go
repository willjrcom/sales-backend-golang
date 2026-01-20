package productrepositorybun

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ProductRepositoryBun struct {
	db *bun.DB
}

func NewProductRepositoryBun(db *bun.DB) model.ProductRepository {
	return &ProductRepositoryBun{db: db}
}

func (r *ProductRepositoryBun) CreateProduct(ctx context.Context, p *model.Product) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(p).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProductRepositoryBun) UpdateProduct(ctx context.Context, p *model.Product) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().Model(p).Where("id = ?", p.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProductRepositoryBun) DeleteProduct(ctx context.Context, id string) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewDelete().Model(&model.Product{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *ProductRepositoryBun) GetProductById(ctx context.Context, id string) (*model.Product, error) {
	product := &model.Product{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(product).Where("product.id = ?", id).Relation("Category").Relation("Size").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductRepositoryBun) GetProductByCode(ctx context.Context, code string) (*model.Product, error) {
	product := &model.Product{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(product).Where("product.code = ?", code).Relation("Category").Relation("Size").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductRepositoryBun) GetAllProducts(ctx context.Context, page, perPage int, isActive bool, categoryID string) ([]model.Product, int, error) {
	products := []model.Product{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	defer cancel()
	defer tx.Rollback()

	// Calculate offset
	offset := page * perPage

	// Get paginated products with filter
	query := tx.NewSelect().
		Model(&products).
		Relation("Category").
		Relation("Size").
		Where("product.is_active = ?", isActive).
		Order("product.name ASC").
		Limit(perPage).
		Offset(offset)

	if categoryID != "" {
		query = query.Where("product.category_id = ?", categoryID)
	}
	if err := query.Scan(ctx); err != nil {
		return nil, 0, err
	}

	// Get total count
	total, err := tx.NewSelect().
		Model(&model.Product{}).
		Where("is_active = ?", isActive).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *ProductRepositoryBun) GetDefaultProducts(ctx context.Context, page, perPage int, isActive bool) ([]model.Product, int, error) {
	products := []model.Product{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, 0, err
	}

	defer cancel()
	defer tx.Rollback()

	// Calculate offset
	offset := page * perPage

	// Get paginated products with filter
	if err := tx.NewSelect().
		Model(&products).
		Relation("Category").
		Relation("Size").
		Join("JOIN product_categories AS cat ON cat.id = product.category_id").
		Where("product.is_active = ?", isActive).
		Where("cat.is_additional = ?", false).
		Where("cat.is_complement = ?", false).
		Order("product.name ASC").
		Limit(perPage).
		Offset(offset).
		Scan(ctx); err != nil {
		return nil, 0, err
	}

	// Get total count
	total, err := tx.NewSelect().
		Model(&model.Product{}).
		Join("JOIN product_categories AS cat ON cat.id = product.category_id").
		Where("product.is_active = ?", isActive).
		Where("cat.is_additional = ?", false).
		Where("cat.is_complement = ?", false).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(); err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (r *ProductRepositoryBun) GetAllProductsMap(ctx context.Context, isActive bool, categoryID string) ([]model.Product, error) {
	products := []model.Product{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	query := tx.NewSelect().
		Model(&products).
		Relation("Size").
		Column("product.id", "product.name", "product.size_id").
		Where("product.is_active = ?", isActive)

	if categoryID != "" {
		query = query.Where("product.category_id = ?", categoryID)
	}

	err = query.Scan(ctx)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return products, nil
}
