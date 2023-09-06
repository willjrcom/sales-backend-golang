package productentity

type Repository interface {
	RegisterProduct(p *Product) error
	UpdateProduct(p *Product) error
	DeleteProduct(id string) error
	GetProductById(id string) (*Product, error)
	GetProductBy(key string, value string) (*Product, error)
	GetAllProduct(key string, value string) ([]Product, error)
}

type RepositoryCategory interface {
	RegisterCategoryProduct(category *CategoryProduct) error
	UpdateCategoryProduct(category *CategoryProduct) error
	DeleteCategoryProduct(id string) error
	GetCategoryProductById(id string) (*CategoryProduct, error)
	GetAllCategoryProduct() ([]CategoryProduct, error)
}
