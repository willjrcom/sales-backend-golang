package billing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

func (uc *CheckoutUseCase) CreateSubscriptionCheckout(ctx context.Context, req *billingdto.CreateSubscriptionCheckoutDTO) (*billingdto.CheckoutResponseDTO, error) {
	companyModel, err := uc.companyRepo.GetCompany(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	if companyModel.Email == "" {
		return nil, fmt.Errorf("company email is required for subscription")
	}

	// Check if user has active subscription
	hasActiveSubscription := companyModel.SubscriptionExpiresAt != nil && companyModel.SubscriptionExpiresAt.After(time.Now().UTC())

	if hasActiveSubscription {
		// User has active subscription - will create UPCOMING/SCHEDULED subscription
		// Do NOT cancel current subscription
		fmt.Printf("Creating scheduled subscription for company %s, will start after %s\n",
			companyModel.ID.String(), companyModel.SubscriptionExpiresAt.Format("2006-01-02"))
	} else {
		// No active subscription or expired - cancel any lingering preapproval and create new
		if err := uc.CancelSubscription(ctx, companyModel.ID); err != nil {
			// Log but continue - CancelSubscription handles "not found" cases
			fmt.Printf("Warning: failed to cancel previous subscription during new checkout: %v\n", err)
		}
	}

	var basePrice decimal.Decimal
	switch req.ToPlanType() {
	case companyentity.PlanIntermediate:
		price := getEnvFloat("PRICE_INTERMEDIATE", 119.90)
		basePrice = decimal.NewFromFloat(price)
	case companyentity.PlanAdvanced:
		price := getEnvFloat("PRICE_ADVANCED", 129.90)
		basePrice = decimal.NewFromFloat(price)
	default:
		price := getEnvFloat("PRICE_BASIC", 99.90)
		basePrice = decimal.NewFromFloat(price)
	}

	months := 1
	discount := 0.0
	frequency := 1
	frequencyType := "months"
	switch req.ToPeriodicity() {
	case companyentity.PeriodicitySemiannual:
		months = 6
		discount = 0.05
		frequency = 6
	case companyentity.PeriodicityAnnual:
		months = 12
		discount = 0.10
		frequency = 12
	default:
		frequency = 1
	}

	// DEBUG: Force daily frequency for testing if env set
	if os.Getenv("DEBUG_FAST_SUBSCRIPTION") == "true" {
		frequency = 1
		frequencyType = "days"
	}

	totalAmount := basePrice.Mul(decimal.NewFromInt(int64(months)))
	finalAmount := totalAmount.Mul(decimal.NewFromFloat(1.0 - discount))

	title := fmt.Sprintf("Assinatura Gfood Plano %s - %s", translatePlanType(req.ToPlanType()), translatePeriodicity(req.ToPeriodicity()))
	description := fmt.Sprintf("Cobrança recorrente a cada %d meses. Valor: %s", frequency, finalAmount.StringFixed(2))

	externalRef := fmt.Sprintf("SUB:%s:%s:%d", companyModel.ID.String(), req.ToPlanType(), months)

	paymentEntity := entity.NewEntity()
	var checkoutURL string
	var providerID string

	if hasActiveSubscription {
		// SCHEDULED SUBSCRIPTION: Create one-time payment (no Preapproval yet)
		// The subscription will be created as "upcoming" when payment is approved
		description = fmt.Sprintf("Assinatura Agendada - %s. Iniciará em %s",
			title, companyModel.SubscriptionExpiresAt.AddDate(0, 0, 1).Format("02/01/2006"))

		checkoutItem := mercadopagoservice.NewCheckoutItem(
			title,
			description,
			1,
			finalAmount.InexactFloat64(),
		)

		mpReq := &mercadopagoservice.CheckoutRequest{
			CompanyID:         companyModel.ID.String(),
			Schema:            companyModel.SchemaName,
			PaymentType:       mercadopagoservice.PaymentCheckoutTypeSubscription,
			Item:              checkoutItem,
			ExternalReference: paymentEntity.ID.String(),
			Metadata: map[string]any{
				"plan_type": string(req.ToPlanType()),
				"months":    months,
			},
		}

		pref, err := uc.mpService.CreateCheckoutPreference(ctx, mpReq)
		if err != nil {
			return nil, fmt.Errorf("failed to create scheduled subscription checkout: %w", err)
		}

		checkoutURL = pref.InitPoint
		providerID = paymentEntity.ID.String()
	} else {
		// IMMEDIATE SUBSCRIPTION: Create Mercado Pago Preapproval (recurring)
		subReq := &mercadopagoservice.SubscriptionRequest{
			Title:         title,
			Description:   description,
			Price:         finalAmount.InexactFloat64(),
			Frequency:     frequency,
			FrequencyType: frequencyType,
			ExternalRef:   externalRef,
			PayerEmail:    companyModel.Email,
			BackURL:       uc.mpService.SuccessURL(),
		}

		subResp, err := uc.mpService.CreateSubscription(ctx, subReq)
		if err != nil {
			return nil, fmt.Errorf("failed to create subscription: %w", err)
		}

		checkoutURL = subResp.InitPoint
		providerID = subResp.ID
	}

	payment := &companyentity.CompanyPayment{
		Entity:            paymentEntity,
		CompanyID:         companyModel.ID,
		Provider:          mercadoPagoProvider,
		ProviderPaymentID: providerID,
		PreapprovalID: func() string {
			if hasActiveSubscription {
				return ""
			} else {
				return providerID
			}
		}(),
		Status:            companyentity.PaymentStatusPending,
		Currency:          "BRL",
		Amount:            finalAmount,
		Months:            months,
		PlanType:          companyentity.PlanType(req.ToPlanType()),
		PaymentURL:        checkoutURL,
		ExternalReference: externalRef, // Always use SUB: format so frontend detects as "Assinatura"
		ExpiresAt:         func() *time.Time { t := time.Now().UTC().AddDate(0, 0, 5); return &t }(),
		IsMandatory:       false,
	}

	paymentModel := &model.CompanyPayment{}
	paymentModel.FromDomain(payment)

	if err := uc.companyPaymentRepo.CreateCompanyPayment(ctx, paymentModel); err != nil {
		return nil, fmt.Errorf("failed to create pending subscription payment: %w", err)
	}

	return &billingdto.CheckoutResponseDTO{
		PaymentID:   payment.ID.String(),
		CheckoutUrl: checkoutURL,
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

	// add schema_name to context from metadata
	ctx = context.WithValue(ctx, model.Schema("schema"), details.Metadata.SchemaName)
	companyModel, err := s.companyRepo.GetCompany(ctx)
	if err != nil {
		return err
	}

	company := companyModel.ToDomain()

	// Route based on payment type
	paymentType := details.Metadata.PaymentType

	// Fallback: If payment type is missing, check ExternalReference for legacy/recurring format
	if paymentType == "" && strings.HasPrefix(details.ExternalReference, "SUB:") {
		paymentType = string(mercadopagoservice.PaymentCheckoutTypeSubscription)
	}

	switch paymentType {
	case string(mercadopagoservice.PaymentCheckoutTypeSubscription):
		return s.runSubscriptionWebhook(ctx, details, company, dto)
	case string(mercadopagoservice.PaymentCheckoutTypeCost):
		return s.runCostWebhook(ctx, details, company, dto)
	case string(mercadopagoservice.PaymentCheckoutTypeSubscriptionUpgrade):
		return s.runSubscriptionUpgradeWebhook(ctx, details, company, dto)
	default:
		return fmt.Errorf("unknown payment type: %s", details.Metadata.PaymentType)
	}
}

// runSubscriptionWebhook handles payments from Mercado Pago Preapprovals (recurrent subscriptions)

func (s *CheckoutUseCase) runSubscriptionWebhook(ctx context.Context, details *mercadopagoservice.PaymentDetails, company *companyentity.Company, dto *companydto.MercadoPagoWebhookDTO) error {
	companyID := details.Metadata.CompanyID
	planType := details.Metadata.PlanType
	months := details.Metadata.Months

	// Fallback: If metadata is missing (recurring payment), parse from ExternalReference
	// Format: SUB:<CompanyID>:<PlanType>:<Months>
	if planType == "" && strings.HasPrefix(details.ExternalReference, "SUB:") {
		parts := strings.Split(details.ExternalReference, ":")
		if len(parts) >= 4 {
			// parts[0] = "SUB"
			// parts[1] = CompanyID
			// parts[2] = PlanType
			// parts[3] = Months
			if companyID == "" {
				companyID = parts[1]
			}
			planType = parts[2]
			if m, err := strconv.Atoi(parts[3]); err == nil {
				months = m
			}
		}
	}

	existing, err := s.companyPaymentRepo.GetCompanyPaymentByProviderID(ctx, details.ID)
	if err == nil && existing != nil {
		return nil
	}

	paidAt := time.Now().UTC()
	if details.DateApproved != nil {
		paidAt = details.DateApproved.UTC()
	}

	amount := decimal.NewFromFloat(details.TransactionAmount)
	startDate := paidAt
	endDate := startDate.AddDate(0, months, 0)

	if company.SubscriptionExpiresAt != nil && company.SubscriptionExpiresAt.After(paidAt) {
		newExpire := company.SubscriptionExpiresAt.AddDate(0, months, 0)
		endDate = newExpire
	}

	if err := s.companyRepo.UpdateCompanySubscription(ctx, company.ID, company.SchemaName, &endDate, string(planType)); err != nil {
		return err
	}

	var paymentToSave *model.CompanyPayment
	pending, _ := s.companyPaymentRepo.GetPendingPaymentByExternalReference(ctx, details.ExternalReference)
	// First payment will exists
	if pending != nil {
		paymentToSave = pending

		rawPayload, _ := json.Marshal(dto)
		paymentToSave.Status = details.Status
		paymentToSave.ProviderPaymentID = details.ID
		paymentToSave.PaidAt = &paidAt
		paymentToSave.RawPayload = rawPayload

		if err := s.companyPaymentRepo.UpdateCompanyPayment(ctx, paymentToSave); err != nil {
			return err
		}

		sub := companyentity.NewCompanySubscription(company.ID, companyentity.PlanType(planType), startDate, endDate)
		sub.PaymentID = &paymentToSave.ID

		subModel := &model.CompanySubscription{}
		subModel.FromDomain(sub)
		if err := s.companyRepo.CreateSubscription(ctx, subModel); err != nil {
			return err
		}

		return nil
	}

	// Recurrency payments must create
	paymentEntity := entity.NewEntity()
	domPay := &companyentity.CompanyPayment{
		Entity:            paymentEntity,
		CompanyID:         uuid.MustParse(companyID),
		Provider:          mercadoPagoProvider,
		Status:            companyentity.PaymentStatusPending,
		Currency:          "BRL",
		Amount:            amount,
		Months:            months,
		PlanType:          companyentity.PlanType(planType),
		ProviderPaymentID: details.ID,
		ExternalReference: details.ExternalReference,
		IsMandatory:       false,
	}
	paymentToSave = &model.CompanyPayment{}
	paymentToSave.FromDomain(domPay)

	rawPayload, _ := json.Marshal(dto)
	paymentToSave.Status = details.Status
	paymentToSave.ProviderPaymentID = details.ID
	paymentToSave.PaidAt = &paidAt
	paymentToSave.RawPayload = rawPayload

	if err := s.companyPaymentRepo.CreateCompanyPayment(ctx, paymentToSave); err != nil {
		return err
	}

	sub := companyentity.NewCompanySubscription(company.ID, companyentity.PlanType(planType), startDate, endDate)
	sub.PaymentID = &paymentToSave.ID

	subModel := &model.CompanySubscription{}
	subModel.FromDomain(sub)
	if err := s.companyRepo.CreateSubscription(ctx, subModel); err != nil {
		return err
	}

	return nil
}

// runCostWebhook handles payments for extra costs (NFC-e, etc)
func (s *CheckoutUseCase) runCostWebhook(ctx context.Context, details *mercadopagoservice.PaymentDetails, company *companyentity.Company, dto *companydto.MercadoPagoWebhookDTO) error {
	// Parse payment ID from external reference
	paymentID, err := uuid.Parse(details.ExternalReference)
	if err != nil {
		return fmt.Errorf("invalid external reference (payment_id): %w", err)
	}

	// Get payment
	paymentModel, err := s.companyPaymentRepo.GetCompanyPaymentByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("payment %s not found: %w", paymentID, err)
	}

	// Idempotency check
	if paymentModel.Status == string(companyentity.PaymentStatusApproved) || paymentModel.Status == string(companyentity.PaymentStatusPaid) {
		return nil
	}

	// Update payment status
	paidAt := time.Now().UTC()
	if details.DateApproved != nil {
		paidAt = details.DateApproved.UTC()
	}

	rawPayload, _ := json.Marshal(dto)
	paymentModel.Status = details.Status
	paymentModel.ProviderPaymentID = details.ID
	paymentModel.PaidAt = &paidAt
	paymentModel.RawPayload = rawPayload

	if err := s.companyPaymentRepo.UpdateCompanyPayment(ctx, paymentModel); err != nil {
		return err
	}

	// Mark linked costs as PAID
	costs, err := s.costRepo.GetByPaymentID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get associated costs: %w", err)
	}

	for _, cost := range costs {
		cost.Status = "PAID"
		if err := s.costRepo.Update(ctx, cost); err != nil {
			return fmt.Errorf("failed to update cost status: %w", err)
		}
	}

	return nil
}

