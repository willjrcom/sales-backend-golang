package companyusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	userdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/user"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/cnpj"
	geocodeservice "github.com/willjrcom/sales-backend-go/internal/infra/service/geocode"
	schemaservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
	userusecases "github.com/willjrcom/sales-backend-go/internal/usecases/user"
)

type Service struct {
	r  model.CompanyRepository
	a  model.AddressRepository
	s  schemaservice.Service
	u  model.UserRepository
	us userusecases.Service
}

func NewService(r model.CompanyRepository) *Service {
	return &Service{r: r}
}

func (s *Service) AddDependencies(a model.AddressRepository, ss schemaservice.Service, u model.UserRepository, us userusecases.Service) {
	s.a = a
	s.s = ss
	s.u = u
	s.us = us
}

func (s *Service) NewCompany(ctx context.Context, dto *companydto.CompanyCreateDTO) (id uuid.UUID, schemaName *string, err error) {
	cnpjString, tradeName, email, contacts, err := dto.ToDomain()
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

	coordinates, _ := geocodeservice.GetCoordinates(&company.Address.AddressCommonAttributes)

	if coordinates != nil {
		company.Address.Coordinates = *coordinates
	}

	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), company.SchemaName)

	if err := s.s.NewSchema(ctx); err != nil {
		return uuid.Nil, nil, err
	}

	if err = s.r.NewCompany(ctx, company); err != nil {
		return uuid.Nil, nil, err
	}

	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), company.SchemaName)

	userInput := &companydto.UserToCompanyDTO{
		Email: email,
	}

	if err = s.AddUserToCompany(ctx, userInput); err != nil {
		return uuid.Nil, nil, err
	}

	return company.ID, &company.SchemaName, nil
}

func (s *Service) GetCompany(ctx context.Context) (*companydto.CompanyDTO, error) {
	if company, err := s.r.GetCompany(ctx); err != nil {
		return nil, err
	} else {
		output := &companydto.CompanyDTO{}
		output.FromDomain(company)
		return output, nil
	}
}

func (s *Service) AddUserToCompany(ctx context.Context, dto *companydto.UserToCompanyDTO) error {
	email, err := dto.ToDomain()

	if err != nil {
		return err
	}

	userID, _ := s.u.GetIDByEmail(ctx, email)

	if userID != nil {
		if exists, _ := s.r.ValidateUserToPublicCompany(ctx, *userID); exists {
			return errors.New("user already added to company")
		}
	}

	if userID == nil {
		createUserInput := &userdto.UserCreateDTO{
			Email:            email,
			GeneratePassword: true,
		}

		if newUserID, err := s.us.CreateUser(ctx, createUserInput); err != nil {
			return err
		} else {
			userID = newUserID
		}
	}

	if err := s.r.AddUserToPublicCompany(ctx, *userID); err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveUserFromCompany(ctx context.Context, dto *companydto.UserToCompanyDTO) error {
	email, err := dto.ToDomain()

	if err != nil {
		return err
	}

	userID, err := s.u.GetIDByEmail(ctx, email)

	if err != nil {
		return err
	}

	if userID != nil {
		if exists, _ := s.r.ValidateUserToPublicCompany(ctx, *userID); !exists {
			return errors.New("user already removed from company")
		}
	}

	if err := s.r.RemoveUserFromPublicCompany(ctx, *userID); err != nil {
		return err
	}

	return nil
}

func (s *Service) Test(ctx context.Context) error {
	//go kafka.ReadMessages("order_process")
	return nil
}
