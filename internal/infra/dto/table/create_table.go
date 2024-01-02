package tabledto

import (
	"errors"

	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
)

var (
	ErrNameRequired = errors.New("name is required")
)

type CreateTableInput struct {
	tableentity.TableCommonAttributes
}

func (o *CreateTableInput) Validate() error {
	if o.Name == "" {
		return ErrNameRequired
	}

	return nil
}

func (o *CreateTableInput) ToModel() (*tableentity.Table, error) {
	if err := o.Validate(); err != nil {
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
