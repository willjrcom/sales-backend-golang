package shiftrepositorylocal

import (
	"context"
	"sync"

	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ShiftRepositoryLocal struct {
	mu     sync.RWMutex
	shifts map[string]*model.Shift
}

func NewShiftRepositoryLocal() model.ShiftRepository {
	return &ShiftRepositoryLocal{shifts: make(map[string]*model.Shift)}
}

func (r *ShiftRepositoryLocal) CreateShift(ctx context.Context, shift *model.Shift) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.shifts[shift.ID.String()] = shift
	return nil
}

func (r *ShiftRepositoryLocal) UpdateShift(ctx context.Context, shift *model.Shift) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.shifts[shift.ID.String()] = shift
	return nil
}

func (r *ShiftRepositoryLocal) DeleteShift(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.shifts, id)
	return nil
}

func (r *ShiftRepositoryLocal) GetShiftByID(ctx context.Context, id string) (*model.Shift, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if s, ok := r.shifts[id]; ok {
		return s, nil
	}
	return nil, nil
}

func (r *ShiftRepositoryLocal) GetCurrentShift(ctx context.Context) (*model.Shift, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.shifts {
		if s.EndChange == nil {
			return s, nil
		}
	}
	return nil, nil
}

func (r *ShiftRepositoryLocal) GetFullCurrentShift(ctx context.Context) (*model.Shift, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.shifts {
		if s.EndChange == nil {
			return s, nil
		}
	}
	return nil, nil
}

func (r *ShiftRepositoryLocal) GetAllShifts(ctx context.Context, page int, perPage int) ([]model.Shift, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]model.Shift, 0, len(r.shifts))
	for _, s := range r.shifts {
		list = append(list, *s)
	}
	return list, nil
}

func (r *ShiftRepositoryLocal) IncrementCurrentOrder(ctx context.Context, id string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s, ok := r.shifts[id]; ok {
		s.CurrentOrderNumber++
		return s.CurrentOrderNumber, nil
	}
	return 0, nil
}
