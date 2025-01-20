package companyentity

import (
	"strings"

	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/cnpj"
)

type Company struct {
	entity.Entity
	CompanyCommonAttributes
}

type CompanyCommonAttributes struct {
	SchemaName   string
	BusinessName string
	TradeName    string
	Cnpj         string
	Email        string
	Contacts     []string
	Address      *addressentity.Address
	Users        []User
}

type CompanyToUsers struct {
	CompanyID uuid.UUID
	Company   *Company
	UserID    uuid.UUID
	User      *User
}

func NewCompany(cnpjData *cnpj.Cnpj) *Company {
	schema := generateSchema(cnpjData)

	company := &Company{
		Entity: entity.NewEntity(),
		CompanyCommonAttributes: CompanyCommonAttributes{
			BusinessName: cnpjData.BusinessName,
			TradeName:    cnpjData.TradeName,
			Cnpj:         cnpjData.Cnpj,
			SchemaName:   schema,
		},
	}

	addressCommonAttributes := &addressentity.AddressCommonAttributes{
		Street:       cnpjData.Street,
		Number:       cnpjData.Number,
		Neighborhood: cnpjData.Neighborhood,
		City:         cnpjData.City,
		UF:           cnpjData.UF,
		Cep:          cnpjData.Cep,
		AddressType:  addressentity.AddressTypeWork,
	}

	company.AddAddress(addressCommonAttributes)
	return company

}

func (c *Company) UpdateCompany(cnpjData *cnpj.Cnpj) {
	c.BusinessName = cnpjData.BusinessName
	c.TradeName = cnpjData.TradeName
	c.Cnpj = cnpjData.Cnpj

	addressCommonAttributes := &addressentity.AddressCommonAttributes{
		Street:       cnpjData.Street,
		Number:       cnpjData.Number,
		Neighborhood: cnpjData.Neighborhood,
		City:         cnpjData.City,
		UF:           cnpjData.UF,
		Cep:          cnpjData.Cep,
		AddressType:  addressentity.AddressTypeWork,
	}

	c.AddAddress(addressCommonAttributes)
}

func generateSchema(cnpjData *cnpj.Cnpj) string {
	id, _ := shortid.Generate()
	replacedName := strings.ReplaceAll(strings.ToLower(cnpjData.TradeName), " ", "_")
	replacedName = strings.ReplaceAll(replacedName, "-", "_")
	id = strings.ReplaceAll(id, "-", "_")
	schema := "loja_" + replacedName + "_" + strings.ToLower(id)
	return schema
}

func (c *Company) AddAddress(addressCommonAttributes *addressentity.AddressCommonAttributes) {
	c.Address = addressentity.NewAddress(addressCommonAttributes)
	c.Address.ObjectID = c.ID
}
