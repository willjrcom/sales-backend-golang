package billing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	domainbilling "github.com/willjrcom/sales-backend-go/internal/domain/checkout"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	billingdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/checkout"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	mercadopagoservice "github.com/willjrcom/sales-backend-go/internal/infra/service/mercadopago"
)

type CheckoutUseCase struct {
	costRepo           model.CompanyUsageCostRepository
	companyRepo        model.CompanyRepository
	companyPaymentRepo model.CompanyPaymentRepository
	mpService          *mercadopagoservice.Client
}

func NewCheckoutUseCase(
	costRepo model.CompanyUsageCostRepository,
	companyRepo model.CompanyRepository,
	companyPaymentRepo model.CompanyPaymentRepository,
	mpService *mercadopagoservice.Client,
) *CheckoutUseCase {
	return &CheckoutUseCase{
		costRepo:           costRepo,
		companyRepo:        companyRepo,
		companyPaymentRepo: companyPaymentRepo,
		mpService:          mpService,
	}
}

const mercadoPagoProvider = "mercado_pago"

var (
	ErrMercadoPagoDisabled  = errors.New("mercado pago integration disabled")
	ErrInvalidWebhookSecret = errors.New("invalid mercado pago webhook secret")
)

func (uc *CheckoutUseCase) CreateSubscriptionCheckout(ctx context.Context, req *billingdto.CreateCheckoutDTO) (*billingdto.CheckoutResponseDTO, error) {
	var basePrice decimal.Decimal
	switch req.ToPlanType() {
	case domainbilling.PlanIntermediate:
		price := getEnvFloat("PRICE_INTERMEDIATE", 119.90)
		basePrice = decimal.NewFromFloat(price)
	case domainbilling.PlanEnterprise:
		price := getEnvFloat("PRICE_ENTERPRISE", 129.90)
		basePrice = decimal.NewFromFloat(price)
	default:
		price := getEnvFloat("PRICE_BASIC", 99.90)
		basePrice = decimal.NewFromFloat(price)
	}

	months := 1
	discount := 0.0
	switch req.ToPeriodicity() {
	case domainbilling.PeriodicitySemiannual:
		months = 6
		discount = 0.05
	case domainbilling.PeriodicityAnnual:
		months = 12
		discount = 0.10
	}

	totalAmount := basePrice.Mul(decimal.NewFromInt(int64(months)))
	finalAmount := totalAmount.Mul(decimal.NewFromFloat(1.0 - discount))
	startAt := time.Now()
	endAt := startAt.AddDate(0, months, 0)

	title := fmt.Sprintf("Mensalidade Gfood Plano %s - %s", translatePlanType(req.ToPlanType()), translatePeriodicity(req.ToPeriodicity()))
	description := fmt.Sprintf("Periodo %s - %s", startAt.Format("02/01/2006"), endAt.Format("02/01/2006"))

	// 2. Prepare Checkout Item
	checkoutItem := mercadopagoservice.NewCheckoutItem(
		title,
		description,
		1,
		finalAmount.InexactFloat64(),
	)

	company, err := uc.companyRepo.GetCompanyOnlyByID(ctx, req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	paymentEntity := entity.NewEntity()

	// 3. Create MP Preference linked to PaymentID (ahead of time)
	mpReq := &mercadopagoservice.CheckoutRequest{
		CompanyID:         company.ID.String(),
		Schema:            company.SchemaName,
		Item:              checkoutItem,
		ExternalReference: paymentEntity.ID.String(), // Link to Payment
	}

	pref, err := uc.mpService.CreateCheckoutPreference(ctx, mpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create preference: %w", err)
	}

	// 4. Create Pending Payment
	payment := &companyentity.CompanyPayment{
		Entity:            paymentEntity,
		CompanyID:         company.ID,
		Provider:          mercadoPagoProvider,
		Status:            companyentity.PaymentStatusPending,
		Currency:          "BRL",
		Amount:            finalAmount,
		Months:            months,
		ProviderPaymentID: paymentEntity.ID.String(),
		PaymentURL:        pref.InitPoint,
		// PaidAt is nil
		ExternalReference: paymentEntity.ID.String(), // Self-reference or empty? Using ID as ref.
	}

	paymentModel := &model.CompanyPayment{}
	paymentModel.FromDomain(payment)

	if err := uc.companyPaymentRepo.CreateCompanyPayment(ctx, paymentModel); err != nil {
		return nil, fmt.Errorf("failed to create pending payment: %w", err)
	}

	return &billingdto.CheckoutResponseDTO{
		PaymentID:   payment.ID.String(),
		CheckoutUrl: pref.InitPoint,
	}, nil
}

func (s *CheckoutUseCase) HandleMercadoPagoWebhook(ctx context.Context, dto *companydto.MercadoPagoWebhookDTO) error {
	if s.mpService == nil || !s.mpService.Enabled() {
		return ErrMercadoPagoDisabled
	}

	if dto == nil || dto.Type != "payment" || dto.Data.ID == "" {
		return nil
	}

	// Get data.id for signature validation (prefer query param, fallback to body)
	dataIDForSignature := dto.DataIDFromQuery
	if dataIDForSignature == "" {
		dataIDForSignature = dto.Data.ID
	}

	// Validate signature from Mercado Pago
	if !s.mpService.ValidateSignature(dto.XSignature, dto.XRequestID, dataIDForSignature) {
		return ErrInvalidWebhookSecret
	}

	mpPaymentID := dto.Data.ID
	details, err := s.mpService.GetPayment(ctx, mpPaymentID)
	if err != nil {
		return err
	}

	if details.Status != "approved" {
		return nil
	}

	// 1. Find the pending Payment by External Reference (which is our PaymentID)
	// ExternalReference comes from MP, should be UUID
	paymentID, err := uuid.Parse(details.ExternalReference)
	if err != nil {
		return fmt.Errorf("invalid external reference (payment_id): %w", err)
	}

	paymentModel, err := s.companyPaymentRepo.GetCompanyPaymentByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("payment %s not found: %w", paymentID, err)
	}

	// Idempotency: if already paid, ignore
	if paymentModel.Status == string(companyentity.PaymentStatusApproved) || paymentModel.Status == string(companyentity.PaymentStatusPaid) {
		return nil
	}

	companyModel, err := s.companyRepo.GetCompanyOnlyByID(ctx, paymentModel.CompanyID)
	if err != nil {
		return err
	}

	paidAt := time.Now().UTC()
	if details.DateApproved != nil {
		paidAt = details.DateApproved.UTC()
	}

	// Update Subscription ONLY if it's a subscription payment
	if paymentModel.Months > 0 {
		base := paidAt
		if companyModel.SubscriptionExpiresAt != nil && companyModel.SubscriptionExpiresAt.After(paidAt) {
			base = *companyModel.SubscriptionExpiresAt
		}

		newExpiration := base.AddDate(0, paymentModel.Months, 0)

		if err := s.companyRepo.UpdateCompanySubscription(ctx, companyModel.ID, companyModel.SchemaName, &newExpiration, false); err != nil {
			return err
		}
	}

	// Update Payment
	rawPayload, _ := json.Marshal(dto)
	paymentModel.Status = details.Status
	paymentModel.ProviderPaymentID = details.ID
	paymentModel.PaidAt = &paidAt
	paymentModel.RawPayload = rawPayload

	// Assuming UpdateCompanyPayment exists (we'll implement it next)
	if err := s.companyPaymentRepo.UpdateCompanyPayment(ctx, paymentModel); err != nil {
		return err
	}

	// 2. Settle associated costs
	// Since we don't have GetByPaymentID on repository interface yet (or I missed adding it to the interface file logic above,
	// actually I added GetByPaymentID to usageCostRepo interface in previous step),
	// I should fetch costs by PaymentID and mark them as PAID.
	costs, err := s.costRepo.GetByPaymentID(ctx, paymentID)
	if err != nil {
		// Log error but don't fail the webhook processing as payment is already confirmed?
		// Better to return error so MP retries?
		// If we succeed in UpdateCompanyPayment but fail here, we might have inconsistency.
		// Idempotency check above handles retry of payment update.
		// We should try to update costs.
		return fmt.Errorf("failed to get associated costs: %w", err)
	}

	for _, cost := range costs {
		cost.Status = "PAID"
		// cost.PaidAt = &paidAt // usage cost model doesn't have PaidAt? it has Status.
		if err := s.costRepo.Update(ctx, cost); err != nil {
			return fmt.Errorf("failed to update cost status: %w", err)
		}
	}

	return nil
}

