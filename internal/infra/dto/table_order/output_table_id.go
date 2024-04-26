package tableorderdto

import "github.com/google/uuid"

type TableIDAndOrderIDOutput struct {
	TableID uuid.UUID `json:"table_id"`
	OrderID uuid.UUID `json:"order_id"`
}

func NewOutput(tableID uuid.UUID, orderID uuid.UUID) *TableIDAndOrderIDOutput {
	return &TableIDAndOrderIDOutput{
		TableID: tableID,
		OrderID: orderID,
	}
}
