package companydto

type MercadoPagoWebhookDTO struct {
	ID      int64                     `json:"id"`
	Live    bool                      `json:"live_mode"`
	Type    string                    `json:"type"`
	Action  string                    `json:"action"`
	Data    MercadoPagoWebhookDataDTO `json:"data"`
	DateUTC string                    `json:"date_created"`

	// Headers and query params for signature validation (populated by handler)
	XSignature      string `json:"-"`
	XRequestID      string `json:"-"`
	DataIDFromQuery string `json:"-"` // data.id from query params for signature validation
}

type MercadoPagoWebhookDataDTO struct {
	ID string `json:"id"`
}
