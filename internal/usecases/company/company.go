package companyusecases

import (
	"context"

	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/cnpj"
	schemaservice "github.com/willjrcom/sales-backend-go/internal/infra/service/schema"
)

type Service struct {
	r companyentity.Repository
	s schemaservice.Service
}

func NewService(r companyentity.Repository, s schemaservice.Service) *Service {
	return &Service{r: r, s: s}
}

func (s *Service) NewCompany(ctx context.Context, dto *companydto.CompanyInput) (uuid.UUID, error) {
	cnpjString, email, contacts, err := dto.ToModel()
	if err != nil {
		return uuid.Nil, err
	}

	cnpjData, err := cnpj.Get(cnpjString)

	if err != nil {
		return uuid.Nil, err
	}

	company := companyentity.NewCompany(cnpjData)
	company.Email = email
	company.Contacts = contacts

	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), company.SchemaName)

	if err := s.s.NewSchema(ctx); err != nil {
		return uuid.Nil, err
	}

	if err := s.r.NewCompany(ctx, company); err != nil {
		return uuid.Nil, err
	}

	return company.ID, nil
}

func (s *Service) GetCompanyById(ctx context.Context, dto *entitydto.IdRequest) (*companydto.CompanyOutput, error) {
	if company, err := s.r.GetCompanyById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		output := &companydto.CompanyOutput{}
		output.FromModel(company)
		return output, nil
	}
}

func (s *Service) GetAllCompaniesBySchemaName(ctx context.Context, dto *companydto.CompanyBySchemaName) ([]companydto.CompanyOutput, error) {
	schemaName, err := dto.ToModel()

	if err != nil {
		return nil, err
	}

	if companys, err := s.r.GetAllCompaniesBySchemaName(ctx, schemaName); err != nil {
		return nil, err
	} else {
		dtos := companiesToDtos(companys)
		return dtos, nil
	}
}

func companiesToDtos(companys []companyentity.Company) []companydto.CompanyOutput {
	dtos := make([]companydto.CompanyOutput, len(companys))
	for i, company := range companys {
		dtos[i].FromModel(&company)
	}

	return dtos
}
