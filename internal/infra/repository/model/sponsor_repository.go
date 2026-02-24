package model

import (
	"context"

	"github.com/google/uuid"
)

type SponsorRepository interface {
	Create(ctx context.Context, sponsor *Sponsor) error
	Update(ctx context.Context, sponsor *Sponsor) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Sponsor, error)
	GetAllSponsors(ctx context.Context) ([]Sponsor, error)
}
