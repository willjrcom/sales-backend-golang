package companydto

import companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"

type CompanyOutput struct {
	companyentity.CompanyCommonAttributes
}

func (o *CompanyOutput) FromModel(model *companyentity.Company) {
	o.CompanyCommonAttributes = model.CompanyCommonAttributes
}