// runSubscriptionUpgradeWebhook handles subscription upgrade payments
func (s *CheckoutUseCase) runSubscriptionUpgradeWebhook(ctx context.Context, details *mercadopagoservice.PaymentDetails, company *companyentity.Company, dto *companydto.MercadoPagoWebhookDTO) error {
	// Parse payment ID from external reference
	paymentID, err := uuid.Parse(details.ExternalReference)
	if err != nil {
		return fmt.Errorf("invalid external reference (payment_id): %w", err)
	}

	// Get payment
	paymentModel, err := s.companyPaymentRepo.GetCompanyPaymentByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("payment %s not found: %w", paymentID, err)
	}

	// Idempotency check
	if paymentModel.Status == string(companyentity.PaymentStatusApproved) || paymentModel.Status == string(companyentity.PaymentStatusPaid) {
		return nil
	}

	// Update payment status
	paidAt := time.Now().UTC()
	if details.DateApproved != nil {
		paidAt = details.DateApproved.UTC()
	}

	rawPayload, _ := json.Marshal(dto)
	paymentModel.Status = details.Status
	paymentModel.ProviderPaymentID = details.ID
	paymentModel.PaidAt = &paidAt
	paymentModel.RawPayload = rawPayload

	if err := s.companyPaymentRepo.UpdateCompanyPayment(ctx, paymentModel); err != nil {
		return err
	}

	// Get target plan from metadata
	targetPlan := details.Metadata.UpgradeTargetPlan
	if targetPlan == "" {
		return fmt.Errorf("upgrade_target_plan missing from metadata")
	}

	// Update company plan (keeps same expiration date - proration)
	if err := s.companyRepo.UpdateCompanySubscription(ctx, company.ID, company.SchemaName, company.SubscriptionExpiresAt, targetPlan); err != nil {
		return fmt.Errorf("failed to upgrade company plan: %w", err)
	}

	fmt.Printf("Company %s plan upgraded successfully to %s\n", company.ID, targetPlan)
	return nil
}

