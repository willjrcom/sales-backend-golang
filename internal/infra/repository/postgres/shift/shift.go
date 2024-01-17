package contactrepositorybun

import (
	"sync"

	"github.com/uptrace/bun"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
	"golang.org/x/net/context"
)

type ShiftRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewShiftRepositoryBun(ctx context.Context, db *bun.DB) *ShiftRepositoryBun {
	return &ShiftRepositoryBun{db: db}
}

func (r *ShiftRepositoryBun) SaveShift(ctx context.Context, c *shiftentity.Shift) error {
	r.mu.Lock()
	_, err := r.db.NewInsert().Model(c).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ShiftRepositoryBun) DeleteShift(ctx context.Context, id string) error {
	r.mu.Lock()
	_, err := r.db.NewUpdate().Model(&shiftentity.Shift{}).Where("id = ?", id).Exec(ctx)
	r.mu.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (r *ShiftRepositoryBun) GetShiftById(ctx context.Context, id string) (*shiftentity.Shift, error) {
	contact := &shiftentity.Shift{}

	r.mu.Lock()
	err := r.db.NewSelect().Model(contact).Where("shift.id = ?", id).Scan(ctx)
	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (r *ShiftRepositoryBun) GetAllShifts(ctx context.Context) ([]shiftentity.Shift, error) {
	Shifts := []shiftentity.Shift{}
	r.mu.Lock()
	err := r.db.NewSelect().Model(&Shifts).Scan(ctx)

	r.mu.Unlock()

	if err != nil {
		return nil, err
	}

	return Shifts, nil
}
