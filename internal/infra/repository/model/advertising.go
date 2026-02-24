package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	advertisingentity "github.com/willjrcom/sales-backend-go/internal/domain/advertising"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Advertising struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:public.advertisements"`
	AdvertisingCommonAttributes
}

type AdvertisingCommonAttributes struct {
	Title          string            `bun:"title,notnull"`
	Description    string            `bun:"description"`
	Link           string            `bun:"link"`
	Contact        string            `bun:"contact"`
	CoverImagePath string            `bun:"cover_image_path"`
	Images         []string          `bun:"images,type:jsonb"`
	SponsorID      uuid.UUID         `bun:"sponsor_id,type:uuid"`
	Sponsor        *Sponsor          `bun:"rel:belongs-to,join:sponsor_id=id"`
	Categories     []CompanyCategory `bun:"m2m:public.category_advertisements,join:Advertising=Category"`
}

func (a *Advertising) FromDomain(ad *advertisingentity.Advertising) {
	if ad == nil {
		return
	}
	*a = Advertising{
		Entity: entitymodel.FromDomain(ad.Entity),
		AdvertisingCommonAttributes: AdvertisingCommonAttributes{
			Title:          ad.Title,
			Description:    ad.Description,
			Link:           ad.Link,
			Contact:        ad.Contact,
			CoverImagePath: ad.CoverImagePath,
			Images:         ad.Images,
			SponsorID:      ad.SponsorID,
		},
	}
	if ad.Sponsor != nil {
		sponsorModel := &Sponsor{}
		sponsorModel.FromDomain(ad.Sponsor)
		a.Sponsor = sponsorModel
	}

	for _, cat := range ad.CompanyCategoryAdvertising {
		id, err := uuid.Parse(cat.ID)
		if err == nil {
			c := CompanyCategory{}
			c.ID = id
			a.Categories = append(a.Categories, c)
		}
	}
}

func (a *Advertising) ToDomain() *advertisingentity.Advertising {
	if a == nil {
		return nil
	}
	ad := &advertisingentity.Advertising{
		Entity: a.Entity.ToDomain(),
		AdvertisingCommonAttributes: advertisingentity.AdvertisingCommonAttributes{
			Title:          a.Title,
			Description:    a.Description,
			Link:           a.Link,
			Contact:        a.Contact,
			CoverImagePath: a.CoverImagePath,
			Images:         a.Images,
			SponsorID:      a.SponsorID,
		},
	}

	cats := make([]advertisingentity.CompanyCategory, len(a.Categories))
	for i, cat := range a.Categories {
		cats[i] = advertisingentity.CompanyCategory{ID: cat.ID.String()}
	}
	ad.CompanyCategoryAdvertising = cats

	if a.Sponsor != nil {
		ad.Sponsor = a.Sponsor.ToDomain()
	}
	return ad
}

type AdvertisingRepository interface {
	Create(ctx context.Context, advertising *Advertising) error
	Update(ctx context.Context, advertising *Advertising) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*Advertising, error)
	GetAllAdvertisements(ctx context.Context) ([]Advertising, error)
}
