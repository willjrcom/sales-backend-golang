package companyusecases

import (
	"context"

	"github.com/google/uuid"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/cnpj"
	schemaservice "github.com/willjrcom/sales-backend-go/internal/infra/service/headerservice"
)

type Service struct {
	r companyentity.Repository
	a addressentity.Repository
	s schemaservice.Service
}

func NewService(r companyentity.Repository, a addressentity.Repository, s schemaservice.Service) *Service {
	return &Service{r: r, a: a, s: s}
}

func (s *Service) NewCompany(ctx context.Context, dto *companydto.CompanyInput) (id uuid.UUID, schemaName *string, err error) {
	cnpjString, tradeName, email, contacts, err := dto.ToModel()
	if err != nil {
		return uuid.Nil, nil, err
	}

	cnpjData, err := cnpj.Get(cnpjString)

	if err != nil {
		return uuid.Nil, nil, err
	}

	if tradeName != cnpjData.TradeName {
		cnpjData.TradeName = tradeName
	}

	company := companyentity.NewCompany(cnpjData)
	company.Email = email
	company.Contacts = contacts

	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), company.SchemaName)

	if err := s.s.NewSchema(ctx); err != nil {
		return uuid.Nil, nil, err
	}

	if err := s.r.NewCompany(ctx, company); err != nil {
		return uuid.Nil, nil, err
	}

	return company.ID, &company.SchemaName, nil
}

func (s *Service) GetCompany(ctx context.Context) (*companydto.CompanyOutput, error) {
	if company, err := s.r.GetCompany(ctx); err != nil {
		return nil, err
	} else {
		output := &companydto.CompanyOutput{}
		output.FromModel(company)
		return output, nil
	}
}
