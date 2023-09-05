package orderentity

type Repository interface {
	CreateOrder(order *Order) error
	UpdateOrder(order *Order) error
	DeleteOrder(id string) error
	GetOrderById(id string) (*Order, error)
	GetOrderBy(key string, value string) (*Order, error)
	GetAllOrder(key string, value string) ([]Order, error)
}
