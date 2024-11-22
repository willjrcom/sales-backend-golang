package orderrepositorybun

import (
	"context"
	"database/sql"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type OrderRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewOrderRepositoryBun(db *bun.DB) *OrderRepositoryBun {
	return &OrderRepositoryBun{db: db}
}

func (r *OrderRepositoryBun) CreateOrder(ctx context.Context, order *orderentity.Order) error {
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

func (r *OrderRepositoryBun) PendingOrder(ctx context.Context, p *orderentity.Order) error {
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

	for _, group := range p.Groups {
		if _, err = tx.NewUpdate().Model(&group).WherePK().Exec(ctx); err != nil {
			if errRoolback := tx.Rollback(); errRoolback != nil {
				return errRoolback
			}

			return err
		}

		for _, item := range group.Items {
			if _, err = tx.NewUpdate().Model(&item).WherePK().Exec(ctx); err != nil {
				if errRoolback := tx.Rollback(); errRoolback != nil {
					return errRoolback
				}

				return err
			}

			for _, additionalItem := range item.AdditionalItems {
				if _, err = tx.NewUpdate().Model(&additionalItem).WherePK().Exec(ctx); err != nil {
					if errRoolback := tx.Rollback(); errRoolback != nil {
						return errRoolback
					}

					return err
				}
			}

			if group.ComplementItemID != nil && group.ComplementItem != nil {
				if _, err = tx.NewUpdate().Model(group.ComplementItem).WherePK().Exec(ctx); err != nil {
					if errRoolback := tx.Rollback(); errRoolback != nil {
						return errRoolback
					}

					return err
				}
			}
		}
	}

	if p.Delivery != nil {
		if _, err = tx.NewUpdate().Model(p.Delivery).WherePK().Exec(ctx); err != nil {
			if errRoolback := tx.Rollback(); errRoolback != nil {
				return errRoolback
			}

			return err
		}

	} else if p.Pickup != nil {
		if _, err = tx.NewUpdate().Model(p.Pickup).WherePK().Exec(ctx); err != nil {
			if errRoolback := tx.Rollback(); errRoolback != nil {
				return errRoolback
			}

			return err
		}

	} else if p.Table != nil {
		if _, err = tx.NewUpdate().Model(p.Table).WherePK().Exec(ctx); err != nil {
			if errRoolback := tx.Rollback(); errRoolback != nil {
				return errRoolback
			}
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		if errRoolback := tx.Rollback(); errRoolback != nil {
			return errRoolback
		}

		return err
	}

	return nil
}

func (r *OrderRepositoryBun) UpdateOrder(ctx context.Context, order *orderentity.Order) error {
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

	if _, err := tx.NewDelete().Model(&orderentity.Order{}).Where("id = ?", id).Exec(ctx); err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return errRollBack
		}

		return err
	}

	if _, err := tx.NewDelete().Model(&orderentity.OrderDelivery{}).Where("order_id = ?", id).Exec(ctx); err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return errRollBack
		}

		return err
	}

	if _, err := tx.NewDelete().Model(&orderentity.OrderPickup{}).Where("order_id = ?", id).Exec(ctx); err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return errRollBack
		}

		return err
	}

	if _, err := tx.NewDelete().Model(&orderentity.OrderTable{}).Where("order_id = ?", id).Exec(ctx); err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return errRollBack
		}

		return err
	}

	if _, err := tx.NewDelete().Model(&orderentity.PaymentOrder{}).Where("order_id = ?", id).Exec(ctx); err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return errRollBack
		}

		return err
	}

	groupItems := []groupitementity.GroupItem{}
	if err := tx.NewSelect().Model(&groupItems).Where("order_id = ?", id).Relation("ComplementItem").Relation("Items.AdditionalItems").Scan(ctx); err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return errRollBack
		}

		return err
	}

	if _, err := tx.NewDelete().Model(&groupitementity.GroupItem{}).Where("order_id = ?", id).Exec(ctx); err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return errRollBack
		}

		return err
	}

	for _, groupItem := range groupItems {
		if groupItem.ComplementItem != nil {
			if _, err := tx.NewDelete().Model(groupItem.ComplementItem).WherePK().Exec(ctx); err != nil {
				if errRollBack := tx.Rollback(); errRollBack != nil {
					return errRollBack
				}

				return err
			}
		}

		for _, item := range groupItem.Items {
			if _, err := tx.NewDelete().Model(&item).WherePK().Exec(ctx); err != nil {
				if errRollBack := tx.Rollback(); errRollBack != nil {
					return errRollBack
				}

				return err
			}

			if _, err := tx.NewDelete().Model(&itementity.ItemToAdditional{}).Where("item_id = ?", item.ID).Exec(ctx); err != nil {
				if errRollBack := tx.Rollback(); errRollBack != nil {
					return errRollBack
				}

				return err
			}

			for _, additionalItem := range item.AdditionalItems {
				if _, err := tx.NewDelete().Model(&additionalItem).WherePK().Exec(ctx); err != nil {
					if errRollBack := tx.Rollback(); errRollBack != nil {
						return errRollBack
					}

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

func (r *OrderRepositoryBun) GetOrderById(ctx context.Context, id string) (order *orderentity.Order, err error) {
	order = &orderentity.Order{}
	order.ID = uuid.MustParse(id)

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(order).WherePK().Relation("Groups.Items.AdditionalItems").Relation("Attendant").Relation("Payments").Relation("Groups.ComplementItem").Relation("Table").Relation("Delivery.Address").Relation("Pickup").Scan(ctx); err != nil {
		return nil, err
	}

	order.CalculateTotalPrice()
	return order, nil
}

func (r *OrderRepositoryBun) GetAllOrders(ctx context.Context) ([]orderentity.Order, error) {
	orders := []orderentity.Order{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&orders).Relation("Groups.Items.AdditionalItems").Relation("Attendant").Relation("Payments").Relation("Groups.ComplementItem").Relation("Table").Relation("Delivery.Address").Relation("Pickup").Scan(ctx); err != nil {
		return nil, err
	}

	for i := range orders {
		orders[i].CalculateTotalPrice()
	}

	return orders, nil
}

func (r *OrderRepositoryBun) AddPaymentOrder(ctx context.Context, payment *orderentity.PaymentOrder) error {
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
