package orderentity

import (
	"context"

	"github.com/google/uuid"
)

type ItemRepository interface {
	AddItem(ctx context.Context, item *Item) error
	AddAdditionalItem(ctx context.Context, id uuid.UUID, productID uuid.UUID, additionalItem *Item) error
	DeleteItem(ctx context.Context, id string) error
	DeleteAdditionalItem(ctx context.Context, idAdditional uuid.UUID) error
	UpdateItem(ctx context.Context, item *Item) error
	GetItemById(ctx context.Context, id string) (*Item, error)
}

type GroupItemRepository interface {
	CreateGroupItem(ctx context.Context, groupitem *GroupItem) (err error)
	UpdateGroupItem(ctx context.Context, groupitem *GroupItem) (err error)
	GetGroupByID(ctx context.Context, id string, withRelation bool) (*GroupItem, error)
	DeleteGroupItem(ctx context.Context, id string, complementItemID *string) error
	GetGroupsByOrderIDAndStatus(ctx context.Context, id string, status StatusGroupItem) ([]GroupItem, error)
	GetGroupsByStatus(ctx context.Context, status StatusGroupItem) ([]GroupItem, error)
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	PendingOrder(ctx context.Context, order *Order) error
	UpdateOrder(ctx context.Context, order *Order) error
	DeleteOrder(ctx context.Context, id string) error
	GetOrderById(ctx context.Context, id string) (*Order, error)
	GetAllOrders(ctx context.Context) ([]Order, error)
	GetAllOrdersWithDelivery(ctx context.Context) ([]Order, error)
	AddPaymentOrder(ctx context.Context, payment *PaymentOrder) error
}

type OrderPickupRepository interface {
	CreateOrderPickup(ctx context.Context, pickup *OrderPickup) error
	UpdateOrderPickup(ctx context.Context, pickup *OrderPickup) error
	DeleteOrderPickup(ctx context.Context, id string) error
	GetPickupById(ctx context.Context, id string) (*OrderPickup, error)
	GetAllPickups(ctx context.Context) ([]OrderPickup, error)
}

type OrderDeliveryRepository interface {
	CreateOrderDelivery(ctx context.Context, delivery *OrderDelivery) error
	UpdateOrderDelivery(ctx context.Context, delivery *OrderDelivery) error
	DeleteOrderDelivery(ctx context.Context, id string) error
	GetDeliveryById(ctx context.Context, id string) (*OrderDelivery, error)
	GetDeliveriesByIds(ctx context.Context, ids []string) ([]OrderDelivery, error)
	GetAllDeliveries(ctx context.Context) ([]OrderDelivery, error)
}

type DeliveryDriverRepository interface {
	CreateDeliveryDriver(ctx context.Context, DeliveryDriver *DeliveryDriver) error
	UpdateDeliveryDriver(ctx context.Context, DeliveryDriver *DeliveryDriver) error
	DeleteDeliveryDriver(ctx context.Context, id string) error
	GetDeliveryDriverById(ctx context.Context, id string) (*DeliveryDriver, error)
	GetDeliveryDriverByEmployeeId(ctx context.Context, id string) (*DeliveryDriver, error)
	GetAllDeliveryDrivers(ctx context.Context) ([]DeliveryDriver, error)
}

type OrderTableRepository interface {
	CreateOrderTable(ctx context.Context, table *OrderTable) error
	UpdateOrderTable(ctx context.Context, table *OrderTable) error
	DeleteOrderTable(ctx context.Context, id string) error
	GetOrderTableById(ctx context.Context, id string) (*OrderTable, error)
	GetPendingOrderTablesByTableId(ctx context.Context, id string) ([]OrderTable, error)
	GetAllOrderTables(ctx context.Context) ([]OrderTable, error)
}
