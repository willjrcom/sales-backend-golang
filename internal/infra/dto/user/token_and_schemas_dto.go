package userdto

import (
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
)

type TokenAndSchemasOutput struct {
	Person      personentity.Person              `json:"person"`
	AccessToken string                           `json:"access_token"`
	Companies   []companyentity.CompanyWithUsers `json:"companies"`
}
