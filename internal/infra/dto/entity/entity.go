package entitydto

import "github.com/google/uuid"

type IdRequest struct {
	Id uuid.UUID `json:"id"`
}
