package orderrepositorybun

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderRepositoryBun struct {
	db *bun.DB
}

func NewOrderRepositoryBun(db *bun.DB) model.OrderRepository {
	return &OrderRepositoryBun{db: db}
}

func (r *OrderRepositoryBun) CreateOrder(ctx context.Context, order *model.Order) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(order).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderRepositoryBun) PendingOrder(ctx context.Context, p *model.Order) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err = tx.NewUpdate().Model(p).Where("id = ?", p.ID).Exec(ctx); err != nil {
		return err
	}

	for _, group := range p.GroupItems {
		if _, err = tx.NewUpdate().Model(&group).WherePK().Exec(ctx); err != nil {
			return err
		}

		for _, item := range group.Items {
			if _, err = tx.NewUpdate().Model(&item).WherePK().Exec(ctx); err != nil {
				return err
			}

			for _, additionalItem := range item.AdditionalItems {
				if _, err = tx.NewUpdate().Model(&additionalItem).WherePK().Exec(ctx); err != nil {
					return err
				}
			}

			if group.ComplementItemID != nil && group.ComplementItem != nil {
				if _, err = tx.NewUpdate().Model(group.ComplementItem).WherePK().Exec(ctx); err != nil {
					return err
				}
			}
		}
	}

	if p.Delivery != nil {
		if _, err = tx.NewUpdate().Model(p.Delivery).WherePK().Exec(ctx); err != nil {
			return err
		}

	} else if p.Pickup != nil {
		if _, err = tx.NewUpdate().Model(p.Pickup).WherePK().Exec(ctx); err != nil {
			return err
		}

	} else if p.Table != nil {
		if _, err = tx.NewUpdate().Model(p.Table).WherePK().Exec(ctx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *OrderRepositoryBun) UpdateOrder(ctx context.Context, order *model.Order) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewUpdate().Model(order).Where("id = ?", order.ID).Exec(ctx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *OrderRepositoryBun) DeleteOrder(ctx context.Context, id string) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return err
	}

	defer cancel()
	defer tx.Rollback()

	if _, err := tx.NewDelete().Model(&model.Order{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.OrderDelivery{}).Where("order_id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.OrderPickup{}).Where("order_id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.OrderTable{}).Where("order_id = ?", id).Exec(ctx); err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.PaymentOrder{}).Where("order_id = ?", id).Exec(ctx); err != nil {
		return err
	}

	groupItems := []model.GroupItem{}
	if err := tx.NewSelect().Model(&groupItems).Where("order_id = ?", id).Relation("ComplementItem").Relation("Items.AdditionalItems").Scan(ctx); err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.GroupItem{}).Where("order_id = ?", id).Exec(ctx); err != nil {
		return err
	}

	for _, groupItem := range groupItems {
		if groupItem.ComplementItem != nil {
			if _, err := tx.NewDelete().Model(groupItem.ComplementItem).WherePK().Exec(ctx); err != nil {
				return err
			}
		}

		for _, item := range groupItem.Items {
			if _, err := tx.NewDelete().Model(&item).WherePK().Exec(ctx); err != nil {
				return err
			}

			if _, err := tx.NewDelete().Model(&model.ItemToAdditional{}).Where("item_id = ?", item.ID).Exec(ctx); err != nil {
				return err
			}

			for _, additionalItem := range item.AdditionalItems {
				if _, err := tx.NewDelete().Model(&additionalItem).WherePK().Exec(ctx); err != nil {
					return err
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *OrderRepositoryBun) GetOrderById(ctx context.Context, id string) (order *model.Order, err error) {
	order = &model.Order{}
	order.ID = uuid.MustParse(id)

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if err := tx.NewSelect().Model(order).WherePK().
		Relation("GroupItems.Items", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("is_additional = ?", false)
		}).
		Relation("GroupItems.Items.AdditionalItems").
		Relation("Attendant").
		Relation("Payments").
		Relation("Table").
		Relation("Delivery").
		Relation("Pickup").
		Scan(ctx); err != nil {
		return nil, err
	}
	// load complement items separately
	var complementItems []model.Item
	var complementIDs []uuid.UUID
	for _, g := range order.GroupItems {
		if g.ComplementItemID != nil {
			complementIDs = append(complementIDs, *g.ComplementItemID)
		}
	}
	if len(complementIDs) > 0 {
		if err := tx.NewSelect().Model(&complementItems).
			Where("id IN (?)", bun.In(complementIDs)).
			Scan(ctx); err != nil {
			return nil, err
		}
		compMap := make(map[uuid.UUID]*model.Item, len(complementItems))
		for i := range complementItems {
			ci := complementItems[i]
			compMap[ci.ID] = &ci
		}
		for i := range order.GroupItems {
			g := &order.GroupItems[i]
			if g.ComplementItemID != nil {
				if ci, ok := compMap[*g.ComplementItemID]; ok {
					g.ComplementItem = ci
				}
			}
		}
	}

	if order.Delivery != nil {
		if err := tx.NewSelect().Model(order.Delivery).WherePK().
			Relation("Client.Contact").
			Relation("Client.Address").
			Relation("Address").
			Relation("Driver").
			Scan(ctx); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return order, nil
}

func (r *OrderRepositoryBun) GetAllOrders(ctx context.Context, shiftID string, withStatus []orderentity.StatusOrder, withCategory bool, queryCondition string) ([]model.Order, error) {
	orders := []model.Order{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	if queryCondition == "" {
		queryCondition = "OR"
	}

	query := tx.NewSelect().Model(&orders).
		// use quoted alias for reserved keyword 'order'
		Where(`"order"."status" IN (?) `+queryCondition+` "order"."shift_id" = ?`, bun.In(withStatus), shiftID).
		Relation("GroupItems.Items", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("is_additional = ?", false)
		}).
		Relation("GroupItems.Items.AdditionalItems").
		Relation("Attendant").
		Relation("Payments").
		Relation("Table").
		Relation("Delivery").
		Relation("Pickup")

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	var complementItems []model.Item
	var complementIDs []uuid.UUID
	for i := range orders {
		for j := range orders[i].GroupItems {
			if orders[i].GroupItems[j].ComplementItemID != nil {
				complementIDs = append(complementIDs, *orders[i].GroupItems[j].ComplementItemID)
			}
		}
	}
	if len(complementIDs) > 0 {
		if err := tx.NewSelect().Model(&complementItems).
			Where("id IN (?)", bun.In(complementIDs)).
			Scan(ctx); err != nil {
			return nil, err
		}
		compMap := make(map[uuid.UUID]*model.Item, len(complementItems))
		for k := range complementItems {
			ci := complementItems[k]
			compMap[ci.ID] = &ci
		}
		for i := range orders {
			for j := range orders[i].GroupItems {
				g := &orders[i].GroupItems[j]
				if g.ComplementItemID != nil {
					if ci, ok := compMap[*g.ComplementItemID]; ok {
						g.ComplementItem = ci
					}
				}
			}
		}
	}
	// optionally load categories if requested
	if withCategory {
		var categories []model.ProductCategory
		var categoryIDs []uuid.UUID
		for i := range orders {
			for j := range orders[i].GroupItems {
				categoryIDs = append(categoryIDs, orders[i].GroupItems[j].CategoryID)
			}
		}
		if len(categoryIDs) > 0 {
			if err := tx.NewSelect().Model(&categories).
				Where("id IN (?)", bun.In(categoryIDs)).
				Scan(ctx); err != nil {
				return nil, err
			}
			catMap := make(map[uuid.UUID]*model.ProductCategory, len(categories))
			for k := range categories {
				ci := categories[k]
				catMap[ci.ID] = &ci
			}
			for i := range orders {
				for j := range orders[i].GroupItems {
					g := &orders[i].GroupItems[j]
					if cat, ok := catMap[g.CategoryID]; ok {
						g.Category = cat
					}
				}
			}
		}
	}

	// Coletar todos os DriverIDs dos deliveries
	var driverIDs []uuid.UUID
	for i := range orders {
		if orders[i].Delivery != nil && orders[i].Delivery.DriverID != nil {
			driverIDs = append(driverIDs, *orders[i].Delivery.DriverID)
		}
	}

	// Buscar todos os drivers em uma query separada
	var drivers []model.DeliveryDriver
	if len(driverIDs) > 0 {
		if err := tx.NewSelect().Model(&drivers).
			Where("id IN (?)", bun.In(driverIDs)).
			Scan(ctx); err != nil {
			return nil, err
		}
	}

	// Mapear drivers por ID
	driverMap := make(map[uuid.UUID]*model.DeliveryDriver, len(drivers))
	for i := range drivers {
		d := drivers[i]
		driverMap[d.ID] = &d
	}
	for i := range orders {
		if orders[i].Delivery != nil && orders[i].Delivery.DriverID != nil {
			if driver, ok := driverMap[*orders[i].Delivery.DriverID]; ok {
				orders[i].Delivery.Driver = driver
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepositoryBun) GetAllOrdersWithDelivery(ctx context.Context, shiftID string, page, perPage int) ([]model.Order, error) {
	orders := []model.Order{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	validStatuses := []orderentity.StatusOrder{
		orderentity.OrderStatusReady,
	}

	query := tx.NewSelect().Model(&orders).
		Where("delivery.id IS NOT NULL").
		Relation("Delivery.Client").
		Relation("Delivery.Address").
		Where(`"order"."status" IN (?) OR "order"."shift_id" = ?`, bun.In(validStatuses), shiftID).
		Limit(perPage).
		Offset(page * perPage)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	// Coletar todos os DriverIDs
	var driverIDs []uuid.UUID
	for i := range orders {
		if orders[i].Delivery != nil && orders[i].Delivery.DriverID != nil {
			driverIDs = append(driverIDs, *orders[i].Delivery.DriverID)
		}
	}

	// Buscar todos os drivers em uma query separada
	var drivers []model.DeliveryDriver
	if len(driverIDs) > 0 {
		if err := tx.NewSelect().Model(&drivers).
			Where("id IN (?)", bun.In(driverIDs)).
			Scan(ctx); err != nil {
			return nil, err
		}
	}

	// Mapear drivers por ID
	driverMap := make(map[uuid.UUID]*model.DeliveryDriver, len(drivers))
	for i := range drivers {
		d := drivers[i]
		driverMap[d.ID] = &d
	}
	for i := range orders {
		if orders[i].Delivery != nil && orders[i].Delivery.DriverID != nil {
			if driver, ok := driverMap[*orders[i].Delivery.DriverID]; ok {
				orders[i].Delivery.Driver = driver
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepositoryBun) GetAllOrdersWithPickup(ctx context.Context, shiftID string, page, perPage int) ([]model.Order, error) {
	orders := []model.Order{}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
	if err != nil {
		return nil, err
	}

	defer cancel()
	defer tx.Rollback()

	validStatuses := []string{
		string(orderentity.OrderStatusReady),
	}

	query := tx.NewSelect().Model(&orders).
		Relation("Pickup").
		Where("pickup.id IS NOT NULL").
		Where(`"order"."status" IN (?) AND "order"."shift_id" = ?`, bun.In(validStatuses), shiftID).
		Limit(perPage).
		Offset(page * perPage)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepositoryBun) AddPaymentOrder(ctx context.Context, payment *model.PaymentOrder) error {

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, r.db)
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
