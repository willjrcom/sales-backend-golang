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

func (r *ProductCategoryRepositoryBun) CreateCategory(ctx context.Context, cp *model.ProductCategory) error {

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err = tx.NewInsert().Model(cp).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.updateAdditionalCategories(ctx, tx, cp.ID, cp.AdditionalCategories); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.updateComplementCategories(ctx, tx, cp.ID, cp.ComplementCategories); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProductCategoryRepositoryBun) UpdateCategory(ctx context.Context, c *model.ProductCategory) error {

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err = tx.NewUpdate().Model(c).Where("id = ?", c.ID).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.updateAdditionalCategories(ctx, tx, c.ID, c.AdditionalCategories); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.updateComplementCategories(ctx, tx, c.ID, c.ComplementCategories); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProductCategoryRepositoryBun) updateAdditionalCategories(ctx context.Context, tx *bun.Tx, categoryID uuid.UUID, additionalCategories []model.ProductCategory) error {
	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.ProductCategoryToAdditional{}).Where("category_id = ?", categoryID).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	for _, ac := range additionalCategories {
		categoryToAdditional := &model.ProductCategoryToAdditional{
			CategoryID:           categoryID,
			AdditionalCategoryID: ac.ID,
		}

		if _, err := tx.NewInsert().Model(categoryToAdditional).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *ProductCategoryRepositoryBun) updateComplementCategories(ctx context.Context, tx *bun.Tx, categoryID uuid.UUID, complementCategories []model.ProductCategory) error {
	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.ProductCategoryToComplement{}).Where("category_id = ?", categoryID).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	for _, ac := range complementCategories {
		categoryToComplement := &model.ProductCategoryToComplement{
			CategoryID:           categoryID,
			ComplementCategoryID: ac.ID,
		}

		if _, err := tx.NewInsert().Model(categoryToComplement).Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *ProductCategoryRepositoryBun) DeleteCategory(ctx context.Context, id string) error {

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.ProductCategory{}).Where("id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.ProductCategoryToAdditional{}).Where("category_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.ProductCategoryToAdditional{}).Where("additional_category_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.ProductCategoryToComplement{}).Where("complement_category_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.Size{}).Where("category_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.Quantity{}).Where("category_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.ProcessRule{}).Where("category_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *ProductCategoryRepositoryBun) GetCategoryById(ctx context.Context, id string) (*model.ProductCategory, error) {
	category := &model.ProductCategory{}

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	if err := tx.NewSelect().Model(category).Where("id = ?", id).Relation("Products").Relation("Sizes").Relation("Quantities").Relation("ProcessRules").Relation("AdditionalCategories").Relation("ComplementCategories").Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return category, nil
}

func (r *ProductCategoryRepositoryBun) GetCategoryByName(ctx context.Context, name string, withRelation bool) (*model.ProductCategory, error) {
	category := &model.ProductCategory{}

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

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

func (r *ProductCategoryRepositoryBun) GetAllCategories(ctx context.Context) ([]model.ProductCategory, error) {
	categories := []model.ProductCategory{}

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	// load categories with their simple relations
	if err := tx.NewSelect().
		Model(&categories).
		Relation("Sizes").
		Relation("Quantities").
		Relation("ProcessRules").
		Relation("AdditionalCategories").
		Relation("ComplementCategories").
		Scan(ctx); err != nil {
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

	tx, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

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
