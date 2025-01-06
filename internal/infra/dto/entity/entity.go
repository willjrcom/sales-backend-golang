package entitydto

import "github.com/google/uuid"

type IDRequest struct {
	ID uuid.UUID `json:"id"`
}

func NewIdRequest(id uuid.UUID) *IDRequest {
	return &IDRequest{
		ID: id,
	}
}
