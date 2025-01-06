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
}

func (o *TableCreateDTO) validate() error {
	if o.Name == "" {
		return ErrNameRequired
	}

	return nil
}

func (o *TableCreateDTO) ToModel() (*tableentity.Table, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	tableCommonAttributes := tableentity.TableCommonAttributes{
		Name:        o.Name,
		IsAvailable: true,
	}

	table := &tableentity.Table{
		Entity:                entity.NewEntity(),
		TableCommonAttributes: tableCommonAttributes,
	}

	return table, nil
}
