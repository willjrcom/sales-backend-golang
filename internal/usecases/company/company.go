package companyusecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	userdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/user"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/cnpj"
	schemaservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
	userusecases "github.com/willjrcom/sales-backend-go/internal/usecases/user"
)

type Service struct {
	r  companyentity.CompanyRepository
	a  addressentity.Repository
	s  schemaservice.Service
	u  companyentity.UserRepository
	us userusecases.Service
}

func NewService(r companyentity.CompanyRepository) *Service {
	return &Service{r: r}
}

func (s *Service) AddDependencies(a addressentity.Repository, ss schemaservice.Service, u companyentity.UserRepository, us userusecases.Service) {
	s.a = a
	s.s = ss
	s.u = u
	s.us = us
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

	if userID != uuid.Nil {
		if exists, _ := s.r.ValidateUserToPublicCompany(ctx, userID); exists {
			return errors.New("user already added to company")
		}
	}

	if userID == uuid.Nil {
		userCommonAttributes := &companyentity.UserCommonAttributes{
			Email: user.Email,
		}

		createUserInput := &userdto.CreateUserInput{
			UserCommonAttributes: *userCommonAttributes,
			GeneratePassword:     true,
		}

		if userID, err = s.us.CreateUser(ctx, createUserInput); err != nil {
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

	if userID != uuid.Nil {
		if exists, _ := s.r.ValidateUserToPublicCompany(ctx, userID); !exists {
			return errors.New("user already removed from company")
		}
	}

	if err := s.r.RemoveUserFromPublicCompany(ctx, userID); err != nil {
		return err
	}

	return nil
}
