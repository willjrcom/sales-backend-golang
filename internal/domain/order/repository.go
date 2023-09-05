package orderentity

import itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"

type Repository interface {
	CreateOrder(order *Order) error
	AddItemOrder(item *itementity.Item, items itementity.Items) error
	UpdateOrder(order *Order) error
	GetOrderById(id string) (*Order, error)
	GetOrderBy(key string, value string) (*Order, error)
	GetAllOrder(key string, value string) ([]Order, error)
}
