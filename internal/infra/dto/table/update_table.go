package tabledto

import (
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
)

type UpdateTableInput struct {
	tableentity.PatchTable
}

func (c *UpdateTableInput) UpdateModel(place *tableentity.Table) (err error) {
	if c.Name != nil {
		place.Name = *c.Name
	}
	if c.IsAvailable != nil {
		place.IsAvailable = *c.IsAvailable
	}

	return nil
}
