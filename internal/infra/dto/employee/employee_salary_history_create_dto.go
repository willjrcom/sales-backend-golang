package employeedto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type EmployeeSalaryHistoryCreateDTO struct {
	EmployeeID uuid.UUID       `json:"employee_id"`
	StartDate  time.Time       `json:"start_date"`
	EndDate    *time.Time      `json:"end_date,omitempty"`
	SalaryType string          `json:"salary_type"`
	BaseSalary decimal.Decimal `json:"base_salary"`
	HourlyRate decimal.Decimal `json:"hourly_rate"`
	Commission decimal.Decimal `json:"commission"`
}

func (dto *EmployeeSalaryHistoryCreateDTO) ToDomain() *employeeentity.EmployeeSalaryHistory {
	return &employeeentity.EmployeeSalaryHistory{
		Entity:     entity.NewEntity(),
		EmployeeID: dto.EmployeeID,
		StartDate:  dto.StartDate,
		EndDate:    dto.EndDate,
		SalaryType: dto.SalaryType,
		BaseSalary: dto.BaseSalary,
		HourlyRate: dto.HourlyRate,
		Commission: dto.Commission,
	}
}

type EmployeeSalaryHistoryDTO struct {
	ID         uuid.UUID       `json:"id"`
	EmployeeID uuid.UUID       `json:"employee_id"`
	StartDate  time.Time       `json:"start_date"`
	EndDate    *time.Time      `json:"end_date,omitempty"`
	SalaryType string          `json:"salary_type"`
	BaseSalary decimal.Decimal `json:"base_salary"`
	HourlyRate decimal.Decimal `json:"hourly_rate"`
	Commission decimal.Decimal `json:"commission"`
}

func (d *EmployeeSalaryHistoryDTO) FromDomain(h *employeeentity.EmployeeSalaryHistory) {
	if h == nil {
		return
	}
	*d = EmployeeSalaryHistoryDTO{
		ID:         h.ID,
		EmployeeID: h.EmployeeID,
		StartDate:  h.StartDate,
		EndDate:    h.EndDate,
		SalaryType: h.SalaryType,
		BaseSalary: h.BaseSalary,
		HourlyRate: h.HourlyRate,
		Commission: h.Commission,
	}
}
