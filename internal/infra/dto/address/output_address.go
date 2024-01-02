package addressdto

import (
	"github.com/google/uuid"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
)

type AddressOutput struct {
	ID uuid.UUID `json:"id"`
	addressentity.AddressCommonAttributes
}

func (a *AddressOutput) FromModel(model *addressentity.Address) {
	a.ID = model.ID
	a.AddressCommonAttributes = model.AddressCommonAttributes
}
