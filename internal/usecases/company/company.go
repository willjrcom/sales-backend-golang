package companyusecases

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	employeedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/employee"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/cnpj"
	geocodeservice "github.com/willjrcom/sales-backend-go/internal/infra/service/geocode"
	schemaservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
	mercadopagoservice "github.com/willjrcom/sales-backend-go/internal/infra/service/mercadopago"
	employeeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/employee"
	userusecases "github.com/willjrcom/sales-backend-go/internal/usecases/user"
)

type Service struct {
	r  model.CompanyRepository
	a  model.AddressRepository
	s  schemaservice.Service
	u  model.UserRepository
	es employeeusecases.Service
	us userusecases.Service
	mp *mercadopagoservice.Client
}

const mercadoPagoProvider = "mercado_pago"

var (
	ErrMercadoPagoDisabled     = errors.New("mercado pago integration disabled")
	ErrInvalidWebhookSecret    = errors.New("invalid mercado pago webhook secret")
	ErrInvalidWebhookSignature = errors.New("invalid mercado pago webhook signature")
)

func NewService(r model.CompanyRepository, mp *mercadopagoservice.Client) *Service {
	return &Service{r: r, mp: mp}
}

func (s *Service) AddDependencies(a model.AddressRepository, ss schemaservice.Service, u model.UserRepository, us userusecases.Service, es employeeusecases.Service) {
	s.a = a
	s.s = ss
	s.u = u
	s.us = us
	s.es = es
}

func (s *Service) StartSubscriptionWatcher(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = 24 * time.Hour
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		run := func() {
			if err := s.ProcessSubscriptionExpirations(context.Background()); err != nil {
				log.Printf("subscription watcher error: %v", err)
			}
		}

		run()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				run()
			}
		}
	}()
}

