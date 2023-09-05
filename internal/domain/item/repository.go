package itementity

type Repository interface {
	AddItemOrder(item *Item) error
	RemoveItemOrder(item *Item) error
	UpdateItemOrder(id string, item *Item) error
}
