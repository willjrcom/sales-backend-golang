package model

import "context"

type SponsorRepository interface {
	CreateSponsor(ctx context.Context, Sponsor *Sponsor) (err error)
	UpdateSponsor(ctx context.Context, Sponsor *Sponsor) (err error)
	DeleteSponsor(ctx context.Context, id string) (err error)
	GetSponsorByID(ctx context.Context, id string) (Sponsor *Sponsor, err error)
	GetAllSponsors(ctx context.Context) ([]Sponsor, error)
}
