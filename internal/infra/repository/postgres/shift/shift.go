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

func NewShiftRepositoryBun(db *bun.DB) model.ShiftRepository {
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

func (r *ShiftRepositoryBun) GetCurrentShift(ctx context.Context) (*model.Shift, error) {
	shift := &model.Shift{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(shift).Where("shift.closed_at is NULL AND shift.end_change is NULL").Scan(ctx); err != nil {
		return nil, err
	}

	return shift, nil
}

func (r *ShiftRepositoryBun) GetFullCurrentShift(ctx context.Context) (*model.Shift, error) {
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

func (s *ShiftRepositoryBun) IncrementCurrentOrder(id string) (int, error) {
	ctx := context.Background()

	// Inicia a transação
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Busca o turno com "FOR UPDATE" para bloquear a linha
	shift := new(model.Shift)
	err = tx.NewSelect().
		Model(shift).
		Where("id = ?", id).
		For("UPDATE"). // bloqueia a linha durante a transação
		Scan(ctx)
	if err != nil {
		return 0, err
	}

	// Incrementa o número do pedido
	shift.CurrentOrderNumber++

	// Atualiza o registro
	_, err = tx.NewUpdate().
		Model(shift).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	// Commit da transação
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return shift.CurrentOrderNumber, nil
}
