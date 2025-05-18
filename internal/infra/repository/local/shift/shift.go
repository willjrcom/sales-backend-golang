package shiftrepositorylocal

import (
	"context"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ShiftRepositoryLocal struct {}

func NewShiftRepositoryLocal() model.ShiftRepository {
	return &ShiftRepositoryLocal{}
}

func (r *ShiftRepositoryLocal) CreateShift(ctx context.Context, shift *model.Shift) error {
	return nil
}

func (r *ShiftRepositoryLocal) UpdateShift(ctx context.Context, shift *model.Shift) error {
	return nil
}

func (r *ShiftRepositoryLocal) DeleteShift(ctx context.Context, id string) error {
	return nil
}

func (r *ShiftRepositoryLocal) GetShiftByID(ctx context.Context, id string) (*model.Shift, error) {
	return nil, nil
}

func (r *ShiftRepositoryLocal) GetCurrentShift(ctx context.Context) (*model.Shift, error) {
	return nil, nil
}

func (r *ShiftRepositoryLocal) GetAllShifts(ctx context.Context) ([]model.Shift, error) {
	return nil, nil
}

func (r *ShiftRepositoryLocal) IncrementCurrentOrder(id string) (int, error) {
	return 0, nil
}
