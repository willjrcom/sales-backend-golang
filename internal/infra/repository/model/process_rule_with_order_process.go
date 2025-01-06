package model

type ProcessRuleWithOrderProcess struct {
	ProcessRule
	OrderProcesses []OrderProcess `bun:"rel:has-one,join:process_rule_id=process_rule_id"`
}
