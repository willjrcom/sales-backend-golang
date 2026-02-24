package model

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	sponsorentity "github.com/willjrcom/sales-backend-go/internal/domain/sponsor"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

type Sponsor struct {
	entitymodel.Entity
	bun.BaseModel `bun:"table:public.sponsors"`
	SponsorCommonAttributes
}

type SponsorCommonAttributes struct {
	Name       string            `bun:"name,notnull"`
	CNPJ       string            `bun:"cnpj,notnull"`
	Email      string            `bun:"email"`
	Contact    string            `bun:"contact"`
	Address    *Address          `bun:"rel:has-one,join:id=object_id,notnull"`
	Categories []CompanyCategory `bun:"m2m:public.category_sponsors,join:Sponsor=Category"`
}

func (m *Sponsor) FromDomain(s *sponsorentity.Sponsor) {
	if s == nil {
		return
	}
	*m = Sponsor{
		Entity: entitymodel.FromDomain(s.Entity),
		SponsorCommonAttributes: SponsorCommonAttributes{
			Name:    s.Name,
			CNPJ:    s.CNPJ,
			Email:   s.Email,
			Contact: s.Contact,
		},
	}
	if s.Address != nil {
		addressModel := &Address{}
		addressModel.FromDomain(s.Address)
		m.Address = addressModel
	}

	for _, cat := range s.CompanyCategorySponsor {
		id, err := uuid.Parse(cat.ID)
		if err == nil {
			c := CompanyCategory{}
			c.ID = id
			m.Categories = append(m.Categories, c)
		}
	}
}

func (m *Sponsor) ToDomain() *sponsorentity.Sponsor {
	if m == nil {
		return nil
	}
	s := &sponsorentity.Sponsor{
		Entity: m.Entity.ToDomain(),
		SponsorCommonAttributes: sponsorentity.SponsorCommonAttributes{
			Name:    m.Name,
			CNPJ:    m.CNPJ,
			Email:   m.Email,
			Contact: m.Contact,
		},
	}

	cats := make([]sponsorentity.CompanyCategory, len(m.Categories))
	for i, cat := range m.Categories {
		cats[i] = sponsorentity.CompanyCategory{ID: cat.ID.String()}
	}
	s.CompanyCategorySponsor = cats

	if m.Address != nil {
		s.Address = m.Address.ToDomain()
	}
	return s
}
