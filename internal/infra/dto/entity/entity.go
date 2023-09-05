package entitydto

import "github.com/google/uuid"

type IdRequest struct {
	Id uuid.UUID `json:"id"`
}

func NewIdRequest(id uuid.UUID) *IdRequest {
	return &IdRequest{
		Id: id,
	}
}
