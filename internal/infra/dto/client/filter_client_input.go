package clientdto

import (
	"errors"

	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
)

var (
	ErrFilterIsEmpty = errors.New("filter is empty")
)

type FilterClientInput struct {
	Name *string `json:"name"`
	Cpf  *string `json:"cpf"`
}

func (f *FilterClientInput) validate() error {
	if f.Name == nil && f.Cpf == nil {
		return ErrFilterIsEmpty
	}

	return nil
}

func (f *FilterClientInput) ToModel() (*cliententity.Client, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}

	client := &cliententity.Client{}

	if f.Name != nil {
		client.Name = *f.Name
	}
	if f.Cpf != nil {
		client.Cpf = *f.Cpf
	}

	return client, nil
}
