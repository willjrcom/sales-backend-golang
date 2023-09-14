package contactdto

import "github.com/google/uuid"

type ContactOutput struct {
	ID     uuid.UUID `json:"id"`
	Ddd    string    `json:"ddd"`
	Number string    `json:"number"`
}
