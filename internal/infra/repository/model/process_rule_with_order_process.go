package model

import (
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
)

type ProcessRuleWithOrderProcess struct {
	ProcessRule
	OrderProcesses []orderprocessentity.OrderProcess `bun:"rel:has-one,join:process_rule_id=process_rule_id"`
}
