package categorydto

import (
	"errors"

	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrNameIsEmpty = errors.New("name is empty")
)

type RegisterCategoryInput struct {
	Name string `json:"name"`
}

func (c *RegisterCategoryInput) validate() error {
	if c.Name == "" {
		return ErrNameIsEmpty
	}

	return nil
}

func (c *RegisterCategoryInput) ToModel() (*productentity.Category, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	return &productentity.Category{
		Entity: entity.NewEntity(),
		Name:   c.Name,
	}, nil
}
