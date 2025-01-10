package productentity

import (
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type ProcessRuleWithOrderProcess struct {
	entity.Entity
	ProcessRuleCommonAttributes
	ProcessRuleWithOrderProcessCommonAttributes
}

type ProcessRuleWithOrderProcessCommonAttributes struct {
	TotalOrderProcessLate int
	TotalOrderQueue       int
}
