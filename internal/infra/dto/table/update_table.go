package tabledto

import (
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
)

type TableUpdateDTO struct {
	Name        *string `json:"name"`
	IsAvailable *bool   `json:"is_available"`
}

func (c *TableUpdateDTO) UpdateModel(place *tableentity.Table) (err error) {
	if c.Name != nil {
		place.Name = *c.Name
	}
	if c.IsAvailable != nil {
		place.IsAvailable = *c.IsAvailable
	}

	return nil
}
