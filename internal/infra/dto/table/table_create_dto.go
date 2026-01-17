package tabledto

import (
	"errors"

	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
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

func (o *TableCreateDTO) ToDomain() (*tableentity.Table, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	isActive := true
	if o.IsActive != nil {
		isActive = *o.IsActive
	}

	tableCommonAttributes := tableentity.TableCommonAttributes{
		Name:        o.Name,
		IsAvailable: true,
		IsActive:    isActive,
	}

	table := &tableentity.Table{
		Entity:                entity.NewEntity(),
		TableCommonAttributes: tableCommonAttributes,
	}

	return table, nil
}
