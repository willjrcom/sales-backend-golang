package companydto

type SubscriptionStatusDTO struct {
	CurrentPlan    string    `json:"current_plan"`   // "free", "basic", "intermediate", "advanced"
	ExpiresAt      *string   `json:"expires_at"`     // ISO 8601 ou null se free
	DaysRemaining  *int      `json:"days_remaining"` // null se free
	IsCancelled    bool      `json:"is_cancelled"`   // true se tem assinatura ativa cancelável
	Frequency      string    `json:"frequency"`      // "MONTHLY", "SEMIANNUALLY", "ANNUALLY" (derived from active subscription)
	AvailablePlans []PlanDTO `json:"available_plans,omitempty"`
}

type PlanDTO struct {
	Key          string   `json:"key"`
	Name         string   `json:"name"`
	Price        float64  `json:"price"`
	Features     []string `json:"features"`
	IsCurrent    bool     `json:"is_current"`
	IsUpgrade    bool     `json:"is_upgrade"`
	UpgradePrice *float64 `json:"upgrade_price,omitempty"`
	Order        int      `json:"order"`
}
