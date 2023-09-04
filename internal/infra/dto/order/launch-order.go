package orderdto

import "github.com/google/uuid"

type LaunchOrderInput struct {
	ID uuid.UUID `json:"id"`
}
