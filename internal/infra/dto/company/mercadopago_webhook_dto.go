package companydto

type MercadoPagoWebhookDTO struct {
	ID      int64                     `json:"id"`
	Live    bool                      `json:"live_mode"`
	Type    string                    `json:"type"`
	Action  string                    `json:"action"`
	Data    MercadoPagoWebhookDataDTO `json:"data"`
	DateUTC string                    `json:"date_created"`
}

type MercadoPagoWebhookDataDTO struct {
	ID string `json:"id"`
}
