package sponsorrepositorylocal

import (
	"context"

	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type SponsorRepositoryLocal struct{}

func NewSponsorRepositoryLocal() model.SponsorRepository {
	return &SponsorRepositoryLocal{}
}

func (r *SponsorRepositoryLocal) CreateSponsor(ctx context.Context, p *model.Sponsor) error {
	return nil
}

func (r *SponsorRepositoryLocal) UpdateSponsor(ctx context.Context, p *model.Sponsor) error {
	return nil
}

func (r *SponsorRepositoryLocal) DeleteSponsor(ctx context.Context, id string) error {
	return nil
}

func (r *SponsorRepositoryLocal) GetSponsorByID(ctx context.Context, id string) (*model.Sponsor, error) {
	return nil, nil
}

func (r *SponsorRepositoryLocal) GetAllSponsors(ctx context.Context) ([]model.Sponsor, error) {
	return nil, nil
}
