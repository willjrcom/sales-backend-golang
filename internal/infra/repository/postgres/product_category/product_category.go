package productcategoryrepositorybun

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ProductCategoryRepositoryBun struct {
	db *bun.DB
}

func NewProductCategoryRepositoryBun(db *bun.DB) model.CategoryRepository {
	return &ProductCategoryRepositoryBun{db: db}
}

func (r *ProductCategoryRepositoryBun) GetComplementProducts(ctx context.Context, categoryID string) ([]model.Product, error) {
	products := []model.Product{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	category := model.ProductCategory{}
	if err := tx.NewSelect().
		Model(&category).
		Where("category.id = ?", categoryID).
		Relation("ComplementCategories").
		Scan(ctx); err != nil {
		return nil, err
	}

	complementCategoryIDs := []uuid.UUID{}
	for _, ac := range category.ComplementCategories {
		complementCategoryIDs = append(complementCategoryIDs, ac.ID)
	}

	if err := tx.NewSelect().
		Model(&products).
		Where("product.category_id IN (?)", bun.In(complementCategoryIDs)).
		Where("product.is_active = ?", true).
		Relation("Size").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductCategoryRepositoryBun) GetAdditionalProducts(ctx context.Context, categoryID string) ([]model.Product, error) {
	products := []model.Product{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	category := model.ProductCategory{}
	if err := tx.NewSelect().
		Model(&category).
		Where("category.id = ?", categoryID).
		Relation("AdditionalCategories").
		Scan(ctx); err != nil {
		return nil, err
	}

	additionalCategoryIDs := []uuid.UUID{}
	for _, ac := range category.AdditionalCategories {
		additionalCategoryIDs = append(additionalCategoryIDs, ac.ID)
	}

	if err := tx.NewSelect().
		Model(&products).
		Where("product.category_id IN (?)", bun.In(additionalCategoryIDs)).
		Where("product.is_active = ?", true).
		Relation("Size").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductCategoryRepositoryBun) CreateCategory(ctx context.Context, cp *model.ProductCategory) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err = tx.NewInsert().Model(cp).Exec(ctx); err != nil {
		return err
	}

	if err := r.updateAdditionalCategories(ctx, tx, cp.ID, cp.AdditionalCategories); err != nil {
		return err
	}

	if err := r.updateComplementCategories(ctx, tx, cp.ID, cp.ComplementCategories); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProductCategoryRepositoryBun) UpdateCategory(ctx context.Context, c *model.ProductCategory) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err = tx.NewUpdate().Model(c).Where("id = ?", c.ID).Exec(ctx); err != nil {
		return err
	}

	if err := r.updateAdditionalCategories(ctx, tx, c.ID, c.AdditionalCategories); err != nil {
		return err
	}

	if err := r.updateComplementCategories(ctx, tx, c.ID, c.ComplementCategories); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProductCategoryRepositoryBun) updateAdditionalCategories(ctx context.Context, tx *bun.Tx, categoryID uuid.UUID, additionalCategories []model.ProductCategory) error {
	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewDelete().Model(&model.ProductCategoryToAdditional{}).Where("category_id = ?", categoryID).Exec(ctx); err != nil {
		return err
	}

	for _, ac := range additionalCategories {
		categoryToAdditional := &model.ProductCategoryToAdditional{
			CategoryID:           categoryID,
			AdditionalCategoryID: ac.ID,
		}

		if _, err := tx.NewInsert().Model(categoryToAdditional).Exec(ctx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *ProductCategoryRepositoryBun) updateComplementCategories(ctx context.Context, tx *bun.Tx, categoryID uuid.UUID, complementCategories []model.ProductCategory) error {
	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewDelete().Model(&model.ProductCategoryToComplement{}).Where("category_id = ?", categoryID).Exec(ctx); err != nil {
		return err
	}

	for _, ac := range complementCategories {
		categoryToComplement := &model.ProductCategoryToComplement{
			CategoryID:           categoryID,
			ComplementCategoryID: ac.ID,
		}

		if _, err := tx.NewInsert().Model(categoryToComplement).Exec(ctx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *ProductCategoryRepositoryBun) DeleteCategory(ctx context.Context, id string) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	// Soft delete: set is_active to false on category only
	isActive := false
	if _, err := tx.NewUpdate().
		Model(&model.ProductCategory{}).
		Set("is_active = ?", isActive).
		Where("id = ?", id).
		Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProductCategoryRepositoryBun) GetCategoryById(ctx context.Context, id string) (*model.ProductCategory, error) {
	category := &model.ProductCategory{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(category).Where("category.id = ?", id).
		Relation("Sizes").
		Relation("Quantities").
		Relation("Products").
		Relation("ProcessRules").
		Relation("AdditionalCategories").
		Relation("ComplementCategories").
		Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return category, nil
}

func (r *ProductCategoryRepositoryBun) GetCategoryByName(ctx context.Context, name string, withRelation bool) (*model.ProductCategory, error) {
	category := &model.ProductCategory{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	query := tx.NewSelect().Model(category).Where("name = ?", name)

	if withRelation {
		query.Relation("Products").Relation("Sizes").Relation("Quantities").Relation("ProcessRules").Relation("AdditionalCategories").Relation("ComplementCategories")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return category, nil
}

func (r *ProductCategoryRepositoryBun) GetAllCategories(ctx context.Context, IDs []uuid.UUID, page int, perPage int, isActive ...bool) ([]model.ProductCategory, error) {
	categories := []model.ProductCategory{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	// Calculate offset
	offset := page * perPage

	// Default to active records (true)
	activeFilter := true
	if len(isActive) > 0 {
		activeFilter = isActive[0]
	}

	// load categories with their simple relations
	query := tx.NewSelect().
		Model(&categories).
		Relation("Sizes").
		Relation("Quantities").
		Relation("ProcessRules").
		Relation("AdditionalCategories").
		Relation("ComplementCategories").
		Limit(perPage).
		Offset(offset)

	if len(IDs) > 0 {
		query.Where("category.id IN (?) AND category.is_active = ?", bun.In(IDs), activeFilter)
	} else {
		query.Where("category.is_active = ?", activeFilter)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	// fetch products for all categories and their sizes
	// collect category IDs
	categoryIDs := make([]uuid.UUID, len(categories))
	for i, cat := range categories {
		categoryIDs[i] = cat.ID
	}
	// load products with Size relation
	var products []model.Product
	if len(categoryIDs) > 0 {
		if err := tx.NewSelect().
			Model(&products).
			Relation("Size").
			Where("product.category_id IN (?)", bun.In(categoryIDs)).
			Where("product.is_active = ?", activeFilter).
			Scan(ctx); err != nil {
			return nil, err
		}
	}
	// group products by their category
	prodByCat := make(map[uuid.UUID][]model.Product, len(categories))
	for _, p := range products {
		prodByCat[p.CategoryID] = append(prodByCat[p.CategoryID], p)
	}
	// assign grouped products to categories
	for i := range categories {
		if ps, ok := prodByCat[categories[i].ID]; ok {
			categories[i].Products = ps
		} else {
			categories[i].Products = nil
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ProductCategoryRepositoryBun) GetAllCategoriesWithProcessRulesAndOrderProcess(ctx context.Context) ([]model.ProductCategoryWithOrderProcess, error) {
	categories := []model.ProductCategoryWithOrderProcess{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	// Busca todas as categorias com suas ProcessRules
	err = tx.NewSelect().
		Model(&categories).
		Relation("ProcessRules").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	schemaName, err := database.GetCurrentSchema(ctx)
	if err != nil {
		return nil, err
	}

	processRuleIDs := make([]uuid.UUID, 0)

	// Coletar todos os IDs de ProcessRules
	for i := range categories {
		for j := range categories[i].ProcessRules {
			processRule := categories[i].ProcessRules[j]
			processRuleIDs = append(processRuleIDs, processRule.ID)
		}
	}

	// Query SQL para contar os processos em geral
	rows, err := r.db.QueryContext(ctx, `
		SELECT process_rule_id, 
			COUNT(CASE WHEN status NOT IN ('Finished', 'Canceled') THEN 1 END) AS total_orders, 
			COUNT(CASE WHEN status NOT IN ('Finished', 'Canceled') AND (EXTRACT(EPOCH FROM (NOW() - started_at::timestamptz)) * 1000000000) > pr.ideal_time THEN 1 END) AS late_orders
		FROM `+schemaName+`.order_processes AS process
		JOIN `+schemaName+`.process_rules AS pr ON process.process_rule_id = pr.id
		WHERE process_rule_id IN (?) 
		GROUP BY process_rule_id
	`, bun.In(processRuleIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map para armazenar os resultados
	processCount := map[uuid.UUID]int{}
	lateCount := map[uuid.UUID]int{}

	for rows.Next() {
		var processRuleID uuid.UUID
		var totalOrders int
		var lateOrders int

		if err := rows.Scan(&processRuleID, &totalOrders, &lateOrders); err != nil {
			return nil, err
		}

		processCount[processRuleID] = totalOrders
		lateCount[processRuleID] = lateOrders
	}

	// Preenchimento na struct
	for i := range categories {
		for j := range categories[i].ProcessRules {
			processRule := &categories[i].ProcessRules[j]
			if count, ok := processCount[processRule.ID]; ok {
				processRule.TotalOrderQueue = count // Total de pedidos
			}
			if late, ok := lateCount[processRule.ID]; ok {
				processRule.TotalOrderProcessLate = late // Total de atrasados
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ProductCategoryRepositoryBun) GetAllCategoriesMap(ctx context.Context, isActive bool) ([]model.ProductCategory, error) {
	categories := []model.ProductCategory{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	err = tx.NewSelect().
		Model(&categories).
		Column("id", "name").
		Where("is_active = ?", isActive).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ProductCategoryRepositoryBun) GetComplementCategories(ctx context.Context) ([]model.ProductCategory, error) {
	categories := []model.ProductCategory{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	err = tx.NewSelect().
		Model(&categories).
		Column("id", "name").
		Where("is_complement = ?", true).
		Where("is_active = ?", true).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ProductCategoryRepositoryBun) GetAdditionalCategories(ctx context.Context) ([]model.ProductCategory, error) {
	categories := []model.ProductCategory{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	err = tx.NewSelect().
		Model(&categories).
		Column("id", "name").
		Where("is_additional = ?", true).
		Where("is_active = ?", true).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return categories, nil
}
