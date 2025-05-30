package companyusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
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

func (s *Service) NewCompany(ctx context.Context, dto *companydto.CompanyCreateDTO) (response *companydto.CompanySchemaDTO, err error) {
	cnpjString, tradeName, contacts, err := dto.ToDomain()
	if err != nil {
		return nil, err
	}

	cnpjData, err := cnpj.Get(cnpjString)

	if err != nil {
		return nil, err
	}

	if tradeName != cnpjData.TradeName {
		cnpjData.TradeName = tradeName
	}

	userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)

	if !ok {
		return nil, errors.New("context user not found")
	}

	userIDUUID := uuid.MustParse(userID)
	userModel, err := s.u.GetUserByID(ctx, userIDUUID)
	if err != nil {
		return nil, err
	}

	company := companyentity.NewCompany(cnpjData)
	company.Email = userModel.Email
	company.Contacts = contacts

	coordinates, _ := geocodeservice.GetCoordinates(&company.Address.AddressCommonAttributes)

	if coordinates != nil {
		company.Address.Coordinates = *coordinates
	}

	ctx = context.WithValue(ctx, model.Schema("schema"), company.SchemaName)

	if err := s.s.NewSchema(ctx); err != nil {
		return nil, err
	}

	companyModel := &model.Company{}
	companyModel.FromDomain(company)
	if err = s.r.NewCompany(ctx, companyModel); err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, schemaentity.Schema("schema"), company.SchemaName)

	userInput := &companydto.UserToCompanyDTO{
		Email: userModel.Email,
	}

	if err = s.AddUserToCompany(ctx, userInput); err != nil {
		return nil, err
	}

	companyDTO := &companydto.CompanySchemaDTO{}
	companyDTO.FromDomain(company)
	return companyDTO, nil
}

func (s *Service) UpdateCompany(ctx context.Context, dto *companydto.CompanyUpdateDTO) (err error) {
	companyModel, err := s.r.GetCompany(ctx)
	if err != nil {
		return err
	}

	company := companyModel.ToDomain()

	dto.UpdateDomain(company)

	if company.Cnpj != companyModel.Cnpj {
		cnpjData, err := cnpj.Get(company.Cnpj)
		if err != nil {
			return err
		}

		if dto.TradeName != nil {
			cnpjData.TradeName = *dto.TradeName
		}

		company.UpdateCompany(cnpjData)
	}

	coordinates, _ := geocodeservice.GetCoordinates(&company.Address.AddressCommonAttributes)

	if coordinates != nil {
		company.Address.Coordinates = *coordinates
	}

	companyModel.FromDomain(company)
	if err = s.r.UpdateCompany(ctx, companyModel); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetCompany(ctx context.Context) (*companydto.CompanyDTO, error) {
	if companyModel, err := s.r.GetCompany(ctx); err != nil {
		return nil, err
	} else {
		company := companyModel.ToDomain()
		output := &companydto.CompanyDTO{}
		output.FromDomain(company)
		return output, nil
	}
}

// GetCompanyUsers retrieves a paginated list of users and the total count.
func (s *Service) GetCompanyUsers(ctx context.Context, page, perPage int) ([]companydto.UserDTO, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	userModels, total, err := s.r.GetCompanyUsers(ctx, offset, perPage)
	if err != nil {
		return nil, 0, err
	}
	dtos := make([]companydto.UserDTO, len(userModels))
	for i, userModel := range userModels {
		user := userModel.ToDomain()
		dto := &companydto.UserDTO{}
		dto.FromDomain(user)
		dtos[i] = *dto
	}
	return dtos, total, nil
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
		return errors.New("user not found")
		// createUserInput := &companydto.UserCreateDTO{
		// 	Email:            email,
		// 	GeneratePassword: true,
		// }

		// if newUserID, err := s.us.CreateUser(ctx, createUserInput); err != nil {
		// 	return err
		// } else {
		// 	userID = newUserID
		// }
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
