package productcategoryrepositorybun

import (
	"context"
	"database/sql"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ProductCategoryRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewProductCategoryRepositoryBun(db *bun.DB) model.CategoryRepository {
	return &ProductCategoryRepositoryBun{db: db}
}

func (r *ProductCategoryRepositoryBun) CreateCategory(ctx context.Context, cp *model.ProductCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	if _, err = tx.NewInsert().Model(cp).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.updateAdditionalCategories(ctx, &tx, cp.ID, cp.AdditionalCategories); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.updateComplementCategories(ctx, &tx, cp.ID, cp.ComplementCategories); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProductCategoryRepositoryBun) UpdateCategory(ctx context.Context, c *model.ProductCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	if _, err = tx.NewUpdate().Model(c).Where("id = ?", c.ID).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.updateAdditionalCategories(ctx, &tx, c.ID, c.AdditionalCategories); err != nil {
		tx.Rollback()
		return err
	}

	if err := r.updateComplementCategories(ctx, &tx, c.ID, c.ComplementCategories); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *ProductCategoryRepositoryBun) updateAdditionalCategories(ctx context.Context, tx *bun.Tx, categoryID uuid.UUID, additionalCategories []model.ProductCategory) error {
	if err := database.ChangeSchema(ctx, r.db); err != nil {

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

	return nil
}

func (r *ProductCategoryRepositoryBun) updateComplementCategories(ctx context.Context, tx *bun.Tx, categoryID uuid.UUID, complementCategories []model.ProductCategory) error {
	if err := database.ChangeSchema(ctx, r.db); err != nil {

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

	return nil
}

func (r *ProductCategoryRepositoryBun) DeleteCategory(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

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

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(category).Where("id = ?", id).Relation("Products").Relation("Sizes").Relation("Quantities").Relation("ProcessRules").Relation("AdditionalCategories").Relation("ComplementCategories").Scan(ctx); err != nil {
		return nil, err
	}

	return category, nil
}

func (r *ProductCategoryRepositoryBun) GetCategoryByName(ctx context.Context, name string, withRelation bool) (*model.ProductCategory, error) {
	category := &model.ProductCategory{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	query := r.db.NewSelect().Model(category).Where("name = ?", name)

	if withRelation {
		query.Relation("Products").Relation("Sizes").Relation("Quantities").Relation("ProcessRules").Relation("AdditionalCategories").Relation("ComplementCategories")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return category, nil
}

func (r *ProductCategoryRepositoryBun) GetAllCategories(ctx context.Context) ([]model.ProductCategory, error) {
	categories := []model.ProductCategory{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&categories).
		Relation("Products").
		Relation("Sizes").
		Relation("Quantities").
		Relation("ProcessRules").
		Relation("AdditionalCategories").
		Relation("ComplementCategories").
		Scan(ctx); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *ProductCategoryRepositoryBun) GetAllCategoriesWithProcessRulesAndOrderProcess(ctx context.Context) ([]model.ProductCategoryWithOrderProcess, error) {
	categories := []model.ProductCategoryWithOrderProcess{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	// Busca todas as categorias com suas ProcessRules
	err := r.db.NewSelect().
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
			COUNT(CASE WHEN status != 'Finished' THEN 1 END) AS total_orders, 
			COUNT(CASE WHEN status != 'Finished' AND (EXTRACT(EPOCH FROM (NOW() - started_at::timestamptz)) * 1000000000) > pr.ideal_time THEN 1 END) AS late_orders
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

	return categories, nil
}
