package userdto

import companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"

type TokenAndSchemasOutput struct {
	AccessToken string                           `json:"accessToken"`
	Companies   []companyentity.CompanyWithUsers `json:"companies"`
}
