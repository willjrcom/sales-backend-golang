package productrepositorylocal

import productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"

type ProductRepositoryLocal struct {
}

func NewProductRepositoryLocal() *ProductRepositoryLocal {
	return &ProductRepositoryLocal{}
}

func (r *ProductRepositoryLocal) RegisterProduct(p *productentity.Product) error {
	return nil
}

func (r *ProductRepositoryLocal) UpdateProduct(p *productentity.Product) error {
	return nil
}

func (r *ProductRepositoryLocal) DeleteProduct(id string) error {
	return nil
}

func (r *ProductRepositoryLocal) GetProductById(id string) (*productentity.Product, error) {
	return nil, nil
}

func (r *ProductRepositoryLocal) GetProductBy(key string, value string) (*productentity.Product, error) {
	return nil, nil
}

func (r *ProductRepositoryLocal) GetAllProduct(key string, value string) ([]productentity.Product, error) {
	return nil, nil
}
