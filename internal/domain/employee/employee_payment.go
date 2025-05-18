package employeeentity

import (
	"time"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type PaymentEmployee struct {
	entity.Entity
	PaymentCommonAttributes
	PaymentTimeLogs
}

type PaymentCommonAttributes struct {
	EmployeeID uuid.UUID
	Amount     float64
	Status     PaymentStatus
	Method     PaymentMethod
	Notes      string
}

type PaymentTimeLogs struct {
	PayDate time.Time
}

type PaymentStatus string

const (
	StatusPending   PaymentStatus = "Pending"
	StatusCompleted PaymentStatus = "Completed"
	StatusCancelled PaymentStatus = "Cancelled"
)

func GetAllPaymentStatus() []PaymentStatus {
	return []PaymentStatus{
		StatusPending,
		StatusCompleted,
		StatusCancelled,
	}
}

type PaymentMethod string

const (
	MethodCash         PaymentMethod = "Cash"
	MethodBankTransfer PaymentMethod = "BankTransfer"
	MethodCheck        PaymentMethod = "Check"
	MethodOther        PaymentMethod = "Other"
)

func GetAllPaymentMethods() []PaymentMethod {
	return []PaymentMethod{
		MethodCash,
		MethodBankTransfer,
		MethodCheck,
		MethodOther,
	}
}

func NewPaymentEmployee(
	employeeID uuid.UUID,
	amount float64,
	status PaymentStatus,
	method PaymentMethod,
	payDate time.Time,
	notes string,
) *PaymentEmployee {
	return &PaymentEmployee{
		Entity: entity.NewEntity(),
		PaymentCommonAttributes: PaymentCommonAttributes{
			EmployeeID: employeeID,
			Amount:     amount,
			Status:     status,
			Method:     method,
			Notes:      notes,
		},
		PaymentTimeLogs: PaymentTimeLogs{
			PayDate: payDate,
		},
	}
}
