package employeerepositorybun

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type EmployeeSalaryHistoryRepositoryBun struct {
	db *bun.DB
}

func NewEmployeeSalaryHistoryRepositoryBun(db *bun.DB) *EmployeeSalaryHistoryRepositoryBun {
	return &EmployeeSalaryHistoryRepositoryBun{db: db}
}

func (r *EmployeeSalaryHistoryRepositoryBun) Create(ctx context.Context, h *model.EmployeeSalaryHistory) error {
	_, err := r.db.NewInsert().Model(h).Exec(ctx)
	return err
}

func (r *EmployeeSalaryHistoryRepositoryBun) GetByEmployee(ctx context.Context, employeeID uuid.UUID) ([]model.EmployeeSalaryHistory, error) {
	var history []model.EmployeeSalaryHistory
	err := r.db.NewSelect().Model(&history).Where("employee_id = ?", employeeID).Order("start_date DESC").Scan(ctx)
	return history, err
}

func (r *EmployeeSalaryHistoryRepositoryBun) GetCurrentByEmployee(ctx context.Context, employeeID uuid.UUID) (*model.EmployeeSalaryHistory, error) {
	var h model.EmployeeSalaryHistory
	err := r.db.NewSelect().Model(&h).
		Where("employee_id = ?", employeeID).
		Where("end_date IS NULL").
		Order("start_date DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *EmployeeSalaryHistoryRepositoryBun) EndCurrent(ctx context.Context, employeeID uuid.UUID, endDate time.Time) error {
	_, err := r.db.NewUpdate().
		Model((*model.EmployeeSalaryHistory)(nil)).
		Set("end_date = ?", endDate).
		Where("employee_id = ?", employeeID).
		Where("end_date IS NULL").
		Exec(ctx)
	return err
}
