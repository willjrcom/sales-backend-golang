package productentity

import (
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type ProductCategoryWithOrderProcess struct {
	entity.Entity
	ProductCategoryWithOrderProcessCommonAttributes
}

type ProductCategoryWithOrderProcessCommonAttributes struct {
	Name           string
	ImagePath      string
	UseProcessRule bool
	IsAdditional   bool
	IsComplement   bool
	NeedPrint      bool
	ProcessRules   []ProcessRuleWithOrderProcess
}
