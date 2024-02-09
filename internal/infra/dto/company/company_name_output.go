package companydto

import (
	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

type CompanyNameOutput struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Cnpj string    `json:"cnpj"`
}

func (o *CompanyNameOutput) FromModel(model *companyentity.Company) {
	o.ID = model.ID
	o.Name = model.TradeName
	o.Cnpj = model.Cnpj
}
