package itementity

import (
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Item struct {
	entity.Entity
	Name        string
	Quantity    float64
	Description string
	Price       float64
	Status      StatusItem
	ItemsID     uuid.UUID
	Observation string
}
