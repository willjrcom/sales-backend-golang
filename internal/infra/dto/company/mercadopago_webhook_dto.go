package companydto

import (
	"encoding/json"
	"fmt"
)

type MercadoPagoWebhookDTO struct {
	ID      FlexibleID                `json:"id"`
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

type FlexibleID string

func (fi *FlexibleID) UnmarshalJSON(data []byte) error {
	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		*fi = FlexibleID(fmt.Sprintf("%d", i))
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*fi = FlexibleID(s)
		return nil
	}
	return nil
}
