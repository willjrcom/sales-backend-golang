package companyusecases

import (
	"context"

	"github.com/google/uuid"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
	userentity "github.com/willjrcom/sales-backend-go/internal/domain/user"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/cnpj"
	schemaservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
)

type Service struct {
	r companyentity.Repository
	a addressentity.Repository
	s schemaservice.Service
	u userentity.Repository
}

func NewService(r companyentity.Repository, a addressentity.Repository, s schemaservice.Service, u userentity.Repository) *Service {
	return &Service{r: r, a: a, s: s, u: u}
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

	publicCompanyID, err := s.r.NewCompany(ctx, company)

	if err != nil {
		return uuid.Nil, nil, err
	}

	userID, err := s.u.GetIDByEmail(ctx, company.Email)

	if err != nil {
		return uuid.Nil, nil, err
	}

	if userID == uuid.Nil {
		userID, err = s.newUser(ctx, company.Email)

		if err != nil {
			return uuid.Nil, nil, err
		}
	}

	if err := s.r.AddUser(ctx, publicCompanyID, userID); err != nil {
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

func (s *Service) newUser(ctx context.Context, email string) (id uuid.UUID, err error) {
	userCommonAttributes := userentity.UserCommonAttributes{
		Email:    email,
		Password: "12345",
	}
	user := userentity.NewUser(userCommonAttributes)

	if err = s.u.CreateUser(ctx, user); err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}
