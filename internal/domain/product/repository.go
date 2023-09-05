package productentity

type Repository interface {
	RegisterProduct(p *Product) error
	UpdateProduct(p *Product) error
	DeleteProduct(id string) error
	GetProductById(id string) (*Product, error)
	GetProductBy(key string, value string) (*Product, error)
	GetAllProduct(key string, value string) ([]Product, error)
}
