package companyentity

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	companycategoryentity "github.com/willjrcom/sales-backend-go/internal/domain/company_category"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/cnpj"
)

var (
	ErrInvalidMonth = errors.New("invalid month: must be between 1 and 12")
	ErrInvalidYear  = errors.New("invalid year: must be 2020 or later")
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
	Preferences  Preferences
	IsBlocked    bool

	// Categories
	Categories []companycategoryentity.CompanyCategory

	// Billing
	MonthlyPaymentDueDay          int
	MonthlyPaymentDueDayUpdatedAt *time.Time
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
			Preferences:  NewDefaultPreferences(),
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
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	id, _ := shortid.Generate()
	replacedName := reg.ReplaceAllString(strings.ToLower(cnpjData.TradeName), "_")
	safeID := reg.ReplaceAllString(id, "_")
	schema := "company_" + replacedName + "_" + strings.ToLower(safeID)
	return schema
}

func (c *Company) AddAddress(addressCommonAttributes *addressentity.AddressCommonAttributes) {
	c.Address = addressentity.NewAddress(addressCommonAttributes)
	c.Address.ObjectID = c.ID
}