func (s *Service) ProcessSubscriptionExpirations(ctx context.Context) error {
	companies, err := s.r.ListCompaniesForBilling(ctx)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, company := range companies {
		if company.SubscriptionExpiresAt == nil {
			continue
		}

		expired := company.SubscriptionExpiresAt.Before(now)
		switch {
		case expired && !company.IsBlocked:
			if err := s.r.UpdateCompanySubscription(ctx, company.ID, company.SchemaName, company.SubscriptionExpiresAt, true); err != nil {
				return err
			}
		case !expired && company.IsBlocked:
			if err := s.r.UpdateCompanySubscription(ctx, company.ID, company.SchemaName, company.SubscriptionExpiresAt, false); err != nil {
				return err
			}
		}
	}

	return nil
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
	userModels, total, err := s.r.GetCompanyUsers(ctx, page, perPage)
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

// ListCompanyPayments returns subscription payments for the authenticated company.
func (s *Service) ListCompanyPayments(ctx context.Context, page, perPage int) ([]companydto.CompanyPaymentDTO, int, error) {
	companyModel, err := s.r.GetCompany(ctx)
	if err != nil {
		return nil, 0, err
	}

	payments, total, err := s.r.ListCompanyPayments(ctx, companyModel.ID, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	dtos := make([]companydto.CompanyPaymentDTO, len(payments))
	for i := range payments {
		dto := companydto.CompanyPaymentDTO{}
		dto.FromDomain(payments[i].ToDomain())
		dtos[i] = dto
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

func (s *Service) GetSubscriptionSettings(ctx context.Context) (*companydto.SubscriptionSettingsDTO, error) {
	if s.mp == nil || !s.mp.Enabled() {
		return nil, ErrMercadoPagoDisabled
	}

	return &companydto.SubscriptionSettingsDTO{
		MonthlyPrice:  s.mp.MonthlyPrice(),
		Currency:      "BRL",
		MinMonths:     1,
		MaxMonths:     12,
		DefaultMonths: 1,
	}, nil
}

func (s *Service) CreateSubscriptionCheckout(ctx context.Context, dto *companydto.SubscriptionCheckoutDTO) (*companydto.SubscriptionCheckoutResponseDTO, error) {
	if s.mp == nil || !s.mp.Enabled() {
		return nil, ErrMercadoPagoDisabled
	}

	months := dto.Normalize()
	companyModel, err := s.r.GetCompany(ctx)
	if err != nil {
		return nil, err
	}

	company := companyModel.ToDomain()
	req := &mercadopagoservice.PreferenceRequest{
		Title:   "Mensalidade",
		Company: company.TradeName,
		Months:  months,
		Price:   s.mp.MonthlyPrice(),
		Schema:  company.SchemaName,
		ID:      company.ID.String(),
	}

	pref, err := s.mp.CreateSubscriptionPreference(ctx, req)
	if err != nil {
		return nil, err
	}

	return &companydto.SubscriptionCheckoutResponseDTO{
		PreferenceID:     pref.ID,
		InitPoint:        pref.InitPoint,
		SandboxInitPoint: pref.SandboxInitPoint,
	}, nil
}

func (s *Service) HandleMercadoPagoWebhook(ctx context.Context, dto *companydto.MercadoPagoWebhookDTO, signatureHeader string, payload []byte) error {
	if s.mp == nil || !s.mp.Enabled() {
		return ErrMercadoPagoDisabled
	}

	if err := s.mp.ValidateWebhookSignature(signatureHeader, payload); err != nil {
		if errors.Is(err, mercadopagoservice.ErrWebhookSecretNotConfigured) {
			return ErrInvalidWebhookSecret
		}
		return ErrInvalidWebhookSignature
	}

	if dto == nil || dto.Type != "payment" || dto.Data.ID == "" {
		return nil
	}

	paymentID := dto.Data.ID
	if existing, err := s.r.GetCompanyPaymentByProviderID(ctx, mercadoPagoProvider, paymentID); err == nil && existing != nil {
		return nil
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	details, err := s.mp.GetPayment(ctx, paymentID)
	if err != nil {
		return err
	}

	if details.Status != "approved" {
		return nil
	}

	companyIDRef := details.Metadata.CompanyID
	if companyIDRef == "" {
		companyIDRef = details.ExternalReference
	}

	companyID, err := uuid.Parse(companyIDRef)
	if err != nil {
		return err
	}

	companyModel, err := s.r.GetCompanyByIDPublic(ctx, companyID)
	if err != nil {
		return err
	}

	months := details.Metadata.Months
	if months <= 0 {
		months = 1
	}

	paidAt := time.Now().UTC()
	if details.DateApproved != nil {
		paidAt = details.DateApproved.UTC()
	}

	base := paidAt
	if companyModel.SubscriptionExpiresAt != nil && companyModel.SubscriptionExpiresAt.After(paidAt) {
		base = *companyModel.SubscriptionExpiresAt
	}

	newExpiration := base.AddDate(0, months, 0)

	if err := s.r.UpdateCompanySubscription(ctx, companyModel.ID, companyModel.SchemaName, &newExpiration, false); err != nil {
		return err
	}

	var rawPayload []byte
	switch {
	case len(payload) > 0:
		rawPayload = append([]byte(nil), payload...)
	default:
		rawPayload, _ = json.Marshal(dto)
	}

	payment := &companyentity.SubscriptionPayment{
		Entity:            entity.NewEntity(),
		CompanyID:         companyModel.ID,
		Provider:          mercadoPagoProvider,
		ProviderPaymentID: paymentID,
		Status:            details.Status,
		Currency:          details.CurrencyID,
		Amount:            decimal.NewFromFloat(details.TransactionAmount),
		Months:            months,
		PaidAt:            paidAt,
		ExternalReference: details.ExternalReference,
		RawPayload:        rawPayload,
	}

	paymentModel := &model.CompanyPayment{}
	paymentModel.FromDomain(payment)

	return s.r.CreateCompanyPayment(ctx, paymentModel)
}

func (s *Service) Test(ctx context.Context) error {
	//go kafka.ReadMessages("order_process")
	return nil
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
