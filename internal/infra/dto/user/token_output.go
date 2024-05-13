package userdto

import companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"

type TokenAndSchemasOutput struct {
	AccessToken string                           `json:"access_token"`
	Companies   []companyentity.CompanyWithUsers `json:"companies"`
}
