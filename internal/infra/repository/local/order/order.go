package orderrepositorylocal

import (
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type ProductRepositoryLocal struct {
}

func NewProductRepositoryLocal() *ProductRepositoryLocal {
	return &ProductRepositoryLocal{}
}

func (r *ProductRepositoryLocal) CreateOrder(order *orderentity.Order) error {
	return nil
}

func (r *ProductRepositoryLocal) AddItemOrder(item *itementity.Item, items itementity.Items) error {
	return nil
}

func (r *ProductRepositoryLocal) UpdateOrder(order *orderentity.Order) error {
	return nil
}

func (r *ProductRepositoryLocal) GetOrderById(id string) (*orderentity.Order, error) {
	return nil, nil
}

func (r *ProductRepositoryLocal) GetOrderBy(key string, value string) (*orderentity.Order, error) {
	return nil, nil
}

func (r *ProductRepositoryLocal) GetAllOrder(key string, value string) ([]orderentity.Order, error) {
	return nil, nil
}