func (uc *CheckoutUseCase) CreateCostCheckout(ctx context.Context, companyID uuid.UUID) (*billingdto.CheckoutResponseDTO, error) {
	// 1. Fetch pending costs
	pendingCosts, err := uc.costRepo.GetPendingCosts(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pending costs: %w", err)
	}

	if len(pendingCosts) == 0 {
		return nil, errors.New("no pending costs found")
	}

	// 2. Calculate total amount
	totalAmount := decimal.Zero
	costIDs := make([]uuid.UUID, len(pendingCosts))
	for i, cost := range pendingCosts {
		totalAmount = totalAmount.Add(cost.Amount)
		costIDs[i] = cost.ID
	}

	// 3. Create Checkout Item and Preference first to get URL
	checkoutItem := mercadopagoservice.NewCheckoutItem(
		"Fatura de Custos Extras",
		fmt.Sprintf("Pagamento de %d custos pendentes (NFC-e, etc)", len(pendingCosts)),
		1,
		totalAmount.InexactFloat64(),
	)

	company, err := uc.companyRepo.GetCompanyOnlyByID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	paymentEntity := entity.NewEntity()

	mpReq := &mercadopagoservice.CheckoutRequest{
		CompanyID:         company.ID.String(),
		Schema:            company.SchemaName,
		Item:              checkoutItem,
		ExternalReference: paymentEntity.ID.String(),
	}

	pref, err := uc.mpService.CreateCheckoutPreference(ctx, mpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create preference: %w", err)
	}

	// 4. Create Pending Payment
	payment := &companyentity.CompanyPayment{
		Entity:            paymentEntity,
		CompanyID:         company.ID,
		Provider:          mercadoPagoProvider,
		Status:            companyentity.PaymentStatusPending,
		Currency:          "BRL",
		Amount:            totalAmount,
		Months:            0, // Not a subscription renewal
		ProviderPaymentID: paymentEntity.ID.String(),
		PaymentURL:        pref.InitPoint,
		ExternalReference: paymentEntity.ID.String(),
	}

	paymentModel := &model.CompanyPayment{}
	paymentModel.FromDomain(payment)

	if err := uc.companyPaymentRepo.CreateCompanyPayment(ctx, paymentModel); err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// 5. Link costs to Payment
	if err := uc.costRepo.UpdateCostsPaymentID(ctx, costIDs, payment.ID); err != nil {
		return nil, fmt.Errorf("failed to link costs to payment: %w", err)
	}

	return &billingdto.CheckoutResponseDTO{
		PaymentID:   payment.ID.String(),
		CheckoutUrl: pref.InitPoint,
	}, nil
}

