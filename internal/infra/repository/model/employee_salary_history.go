package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type EmployeeSalaryHistory struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:employee_salary_histories"`
	EmployeeSalaryHistoryCommonAttributes
}

type EmployeeSalaryHistoryCommonAttributes struct {
	EmployeeID uuid.UUID       `bun:"employee_id,type:uuid,notnull"`
	StartDate  time.Time       `bun:"start_date,notnull"`
	EndDate    *time.Time      `bun:"end_date"`
	SalaryType string          `bun:"salary_type,notnull"`
	BaseSalary decimal.Decimal `bun:"base_salary,type:numeric,notnull"`
	HourlyRate decimal.Decimal `bun:"hourly_rate,type:numeric,notnull"`
	Commission decimal.Decimal `bun:"commission,type:numeric,notnull"`
}

func (h *EmployeeSalaryHistory) FromDomain(domain *employeeentity.EmployeeSalaryHistory) {
	if domain == nil {
		return
	}
	*h = EmployeeSalaryHistory{
		Entity: entitymodel.FromDomain(domain.Entity),
		EmployeeSalaryHistoryCommonAttributes: EmployeeSalaryHistoryCommonAttributes{
			EmployeeID: domain.EmployeeID,
			StartDate:  domain.StartDate,
			EndDate:    domain.EndDate,
			SalaryType: domain.SalaryType,
			BaseSalary: domain.BaseSalary,
			HourlyRate: domain.HourlyRate,
			Commission: domain.Commission,
		},
	}
}

func (h *EmployeeSalaryHistory) ToDomain() *employeeentity.EmployeeSalaryHistory {
	if h == nil {
		return nil
	}
	return &employeeentity.EmployeeSalaryHistory{
		Entity:     h.Entity.ToDomain(),
		EmployeeID: h.EmployeeID,
		StartDate:  h.StartDate,
		EndDate:    h.EndDate,
		SalaryType: h.SalaryType,
		BaseSalary: h.BaseSalary,
		HourlyRate: h.HourlyRate,
		Commission: h.Commission,
	}
}
