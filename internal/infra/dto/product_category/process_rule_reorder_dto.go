package productcategorydto

import "github.com/google/uuid"

type ProcessRuleReorderDTO struct {
	ProcessRules []ProcessRuleOrderDTO `json:"process_rules"`
}

type ProcessRuleOrderDTO struct {
	ID    uuid.UUID `json:"id"`
	Order int8      `json:"order"`
}
