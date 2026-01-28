package companyusecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	employeedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/employee"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/cnpj"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/focusnfe"
	geocodeservice "github.com/willjrcom/sales-backend-go/internal/infra/service/geocode"
	schemaservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
	employeeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/employee"
	userusecases "github.com/willjrcom/sales-backend-go/internal/usecases/user"
)

type Service struct {
	r                  model.CompanyRepository
	companyPaymentRepo model.CompanyPaymentRepository
	a                  model.AddressRepository
	s                  schemaservice.Service
	u                  model.UserRepository
	es                 employeeusecases.Service
	us                 userusecases.Service
	focusClient        *focusnfe.Client
}

func NewService(r model.CompanyRepository, companyPaymentRepo model.CompanyPaymentRepository, focusClient *focusnfe.Client) *Service {
	return &Service{r: r, companyPaymentRepo: companyPaymentRepo, focusClient: focusClient}
}

func (s *Service) AddDependencies(a model.AddressRepository, ss schemaservice.Service, u model.UserRepository, us userusecases.Service, es employeeusecases.Service) {
	s.a = a
	s.s = ss
	s.u = u
	s.us = us
	s.es = es
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

	// Define 7 days free trial
	expiration := time.Now().UTC().AddDate(0, 0, 7)
	company.SubscriptionExpiresAt = &expiration

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

	employeeInput := &employeedto.EmployeeCreateDTO{
		UserID: &userIDUUID,
	}

	if _, err = s.es.CreateEmployee(ctx, employeeInput); err != nil {
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

	if dto.MonthlyPaymentDueDay != nil && *dto.MonthlyPaymentDueDay != company.MonthlyPaymentDueDay {
		// Validate range
		if *dto.MonthlyPaymentDueDay < 1 || *dto.MonthlyPaymentDueDay > 28 {
			return errors.New("payment due day must be between 1 and 28")
		}

		// Check 3-month restriction
		if company.MonthlyPaymentDueDayUpdatedAt != nil {
			nextAllowedUpdate := company.MonthlyPaymentDueDayUpdatedAt.AddDate(0, 3, 0)
			if time.Now().Before(nextAllowedUpdate) {
				return errors.New("payment due day can only be changed every 3 months")
			}
		}
		now := time.Now()
		company.MonthlyPaymentDueDayUpdatedAt = &now
	}

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

// ListPublicCompanies returns basic information for every company stored in the public schema.
func (s *Service) ListPublicCompanies(ctx context.Context) ([]companydto.CompanyBasicDTO, error) {
	companyModels, err := s.r.ListPublicCompanies(ctx)
	if err != nil {
		return nil, err
	}

	if len(companyModels) == 0 {
		return []companydto.CompanyBasicDTO{}, nil
	}

	basic := make([]companydto.CompanyBasicDTO, len(companyModels))
	for i := range companyModels {
		company := companyModels[i].ToDomain()
		dto := companydto.CompanyBasicDTO{}
		dto.FromDomain(company)
		basic[i] = dto
	}

	return basic, nil
}
