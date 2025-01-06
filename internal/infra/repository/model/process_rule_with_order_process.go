package model

type ProcessRuleWithOrderProcess struct {
	ProcessRule
	OrderProcesses []OrderProcess `bun:"rel:has-one,join:id=process_rule_id"`
}
