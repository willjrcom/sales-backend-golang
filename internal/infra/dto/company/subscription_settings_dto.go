package companydto

// SubscriptionSettingsDTO exposes configuration used to build Mercado Pago checkouts.
type SubscriptionSettingsDTO struct {
	MonthlyPrice  float64 `json:"monthly_price"`
	Currency      string  `json:"currency"`
	MinMonths     int     `json:"min_months"`
	MaxMonths     int     `json:"max_months"`
	DefaultMonths int     `json:"default_months"`
}
