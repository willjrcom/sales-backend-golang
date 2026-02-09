package employeeentity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Employee struct {
	entity.Entity
	UserID      uuid.UUID
	User        *companyentity.User
	Payments    []PaymentEmployee
	Permissions Permissions
	IsActive    bool
}

func NewEmployee(userID uuid.UUID) *Employee {
	return &Employee{
		Entity:      entity.NewEntity(),
		UserID:      userID,
		Payments:    make([]PaymentEmployee, 0),
		Permissions: make(Permissions),
		IsActive:    true,
	}
}

type EmployeeSalaryHistory struct {
	entity.Entity
	EmployeeID uuid.UUID
	StartDate  time.Time
	EndDate    *time.Time
	SalaryType string
	BaseSalary decimal.Decimal
	HourlyRate decimal.Decimal
	Commission float64
}

// PermissionKey define as permissões possíveis para um funcionário.
type PermissionKey string

// Enum de permissões
const (
	PermissionBilling              PermissionKey = "billing"
	PermissionClient               PermissionKey = "client"
	PermissionEmployee             PermissionKey = "employee"
	PermissionNewOrder             PermissionKey = "new-order"
	PermissionOrderControl         PermissionKey = "order-control"
	PermissionOrderDeliveryControl PermissionKey = "order-delivery-control"
	PermissionOrderPickupControl   PermissionKey = "order-pickup-control"
	PermissionOrderProcess         PermissionKey = "order-process"
	PermissionOrderTableControl    PermissionKey = "order-table-control"
	PermissionPlace                PermissionKey = "place"
	PermissionPrint                PermissionKey = "print"
	PermissionProduct              PermissionKey = "product"
	PermissionManageStock          PermissionKey = "manage-stock"
	PermissionShift                PermissionKey = "shift"
	PermissionManageCompany        PermissionKey = "manage-company"
)

// GetAllPermissions retorna todas as permissões possíveis
func GetAllPermissions() []PermissionKey {
	return []PermissionKey{
		PermissionBilling,
		PermissionClient,
		PermissionEmployee,
		PermissionNewOrder,
		PermissionOrderControl,
		PermissionOrderDeliveryControl,
		PermissionOrderPickupControl,
		PermissionOrderProcess,
		PermissionOrderTableControl,
		PermissionPlace,
		PermissionPrint,
		PermissionProduct,
		PermissionManageStock,
		PermissionShift,
		PermissionManageCompany,
	}
}

// Permission representa um par chave-valor de permissão.
type Permission struct {
	Key   PermissionKey
	Value bool
}

// Permissions é um map de permissões.
type Permissions map[PermissionKey]bool

// NewPermissions cria um map de permissões a partir de um slice de Permission.
func NewPermissions(entries []Permission) Permissions {
	perms := make(Permissions, len(entries))
	for _, e := range entries {
		perms[e.Key] = e.Value
	}
	return perms
}

// Métodos utilitários para acessar permissões
func (p Permissions) GetBool(key PermissionKey) (bool, error) {
	v, ok := p[key]
	if !ok {
		return false, fmt.Errorf("permission %q not found", key)
	}
	return v, nil
}
