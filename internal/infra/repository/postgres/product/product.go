package productrepositorybun

import (
	"context"

	"github.com/google/uuid"
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

	for _, v := range p.Variations {
		v.ProductID = p.ID
		if _, err := tx.NewInsert().Model(v).Exec(ctx); err != nil {
			return err
		}
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

	// Fetch existing variations
	var existingVariations []model.ProductVariation
	if err := tx.NewSelect().Model(&existingVariations).Where("product_id = ?", p.ID).Scan(ctx); err != nil {
		return err
	}

	// Map existing variations by SizeID
	existingMap := make(map[uuid.UUID]model.ProductVariation)
	for _, v := range existingVariations {
		existingMap[v.SizeID] = v
	}

	// Lists to handle db operations
	var toInsert []*model.ProductVariation
	var toUpdate []*model.ProductVariation
	processedSizeIDs := make(map[uuid.UUID]bool)

	for _, v := range p.Variations {
		v.ProductID = p.ID
		processedSizeIDs[v.SizeID] = true

		if existing, ok := existingMap[v.SizeID]; ok {
			// Update existing variation
			v.ID = existing.ID
			v.CreatedAt = existing.CreatedAt
			v.DeletedAt = nil // Ensure it's not deleted
			toUpdate = append(toUpdate, v)
		} else {
			// Insert new variation
			if v.ID == uuid.Nil {
				v.ID = uuid.New()
			}
			toInsert = append(toInsert, v)
		}
	}

	// Delete variations that are not in the new list
	var idsToDelete []uuid.UUID
	for _, v := range existingVariations {
		if !processedSizeIDs[v.SizeID] {
			idsToDelete = append(idsToDelete, v.ID)
		}
	}

	if len(idsToDelete) > 0 {
		if _, err := tx.NewDelete().Model((*model.ProductVariation)(nil)).
			Where("id IN (?)", bun.In(idsToDelete)).
			Exec(ctx); err != nil {
			return err
		}
	}

	if len(toInsert) > 0 {
		if _, err := tx.NewInsert().Model(&toInsert).Exec(ctx); err != nil {
			return err
		}
	}

	for _, v := range toUpdate {
		if _, err := tx.NewUpdate().Model(v).
			Column("price", "cost", "is_available", "deleted_at").
			WherePK().
			Exec(ctx); err != nil {
			return err
		}
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

	if err := tx.NewSelect().Model(product).Where("product.id = ?", id).Relation("Category").Relation("Variations.Size").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductRepositoryBun) GetProductBySKU(ctx context.Context, sku string) (*model.Product, error) {
	product := &model.Product{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(product).Where("product.sku = ?", sku).Relation("Category").Relation("Variations.Size").Scan(ctx); err != nil {
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

	// Get paginated products with filter
	query := tx.NewSelect().
		Model(&products).
		Relation("Category").
		Relation("Variations.Size").
		Where("product.is_active = ?", isActive).
		Order("product.name ASC").
		Limit(perPage).
		Offset(page * perPage)

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

	// Get paginated products with filter
	if err := tx.NewSelect().
		Model(&products).
		Relation("Category").
		Relation("Variations.Size").
		Join("JOIN product_categories AS cat ON cat.id = product.category_id").
		Where("product.is_active = ?", isActive).
		Where("cat.is_additional = ?", false).
		Where("cat.is_complement = ?", false).
		Order("product.name ASC").
		Limit(perPage).
		Offset(page * perPage).
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
		Relation("Variations.Size").
		Column("product.id", "product.name").
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
