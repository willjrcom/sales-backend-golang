package orderentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type OrderTable struct {
	entity.Entity
	OrderTableCommonAttributes
	OrderTableTimeLogs
}

type OrderTableCommonAttributes struct {
	Name        string
	Contact     string
	Status      StatusOrderTable
	TaxRate     decimal.Decimal
	OrderID     uuid.UUID
	TableID     uuid.UUID
	Table       *Table
	OrderNumber int
}

type OrderTableTimeLogs struct {
	PendingAt   *time.Time
	ClosedAt    *time.Time
	CancelledAt *time.Time
}

func NewTable(orderTableCommonAttributes OrderTableCommonAttributes) *OrderTable {
	orderTableCommonAttributes.Status = OrderTableStatusStaging

	return &OrderTable{
		Entity:                     entity.NewEntity(),
		OrderTableCommonAttributes: orderTableCommonAttributes,
	}
}

func (t *OrderTable) Pend() error {
	if t.Status != OrderTableStatusStaging {
		return nil
	}

	t.Status = OrderTableStatusPending
	t.PendingAt = &time.Time{}
	*t.PendingAt = time.Now().UTC()
	return nil
}

func (t *OrderTable) Close() error {
	t.Status = OrderTableStatusClosed
	t.ClosedAt = &time.Time{}
	*t.ClosedAt = time.Now().UTC()
	return nil
}

func (t *OrderTable) Cancel() error {
	t.Status = OrderTableStatusCancelled
	t.CancelledAt = &time.Time{}
	*t.CancelledAt = time.Now().UTC()
	return nil
}

func (t *OrderTable) UpdatePreferences(preferences companyentity.Preferences) {
	t.TaxRate, _ = preferences.GetDecimal(companyentity.TableTaxRate)
}
