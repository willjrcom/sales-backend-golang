package orderentity

import itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"

type Repository interface {
	CreateOrder(order *Order) error
	AddItemOrder(item *itementity.Item, items itementity.Items) error
	UpdateOrder(order *Order) error
	GetOrder(id string) (*Order, error)
	GetAllOrder() ([]Order, error)
}
