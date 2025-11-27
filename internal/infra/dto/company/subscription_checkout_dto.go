package companydto

type SubscriptionCheckoutDTO struct {
	Months int `json:"months"`
}

func (s *SubscriptionCheckoutDTO) Normalize() int {
	if s == nil || s.Months <= 0 {
		return 1
	}

	if s.Months > 12 {
		return 12
	}

	return s.Months
}

type SubscriptionCheckoutResponseDTO struct {
	PreferenceID     string `json:"preference_id"`
	InitPoint        string `json:"init_point"`
	SandboxInitPoint string `json:"sandbox_init_point,omitempty"`
}
