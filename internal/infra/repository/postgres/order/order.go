package orderrepositorybun

import (
	"context"
	"database/sql"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
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
	order.CalculateTotalChange()

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

	if _, err := r.db.NewDelete().Model(&orderentity.Order{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *OrderRepositoryBun) GetOrderById(ctx context.Context, id string) (order *orderentity.Order, err error) {
	order = &orderentity.Order{}
	order.ID = uuid.MustParse(id)
	relation := ""

	if order.Table != nil {
		relation = "Table"
	} else if order.Delivery != nil {
		relation = "Delivery"
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	query := r.db.NewSelect().Model(order).WherePK().Relation("Groups.Items.AdditionalItems").Relation("Attendant").Relation("Payments")

	if relation != "" {
		query = query.Relation(relation)
	}

	if err = query.Scan(ctx); err != nil {
		return nil, err
	}

	order.CalculateTotalChange()
	return order, nil
}

func (r *OrderRepositoryBun) GetAllOrders(ctx context.Context) ([]orderentity.Order, error) {
	orders := []orderentity.Order{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&orders).Relation("Groups.Items.AdditionalItems").Relation("Groups.ComplementItem").Scan(ctx); err != nil {
		return nil, err
	}

	for i := range orders {
		orders[i].CalculateTotalChange()
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
