package productEntity

import (
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Product struct {
	entity.Entity
	Code        string
	Name        string
	Description string
	Size        string
	Price       float64
	Cost        float64
	Category    CategoryProduct
	IsAvailable bool
}