func (uc *CheckoutUseCase) GenerateMonthlyCostPayment(ctx context.Context, companyID uuid.UUID) error {
	// 1. Fetch pending costs
	pendingCosts, err := uc.costRepo.GetPendingCosts(ctx, companyID)
	if err != nil {
		return fmt.Errorf("failed to fetch pending costs: %w", err)
	}

	if len(pendingCosts) == 0 {
		return nil // Nothing to pay
	}

	// 2. Calculate total amount
	totalAmount := decimal.Zero
	costIDs := make([]uuid.UUID, len(pendingCosts))
	for i, cost := range pendingCosts {
		totalAmount = totalAmount.Add(cost.Amount)
		costIDs[i] = cost.ID
	}

	companyModel, err := uc.companyRepo.GetCompany(ctx)
	if err != nil {
		return err
	}

	// Determine expiration date
	now := time.Now().UTC()
	dueDay := companyModel.MonthlyPaymentDueDay
	if dueDay == 0 {
		dueDay = getEnvInt("MONTHLY_PAYMENT_DUE_DAY", 10)
	}
	expiresAt := time.Date(now.Year(), now.Month(), dueDay, 23, 59, 59, 0, now.Location())

	// If generated after the due day (e.g. rerun), set for next month or few days later?
	// Assuming running on day 1, due day 10 is fine.
	if now.Day() > dueDay {
		// Fallback if running late?
		expiresAt = now.AddDate(0, 0, 5) // 5 days from now
	}

	description := fmt.Sprintf("Fatura Mensal - %s/%d", translateMonth(now.Month()), now.Year())

	// 3. Create Checkout Item
	checkoutItem := mercadopagoservice.NewCheckoutItem(
		description,
		"Pagamento de custos mensais acumulados",
		1,
		totalAmount.InexactFloat64(),
	)

	paymentEntity := entity.NewEntity()

	mpReq := &mercadopagoservice.CheckoutRequest{
		CompanyID:         companyModel.ID.String(),
		Schema:            companyModel.SchemaName,
		PaymentType:       mercadopagoservice.PaymentCheckoutTypeCost,
		Item:              checkoutItem,
		ExternalReference: paymentEntity.ID.String(),
	}

	pref, err := uc.mpService.CreateCheckoutPreference(ctx, mpReq)
	if err != nil {
		return fmt.Errorf("failed to create preference: %w", err)
	}

	// 4. Create Mandatory Payment
	payment := &companyentity.CompanyPayment{
		Entity:            paymentEntity,
		CompanyID:         companyModel.ID,
		Provider:          mercadoPagoProvider,
		Status:            companyentity.PaymentStatusPending,
		Currency:          "BRL",
		Amount:            totalAmount,
		Months:            0,
		ProviderPaymentID: paymentEntity.ID.String(),
		PaymentURL:        pref.InitPoint,
		ExternalReference: paymentEntity.ID.String(),
		ExpiresAt:         &expiresAt,
		IsMandatory:       true,
	}

	paymentModel := &model.CompanyPayment{}
	paymentModel.FromDomain(payment)
	// Add Description manually as it's not in Domain yet?
	// Wait, I updated Domain to have ExpiresAt and IsMandatory.
	// I didn't add Description to Domain/Model?
	// Let's check repository model again.
	// Repository model has Description. Domain CompanyPayment DOES NOT have Description in my previous read?
	// Let's check 'internal/domain/company/payment.go'.
	// It doesn't have Description.
	// I should probably add Description to Domain too or just rely on MP Title.
	// For now, I'll skip Description mapping if it's missing in Domain.

	if err := uc.companyPaymentRepo.CreateCompanyPayment(ctx, paymentModel); err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	// 5. Link costs to Payment & Set Status
	if err := uc.costRepo.UpdateCostsPaymentID(ctx, costIDs, payment.ID); err != nil {
		return fmt.Errorf("failed to link costs to payment: %w", err)
	}

	return nil
}

