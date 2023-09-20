package addressdto

import addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"

type AddressOutput struct {
	ID           string `json:"id"`
	Street       string `json:"street"`
	Number       string `json:"number"`
	Complement   string `json:"complement"`
	Reference    string `json:"reference"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	Cep          string `json:"cep"`
}

func (a *AddressOutput) FromModel(model *addressentity.Address) {
	a.ID = model.ID.String()
	a.Street = model.Street
	a.Number = model.Number
	a.Complement = model.Complement
	a.Reference = model.Reference
	a.Neighborhood = model.Neighborhood
	a.City = model.City
	a.State = model.State
	a.Cep = model.Cep
}
