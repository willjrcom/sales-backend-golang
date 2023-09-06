package entitydto

import "github.com/google/uuid"

type IdRequest struct {
	ID uuid.UUID `json:"id"`
}

func NewIdRequest(id uuid.UUID) *IdRequest {
	return &IdRequest{
		ID: id,
	}
}
