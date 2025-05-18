package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type PaymentEmployee struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:employee_payments,alias:payment"`
	EmployeePaymentCommonAttributes
	EmployeePaymentTimeLogs
}

type EmployeePaymentCommonAttributes struct {
	EmployeeID uuid.UUID                    `bun:"employee_id,type:uuid,notnull"`
	Amount     decimal.Decimal              `bun:"amount,notnull"`
	Status     employeeentity.PaymentStatus `bun:"status,notnull"`
	Method     employeeentity.PaymentMethod `bun:"method,notnull"`
	Notes      string                       `bun:"notes"`
}

type EmployeePaymentTimeLogs struct {
	PayDate time.Time `bun:"pay_date,notnull"`
}

func (p *PaymentEmployee) FromDomain(payment *employeeentity.PaymentEmployee) {
	if payment == nil {
		return
	}
	*p = PaymentEmployee{
		Entity: entitymodel.FromDomain(payment.Entity),
		EmployeePaymentCommonAttributes: EmployeePaymentCommonAttributes{
			EmployeeID: payment.EmployeeID,
			Amount:     payment.Amount,
			Status:     payment.Status,
			Method:     payment.Method,
			Notes:      payment.Notes,
		},
		EmployeePaymentTimeLogs: EmployeePaymentTimeLogs{
			PayDate: payment.PayDate,
		},
	}
}

func (p *PaymentEmployee) ToDomain() *employeeentity.PaymentEmployee {
	if p == nil {
		return nil
	}
	return &employeeentity.PaymentEmployee{
		Entity: p.Entity.ToDomain(),
		PaymentCommonAttributes: employeeentity.PaymentCommonAttributes{
			EmployeeID: p.EmployeeID,
			Amount:     p.Amount,
			Status:     p.Status,
			Method:     p.Method,
			Notes:      p.Notes,
		},
		PaymentTimeLogs: employeeentity.PaymentTimeLogs{
			PayDate: p.PayDate,
		},
	}
}
