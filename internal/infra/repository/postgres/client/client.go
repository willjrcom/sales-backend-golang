package clientrepositorylocal

import cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"

type ProductRepositoryLocal struct {
}

func NewProductRepositoryLocal() *ProductRepositoryLocal {
	return &ProductRepositoryLocal{}
}

func (r *ProductRepositoryLocal) RegisterClient(p *cliententity.Client) error {
	return nil
}

func (r *ProductRepositoryLocal) UpdateClient(p *cliententity.Client) error {
	return nil
}

func (r *ProductRepositoryLocal) DeleteClient(id string) error {
	return nil
}

func (r *ProductRepositoryLocal) GetClientById(id string) (*cliententity.Client, error) {
	return nil, nil
}

func (r *ProductRepositoryLocal) GetClientBy(key string, value string) (*cliententity.Client, error) {
	return nil, nil
}

func (r *ProductRepositoryLocal) GetAllClient(key string, value string) ([]cliententity.Client, error) {
	return nil, nil
}