func getEnvFloat(key string, fallback float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return fallback
}

func translatePeriodicity(p domainbilling.Periodicity) string {
	switch p {
	case domainbilling.PeriodicityMonthly:
		return "Mensal"
	case domainbilling.PeriodicitySemiannual:
		return "Semestral"
	case domainbilling.PeriodicityAnnual:
		return "Anual"
	default:
		return string(p)
	}
}

func translatePlanType(p domainbilling.PlanType) string {
	switch p {
	case domainbilling.PlanBasic:
		return "Básico"
	case domainbilling.PlanIntermediate:
		return "Intermediário"
	case domainbilling.PlanEnterprise:
		return "Enterprise"
	default:
		return string(p)
	}
}

func (uc *CheckoutUseCase) CancelPayment(ctx context.Context, paymentID uuid.UUID) error {
	payment, err := uc.companyPaymentRepo.GetCompanyPaymentByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	if payment.Status != string(companyentity.PaymentStatusPending) {
		return fmt.Errorf("payment cannot be cancelled (status: %s)", payment.Status)
	}

	// Unlink costs if any
	if err := uc.costRepo.UnlinkCostsFromPayment(ctx, paymentID); err != nil {
		return fmt.Errorf("failed to unlink costs: %w", err)
	}

	payment.Status = string(companyentity.PaymentStatusCancelled)
	// Ensure UpdateCompanyPayment is available and works as expected
	if err := uc.companyPaymentRepo.UpdateCompanyPayment(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}
