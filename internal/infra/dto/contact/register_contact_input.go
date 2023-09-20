package contactdto

import (
	"errors"

	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

var (
	ErrDddIsEmpty          = errors.New("ddd is empty")
	ErrNumberLengthInvalid = errors.New("number length invalid")
)

type RegisterContactInput struct {
	Ddd    string `json:"ddd"`
	Number string `json:"number"`
}

func (c *RegisterContactInput) validate() error {
	if c.Ddd == "" {
		return ErrDddIsEmpty
	}
	if len(c.Number) < 8 || len(c.Number) > 10 {
		return ErrNumberLengthInvalid
	}

	return nil
}

func (c *RegisterContactInput) ToModel() (*personentity.Contact, error) {
	if err := c.validate(); err != nil {
		return nil, err
	}

	return &personentity.Contact{
		Entity: entity.NewEntity(),
		Ddd:    c.Ddd,
		Number: c.Number,
	}, nil
}
