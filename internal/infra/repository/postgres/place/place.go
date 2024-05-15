package placerepositorybun

import (
	"context"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
)

type PlaceRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewPlaceRepositoryBun(db *bun.DB) *PlaceRepositoryBun {
	return &PlaceRepositoryBun{db: db}
}

func (r *PlaceRepositoryBun) CreatePlace(ctx context.Context, s *tableentity.Place) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewInsert().Model(s).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PlaceRepositoryBun) UpdatePlace(ctx context.Context, s *tableentity.Place) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(s).Where("id = ?", s.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PlaceRepositoryBun) DeletePlace(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewDelete().Model(&tableentity.Place{}).Where("id = ?", id).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PlaceRepositoryBun) GetPlaceById(ctx context.Context, id string) (*tableentity.Place, error) {
	place := &tableentity.Place{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(place).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return place, nil
}

func (r *PlaceRepositoryBun) GetAllPlaces(ctx context.Context) ([]tableentity.Place, error) {
	places := make([]tableentity.Place, 0)

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&places).Scan(ctx); err != nil {
		return nil, err
	}

	return places, nil
}
