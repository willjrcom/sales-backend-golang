package orderrepositorybun

import (
	"context"
	"database/sql"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type OrderRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewOrderRepositoryBun(db *bun.DB) model.OrderRepository {
	return &OrderRepositoryBun{db: db}
}

func (r *OrderRepositoryBun) CreateOrder(ctx context.Context, order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewInsert().Model(order).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *OrderRepositoryBun) PendingOrder(ctx context.Context, p *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	if _, err = tx.NewUpdate().Model(p).Where("id = ?", p.ID).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	for _, group := range p.GroupItems {
		if _, err = tx.NewUpdate().Model(&group).WherePK().Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}

		for _, item := range group.Items {
			if _, err = tx.NewUpdate().Model(&item).WherePK().Exec(ctx); err != nil {
				tx.Rollback()
				return err
			}

			for _, additionalItem := range item.AdditionalItems {
				if _, err = tx.NewUpdate().Model(&additionalItem).WherePK().Exec(ctx); err != nil {
					tx.Rollback()
					return err
				}
			}

			if group.ComplementItemID != nil && group.ComplementItem != nil {
				if _, err = tx.NewUpdate().Model(group.ComplementItem).WherePK().Exec(ctx); err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	if p.Delivery != nil {
		if _, err = tx.NewUpdate().Model(p.Delivery).WherePK().Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}

	} else if p.Pickup != nil {
		if _, err = tx.NewUpdate().Model(p.Pickup).WherePK().Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}

	} else if p.Table != nil {
		if _, err = tx.NewUpdate().Model(p.Table).WherePK().Exec(ctx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *OrderRepositoryBun) UpdateOrder(ctx context.Context, order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(order).Where("id = ?", order.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *OrderRepositoryBun) DeleteOrder(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	if _, err := tx.NewDelete().Model(&model.Order{}).Where("id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.OrderDelivery{}).Where("order_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.OrderPickup{}).Where("order_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.OrderTable{}).Where("order_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.PaymentOrder{}).Where("order_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	groupItems := []model.GroupItem{}
	if err := tx.NewSelect().Model(&groupItems).Where("order_id = ?", id).Relation("ComplementItem").Relation("Items.AdditionalItems").Scan(ctx); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.NewDelete().Model(&model.GroupItem{}).Where("order_id = ?", id).Exec(ctx); err != nil {
		tx.Rollback()
		return err
	}

	for _, groupItem := range groupItems {
		if groupItem.ComplementItem != nil {
			if _, err := tx.NewDelete().Model(groupItem.ComplementItem).WherePK().Exec(ctx); err != nil {
				tx.Rollback()
				return err
			}
		}

		for _, item := range groupItem.Items {
			if _, err := tx.NewDelete().Model(&item).WherePK().Exec(ctx); err != nil {
				tx.Rollback()
				return err
			}

			if _, err := tx.NewDelete().Model(&model.ItemToAdditional{}).Where("item_id = ?", item.ID).Exec(ctx); err != nil {
				tx.Rollback()
				return err
			}

			for _, additionalItem := range item.AdditionalItems {
				if _, err := tx.NewDelete().Model(&additionalItem).WherePK().Exec(ctx); err != nil {
					tx.Rollback()
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

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(order).WherePK().
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

	var complementItems []model.Item
	var complementIDs []uuid.UUID
	for _, g := range order.GroupItems {
		if g.ComplementItemID != nil {
			complementIDs = append(complementIDs, *g.ComplementItemID)
		}
	}
	if len(complementIDs) > 0 {
		if err := r.db.NewSelect().Model(&complementItems).
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
		if err := r.db.NewSelect().Model(order.Delivery).WherePK().
			Relation("Client.Contact").
			Relation("Client.Address").
			Relation("Address").
			Relation("Driver").
			Scan(ctx); err != nil {
			return nil, err
		}
	}

	return order, nil
}

func (r *OrderRepositoryBun) GetAllOrders(ctx context.Context, shiftID string, withStatus []orderentity.StatusOrder) ([]model.Order, error) {
	orders := []model.Order{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	query := r.db.NewSelect().Model(&orders).
		// use quoted alias for reserved keyword 'order'
		Where(`"order"."status" IN (?) OR "order"."shift_id" = ?`, bun.In(withStatus), shiftID).
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
		if err := r.db.NewSelect().Model(&complementItems).
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

	return orders, nil
}

func (r *OrderRepositoryBun) GetAllOrdersWithDelivery(ctx context.Context, shiftID string, page, perPage int) ([]model.Order, error) {
	orders := []model.Order{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	validStatuses := []string{
		string(orderentity.OrderStatusStaging),
		string(orderentity.OrderStatusPending),
		string(orderentity.OrderStatusReady),
	}

	query := r.db.NewSelect().Model(&orders).
		Where("delivery.id IS NOT NULL").
		Relation("Delivery.Client").
		Relation("Delivery.Address").
		Relation("Delivery.Driver").
		Where(`"order"."status" IN (?) OR "order"."shift_id" = ?`, bun.In(validStatuses), shiftID).
		Limit(perPage).
		Offset(page * perPage)

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepositoryBun) AddPaymentOrder(ctx context.Context, payment *model.PaymentOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewInsert().Model(payment).Exec(ctx); err != nil {
		return err
	}

	return nil
}
