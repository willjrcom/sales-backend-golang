package userdto

import userentity "github.com/willjrcom/sales-backend-go/internal/domain/user"

type TokenAndSchemasOutput struct {
	AccessToken string                     `json:"accessToken"`
	Schemas     []userentity.SchemaCompany `json:"schemas"`
}