func translateMonth(m time.Month) string {
	months := []string{"", "Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho", "Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro"}
	if m >= 1 && m <= 12 {
		return months[m]
	}
	return m.String()
}

func getEnvFloat(key string, fallback float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func translatePeriodicity(p companyentity.Periodicity) string {
	switch p {
	case companyentity.PeriodicityMonthly:
		return "Mensal"
	case companyentity.PeriodicitySemiannual:
		return "Semestral"
	case companyentity.PeriodicityAnnual:
		return "Anual"
	default:
		return string(p)
	}
}

func translatePlanType(p companyentity.PlanType) string {
	switch p {
	case companyentity.PlanBasic:
		return "Básico"
	case companyentity.PlanIntermediate:
		return "Intermediário"
	case companyentity.PlanAdvanced:
		return "Avançado"
	default:
		return string(p)
	}
}

// CancelSubscription cancels the company's active subscription (Preapproval).
func (uc *CheckoutUseCase) CancelSubscription(ctx context.Context, companyID uuid.UUID) error {
	// Get company to find the active subscription
	companyModel, err := uc.companyRepo.GetCompany(ctx)
	if err != nil {
		return fmt.Errorf("failed to get company: %w", err)
	}

	// Find the pending subscription payment with Preapproval ID
	// External reference format: SUB:<CompanyID>:<PlanType>:<Months>
	externalRef := fmt.Sprintf("SUB:%s:", companyModel.ID.String())

	// Get the most recent payment (Active/Approved) with this company's subscription reference
	payment, err := uc.companyPaymentRepo.GetLastPaymentByExternalReferencePrefix(ctx, externalRef)
	if err != nil || payment == nil {
		return fmt.Errorf("no active subscription found")
	}

	// Use the dedicated PreapprovalID field
	if payment.PreapprovalID == "" {
		return fmt.Errorf("subscription has no preapproval ID")
	}

	// Cancel the subscription (Preapproval) at Mercado Pago - stops future renewals
	if err := uc.mpService.CancelSubscription(ctx, payment.PreapprovalID); err != nil {
		return fmt.Errorf("failed to cancel subscription at Mercado Pago: %w", err)
	}

	// Mark the active subscription as canceled directly (no fetch needed)
	if err := uc.companyRepo.MarkActiveSubscriptionAsCanceled(ctx, companyID); err != nil {
		return fmt.Errorf("failed to mark subscription as canceled: %w", err)
	}

	return nil
}

// CalculateUpgradeProration calculates the credit and new cost for upgrading a plan
func (uc *CheckoutUseCase) CalculateUpgradeProration(ctx context.Context, targetPlan companyentity.PlanType) (*billingdto.UpgradeSimulationDTO, error) {
	companyModel, err := uc.companyRepo.GetCompany(ctx)
	if err != nil {
		return nil, err
	}

	currentPlan := companyentity.PlanType(companyModel.CurrentPlan)
	if currentPlan == targetPlan {
		return nil, errors.New("target plan is same as current plan")
	}

	// Get current subscription to determine periodicity (months)
	subRefPrefix := fmt.Sprintf("SUB:%s:", companyModel.ID.String())
	currentSubscription, err := uc.companyPaymentRepo.GetLastPaymentByExternalReferencePrefix(ctx, subRefPrefix)
	months := 1 // Default to monthly if no subscription found
	if err == nil && currentSubscription != nil && currentSubscription.Months > 0 {
		months = currentSubscription.Months
	}

	// Prices with discount applied based on periodicity
	currentPrice := getPlanPriceWithDiscount(currentPlan, months)
	targetPrice := getPlanPriceWithDiscount(targetPlan, months)

	if targetPrice <= currentPrice {
		return nil, errors.New("target plan must be higher value (downgrade not supported via this flow)")
	}

	// Calculate remaining days
	daysRemaining := 0
	if companyModel.SubscriptionExpiresAt != nil && companyModel.SubscriptionExpiresAt.After(time.Now().UTC()) {
		daysRemaining = int(time.Until(*companyModel.SubscriptionExpiresAt).Hours() / 24)
	}
	isFullRenewal := false
	var upgradeAmount float64

	if daysRemaining < 1 || currentPrice == 0 {
		// Early Renewal Logic OR Upgrading from Free Plan
		// User pays full price of new plan, and subscription is extended by 1 full period
		daysRemaining = 0
		upgradeAmount = targetPrice
		isFullRenewal = true
	} else {
		// Standard Proration Logic
		diffPerMonth := targetPrice - currentPrice
		dailyDiff := diffPerMonth / 30.0
		upgradeAmount = dailyDiff * float64(daysRemaining)
		upgradeAmount = decimal.NewFromFloat(upgradeAmount).Round(2).InexactFloat64()
	}

	return &billingdto.UpgradeSimulationDTO{
		TargetPlan:     string(targetPlan),
		OldPlan:        string(currentPlan),
		DaysRemaining:  daysRemaining,
		UpgradeAmount:  upgradeAmount,
		NewMonthlyCost: targetPrice,
		IsFullRenewal:  isFullRenewal,
	}, nil
}

// CreateUpgradeCheckout generates a payment link for the upgrade difference
func (uc *CheckoutUseCase) CreateUpgradeCheckout(ctx context.Context, targetPlan companyentity.PlanType) (*billingdto.CheckoutResponseDTO, error) {
	companyModel, err := uc.companyRepo.GetCompany(ctx)
	if err != nil {
		return nil, err
	}

	sim, err := uc.CalculateUpgradeProration(ctx, targetPlan)
	if err != nil {
		return nil, err
	}

	amount := decimal.NewFromFloat(sim.UpgradeAmount)
	description := fmt.Sprintf("Upgrade para plano %s (%d dias restantes)", translatePlanType(targetPlan), sim.DaysRemaining)

	// Create Checkout Item
	checkoutItem := mercadopagoservice.NewCheckoutItem(
		description,
		"Upgrade de Plano",
		1,
		amount.InexactFloat64(),
	)

	paymentEntity := checkoutPaymentEntity()

	mpReq := &mercadopagoservice.CheckoutRequest{
		CompanyID:         companyModel.ID.String(),
		Schema:            companyModel.SchemaName,
		PaymentType:       mercadopagoservice.PaymentCheckoutTypeSubscriptionUpgrade,
		Item:              checkoutItem,
		ExternalReference: paymentEntity.ID.String(),
		Metadata: map[string]interface{}{
			"upgrade_target_plan": string(targetPlan),
			"is_full_renewal":     sim.IsFullRenewal,
			"company_id":          companyModel.ID.String(),
		},
	}

	pref, err := uc.mpService.CreateCheckoutPreference(ctx, mpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create preference: %w", err)
	}

	// Create Pending Payment
	payment := &companyentity.CompanyPayment{
		Entity:            paymentEntity,
		CompanyID:         companyModel.ID,
		Provider:          mercadoPagoProvider,
		Status:            companyentity.PaymentStatusPending,
		Currency:          "BRL",
		Amount:            amount,
		Months:            0,
		ProviderPaymentID: paymentEntity.ID.String(),
		PaymentURL:        pref.InitPoint,
		ExternalReference: paymentEntity.ID.String(), // Payment ID
		ExpiresAt:         func() *time.Time { t := time.Now().UTC().AddDate(0, 0, 2); return &t }(),
		IsMandatory:       false,
		Description:       description,
		PlanType:          companyentity.PlanType(targetPlan), // Track the target plan
	}

	paymentModel := &model.CompanyPayment{}
	paymentModel.FromDomain(payment)

	fmt.Printf("DEBUG: Creating upgrade payment for company %s, amount: %v, plan: %s\n", companyModel.ID, amount, targetPlan)

	if err := uc.companyPaymentRepo.CreateCompanyPayment(ctx, paymentModel); err != nil {
		fmt.Printf("ERROR: Failed to create payment: %v\n", err)
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	fmt.Printf("DEBUG: Payment created successfully with ID: %s\n", payment.ID)

	return &billingdto.CheckoutResponseDTO{
		PaymentID:   payment.ID.String(),
		CheckoutUrl: pref.InitPoint,
	}, nil
}

func getPlanPrice(p companyentity.PlanType) float64 {
	switch p {
	case companyentity.PlanIntermediate:
		return getEnvFloat("PRICE_INTERMEDIATE", 119.90)
	case companyentity.PlanAdvanced:
		return getEnvFloat("PRICE_ADVANCED", 129.90)
	case companyentity.PlanBasic:
		return getEnvFloat("PRICE_BASIC", 99.90)
	default:
		return 0
	}
}

// getPlanPriceWithDiscount returns the plan price with discount applied based on periodicity
func getPlanPriceWithDiscount(p companyentity.PlanType, months int) float64 {
	basePrice := getPlanPrice(p)
	if basePrice == 0 {
		return 0
	}

	// Apply discount based on periodicity
	var discountPercent float64
	if months >= 12 {
		// Annual discount
		discountPercent = getEnvFloat("DISCOUNT_ANNUAL_PERCENT", 10)
	} else if months >= 6 {
		// Semester discount
		discountPercent = getEnvFloat("DISCOUNT_SEMESTER_PERCENT", 5)
	}

	return basePrice * (1 - discountPercent/100)
}

func checkoutPaymentEntity() entity.Entity {
	return entity.NewEntity()
}
