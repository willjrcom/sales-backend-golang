package companyusecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/cnpj"
	schemaservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
)

type Service struct {
	r companyentity.CompanyRepository
	a addressentity.Repository
	s schemaservice.Service
	u companyentity.UserRepository
}

func NewService(r companyentity.CompanyRepository, a addressentity.Repository, s schemaservice.Service, u companyentity.UserRepository) *Service {
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
	fmt.Println(company.SchemaName)
	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), company.SchemaName)

	if err := s.s.NewSchema(ctx); err != nil {
		return uuid.Nil, nil, err
	}

	if err = s.r.NewCompany(ctx, company); err != nil {
		return uuid.Nil, nil, err
	}

	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), company.SchemaName)

	userCommonAttributes := companyentity.UserCommonAttributes{
		Email:    email,
		Password: "12345",
	}

	userInput := &companydto.UserInput{UserCommonAttributes: userCommonAttributes}

	if err = s.AddUserToCompany(ctx, userInput); err != nil {
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

func (s *Service) AddUserToCompany(ctx context.Context, dto *companydto.UserInput) error {
	user, err := dto.ToModel()

	if err != nil {
		return err
	}

	userID, _ := s.u.GetIDByEmail(ctx, user.Email)

	if userID == uuid.Nil {
		userID, err = s.createUser(ctx, user.Email)

		if err != nil {
			return err
		}
	}

	if err := s.r.AddUserToPublicCompany(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveUserFromCompany(ctx context.Context, dto *companydto.UserInput) error {
	user, err := dto.ToModel()

	if err != nil {
		return err
	}

	userID, err := s.u.GetIDByEmail(ctx, user.Email)

	if err != nil {
		return err
	}

	if err := s.r.RemoveUserFromPublicCompany(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (s *Service) createUser(ctx context.Context, email string) (id uuid.UUID, err error) {
	userCommonAttributes := companyentity.UserCommonAttributes{
		Email:    email,
		Password: "12345",
	}
	user := companyentity.NewUser(userCommonAttributes)

	if err = s.u.CreateUser(ctx, user); err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}
