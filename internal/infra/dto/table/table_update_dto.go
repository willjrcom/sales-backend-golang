package tabledto

import orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"

type TableUpdateDTO struct {
	Name        *string `json:"name"`
	IsAvailable *bool   `json:"is_available"`
	IsActive    *bool   `json:"is_active"`
}

func (c *TableUpdateDTO) UpdateDomain(table *orderentity.Table) (err error) {
	if c.Name != nil {
		table.Name = *c.Name
	}
	if c.IsAvailable != nil {
		table.IsAvailable = *c.IsAvailable
	}
	if c.IsActive != nil {
		table.IsActive = *c.IsActive
	}

	return nil
}
