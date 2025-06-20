package employeedto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
)

type EmployeePaymentDTO struct {
	ID uuid.UUID `json:"id"`
	PaymentTimeLogs
	PaymentCommonAttributes
}

type PaymentCommonAttributes struct {
	EmployeeID uuid.UUID                    `json:"employee_id"`
	Amount     decimal.Decimal              `json:"amount"`
	Status     employeeentity.PaymentStatus `json:"status"`
	Method     employeeentity.PaymentMethod `json:"method"`
	Notes      string                       `json:"notes"`
}

type PaymentTimeLogs struct {
	PaymentDate time.Time `json:"payment_date"`
}

func (d *EmployeePaymentDTO) FromDomain(payment *employeeentity.PaymentEmployee) {
	if payment == nil {
		return
	}
	*d = EmployeePaymentDTO{
		ID: payment.ID,
		PaymentCommonAttributes: PaymentCommonAttributes{
			EmployeeID: payment.EmployeeID,
			Amount:     payment.Amount,
			Status:     payment.Status,
			Method:     payment.Method,
			Notes:      payment.Notes,
		},
		PaymentTimeLogs: PaymentTimeLogs{
			PaymentDate: payment.PaymentDate,
		},
	}
}
