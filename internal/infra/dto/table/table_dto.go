package tabledto

import (
	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type TableDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        *string   `json:"name"`
	IsAvailable *bool     `json:"is_available"`
	IsActive    bool      `json:"is_active"`
}

func (c *TableDTO) FromDomain(table *orderentity.Table) (err error) {
	if table == nil {
		return
	}
	*c = TableDTO{
		ID:          table.ID,
		Name:        &table.Name,
		IsAvailable: &table.IsAvailable,
		IsActive:    table.IsActive,
	}

	return nil
}
