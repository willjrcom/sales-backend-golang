package sponsorrepositorylocal

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type SponsorRepositoryLocal struct {
	mu       sync.RWMutex
	sponsors map[string]*model.Sponsor
}

func NewSponsorRepositoryLocal() model.SponsorRepository {
	return &SponsorRepositoryLocal{sponsors: make(map[string]*model.Sponsor)}
}

func (r *SponsorRepositoryLocal) Create(ctx context.Context, p *model.Sponsor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sponsors[p.ID.String()] = p
	return nil
}

func (r *SponsorRepositoryLocal) Update(ctx context.Context, p *model.Sponsor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sponsors[p.ID.String()] = p
	return nil
}

func (r *SponsorRepositoryLocal) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sponsors, id.String())
	return nil
}

func (r *SponsorRepositoryLocal) GetByID(ctx context.Context, id uuid.UUID) (*model.Sponsor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if s, ok := r.sponsors[id.String()]; ok {
		return s, nil
	}
	return nil, nil
}

func (r *SponsorRepositoryLocal) GetAllSponsors(ctx context.Context) ([]model.Sponsor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]model.Sponsor, 0, len(r.sponsors))
	for _, s := range r.sponsors {
		out = append(out, *s)
	}
	return out, nil
}
