package ordertabledto

import (
	"github.com/google/uuid"
)

type OrderTableContactInput struct {
	Contact string    `json:"contact,omitempty"`
	TableID uuid.UUID `json:"table_id"`
}

func (o *OrderTableContactInput) validate() error {
	if o.TableID == uuid.Nil {
		return ErrTableIDRequired
	}

	return nil
}

func (o *OrderTableContactInput) ToDomain() (uuid.UUID, string, error) {
	if err := o.validate(); err != nil {
		return uuid.Nil, "", err
	}

	return o.TableID, o.Contact, nil
}
