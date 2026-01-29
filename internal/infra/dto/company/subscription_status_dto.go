package companydto

type SubscriptionStatusDTO struct {
	CurrentPlan     string  `json:"current_plan"`      // "free", "basic", "intermediate", "advanced"
	ExpiresAt       *string `json:"expires_at"`        // ISO 8601 ou null se free
	DaysRemaining   *int    `json:"days_remaining"`    // null se free
	UpcomingPlan    *string `json:"upcoming_plan"`     // null se não houver
	UpcomingStartAt *string `json:"upcoming_start_at"` // null se não houver
}
