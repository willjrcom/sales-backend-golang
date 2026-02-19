package tabledto

import (
	"errors"

	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

var (
	ErrNameRequired = errors.New("name is required")
)

type TableCreateDTO struct {
	Name        string `json:"name"`
	IsAvailable bool   `json:"is_available"`
	IsActive    *bool  `json:"is_active"`
}

func (o *TableCreateDTO) validate() error {
	if o.Name == "" {
		return ErrNameRequired
	}

	return nil
}

func (o *TableCreateDTO) ToDomain() (*orderentity.Table, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	isActive := true
	if o.IsActive != nil {
		isActive = *o.IsActive
	}

	tableCommonAttributes := orderentity.TableCommonAttributes{
		Name:        o.Name,
		IsAvailable: true,
		IsActive:    isActive,
	}

	table := &orderentity.Table{
		Entity:                entity.NewEntity(),
		TableCommonAttributes: tableCommonAttributes,
	}

	return table, nil
}
