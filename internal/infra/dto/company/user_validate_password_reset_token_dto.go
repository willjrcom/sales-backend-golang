package companydto

type UserResetTokenRequestDTO struct {
	Token string `json:"token"`
}

type UserResetTokenResponseDTO struct {
	Valid bool   `json:"valid"`
	Email string `json:"email,omitempty"`
}
