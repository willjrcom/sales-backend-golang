package tabledto

import (
	"github.com/google/uuid"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
)

type TableDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        *string   `json:"name"`
	IsAvailable *bool     `json:"is_available"`
}

func (c *TableDTO) FromDomain(table *tableentity.Table) (err error) {
	*c = TableDTO{
		ID:          table.ID,
		Name:        &table.Name,
		IsAvailable: &table.IsAvailable,
	}

	return nil
}
