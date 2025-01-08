package shiftrepositorybun

import (
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	"golang.org/x/net/context"
)

type ShiftRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewShiftRepositoryBun(db *bun.DB) *ShiftRepositoryBun {
	return &ShiftRepositoryBun{db: db}
}

func (r *ShiftRepositoryBun) CreateShift(ctx context.Context, c *model.Shift) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewInsert().Model(c).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ShiftRepositoryBun) UpdateShift(ctx context.Context, c *model.Shift) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(c).Where("id = ?", c.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ShiftRepositoryBun) DeleteShift(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(&model.Shift{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ShiftRepositoryBun) GetShiftByID(ctx context.Context, id string) (*model.Shift, error) {
	shift := &model.Shift{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(shift).Where("shift.id = ?", id).Relation("Attendant").Relation("Orders").Scan(ctx); err != nil {
		return nil, err
	}

	return shift, nil
}

func (r *ShiftRepositoryBun) GetOpenedShift(ctx context.Context) (*model.Shift, error) {
	shift := &model.Shift{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(shift).Where("shift.closed_at is NULL AND shift.end_change is NULL").Relation("Attendant").Relation("Orders").Scan(ctx); err != nil {
		return nil, err
	}

	return shift, nil
}

func (r *ShiftRepositoryBun) GetAllShifts(ctx context.Context) ([]model.Shift, error) {
	Shifts := []model.Shift{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&Shifts).Scan(ctx); err != nil {
		return nil, err
	}

	return Shifts, nil
}
