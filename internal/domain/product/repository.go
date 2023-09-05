package productentity

type Repository interface {
	RegisterProduct(p *Product) error
	UpdateProduct(p *Product) error
	DeleteProduct(id string) error
	GetProduct(id string) (*Product, error)
	GetAllProduct(category string) ([]Product, error)
}
