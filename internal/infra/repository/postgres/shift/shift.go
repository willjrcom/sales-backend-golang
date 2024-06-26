package shiftrepositorybun

import (
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
	"golang.org/x/net/context"
)

type ShiftRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewShiftRepositoryBun(db *bun.DB) *ShiftRepositoryBun {
	return &ShiftRepositoryBun{db: db}
}

func (r *ShiftRepositoryBun) CreateShift(ctx context.Context, c *shiftentity.Shift) error {
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

func (r *ShiftRepositoryBun) UpdateShift(ctx context.Context, c *shiftentity.Shift) error {
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

	if _, err := r.db.NewUpdate().Model(&shiftentity.Shift{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ShiftRepositoryBun) GetShiftByID(ctx context.Context, id string) (*shiftentity.Shift, error) {
	shift := &shiftentity.Shift{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(shift).Where("shift.id = ?", id).Relation("Attendant").Scan(ctx); err != nil {
		return nil, err
	}

	return shift, nil
}

func (r *ShiftRepositoryBun) GetOpenedShift(ctx context.Context) (*shiftentity.Shift, error) {
	shift := &shiftentity.Shift{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(shift).Where("shift.closed_at is NULL AND shift.end_change is NULL").Relation("Attendant").Scan(ctx); err != nil {
		return nil, err
	}

	return shift, nil
}

func (r *ShiftRepositoryBun) GetAllShifts(ctx context.Context) ([]shiftentity.Shift, error) {
	Shifts := []shiftentity.Shift{}

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
