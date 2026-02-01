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
	company, err := uc.companyRepo.GetCompanyOnlyByID(ctx, req.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	if company.Email == "" {
		return nil, fmt.Errorf("company email is required for subscription")
	}

	// Check if user has active subscription
	hasActiveSubscription := company.SubscriptionExpiresAt != nil && company.SubscriptionExpiresAt.After(time.Now().UTC())

	if hasActiveSubscription {
		// User has active subscription - will create UPCOMING/SCHEDULED subscription
		// Do NOT cancel current subscription
		fmt.Printf("Creating scheduled subscription for company %s, will start after %s\n",
			company.ID.String(), company.SubscriptionExpiresAt.Format("2006-01-02"))
	} else {
		// No active subscription or expired - cancel any lingering preapproval and create new
		if err := uc.CancelSubscription(ctx, req.CompanyID); err != nil {
			// Log but continue - CancelSubscription handles "not found" cases
			fmt.Printf("Warning: failed to cancel previous subscription during new checkout: %v\n", err)
		}
	}

	var basePrice decimal.Decimal
	switch req.ToPlanType() {
	case domainbilling.PlanIntermediate:
		price := getEnvFloat("PRICE_INTERMEDIATE", 119.90)
		basePrice = decimal.NewFromFloat(price)
	case domainbilling.PlanAdvanced:
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
	case domainbilling.PeriodicitySemiannual:
		months = 6
		discount = 0.05
		frequency = 6
	case domainbilling.PeriodicityAnnual:
		months = 12
		discount = 0.10
		frequency = 12
	default:
		frequency = 1
	}

	totalAmount := basePrice.Mul(decimal.NewFromInt(int64(months)))
	finalAmount := totalAmount.Mul(decimal.NewFromFloat(1.0 - discount))

	title := fmt.Sprintf("Assinatura Gfood Plano %s - %s", translatePlanType(req.ToPlanType()), translatePeriodicity(req.ToPeriodicity()))
	description := fmt.Sprintf("Cobrança recorrente a cada %d meses. Valor: %s", frequency, finalAmount.StringFixed(2))

	externalRef := fmt.Sprintf("SUB:%s:%s:%d", company.ID.String(), req.ToPlanType(), months)

	paymentEntity := entity.NewEntity()
	var checkoutURL string
	var providerID string

	if hasActiveSubscription {
		// SCHEDULED SUBSCRIPTION: Create one-time payment (no Preapproval yet)
		// The subscription will be created as "upcoming" when payment is approved
		description = fmt.Sprintf("Assinatura Agendada - %s. Iniciará em %s",
			title, company.SubscriptionExpiresAt.AddDate(0, 0, 1).Format("02/01/2006"))

		checkoutItem := mercadopagoservice.NewCheckoutItem(
			title,
			description,
			1,
			finalAmount.InexactFloat64(),
		)

		mpReq := &mercadopagoservice.CheckoutRequest{
			CompanyID:         company.ID.String(),
			Schema:            company.SchemaName,
			Item:              checkoutItem,
			ExternalReference: paymentEntity.ID.String(),
			Metadata: map[string]any{
				"plan_type":   string(req.ToPlanType()),
				"months":      months,
				"is_upcoming": true, // Flag to indicate this is a scheduled subscription
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
			PayerEmail:    company.Email,
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
		CompanyID:         company.ID,
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

	if strings.HasPrefix(details.ExternalReference, "SUB:") {
		parts := strings.Split(details.ExternalReference, ":")
		if len(parts) < 4 {
			return fmt.Errorf("invalid subscription reference format: %s", details.ExternalReference)
		}

		companyID, err := uuid.Parse(parts[1])
		if err != nil {
			return fmt.Errorf("invalid company ID in subscription reference: %w", err)
		}

		planType := companyentity.PlanType(parts[2])
		months, err := strconv.Atoi(parts[3])
		if err != nil || months <= 0 {
			return fmt.Errorf("invalid months in subscription reference: %s", parts[3])
		}

		existing, err := s.companyPaymentRepo.GetCompanyPaymentByProviderID(ctx, details.ID)
		if err == nil && existing != nil {
			return nil
		}

		companyModel, err := s.companyRepo.GetCompanyOnlyByID(ctx, companyID)
		if err != nil {
			return err
		}

		paidAt := time.Now().UTC()
		if details.DateApproved != nil {
			paidAt = details.DateApproved.UTC()
		}

		amount := decimal.NewFromFloat(details.TransactionAmount)
		startDate := paidAt
		endDate := startDate.AddDate(0, months, 0)

		if companyModel.SubscriptionExpiresAt != nil && companyModel.SubscriptionExpiresAt.After(paidAt) {
			newExpire := companyModel.SubscriptionExpiresAt.AddDate(0, months, 0)
			endDate = newExpire
		}

		if err := s.companyRepo.UpdateCompanySubscription(ctx, companyModel.ID, companyModel.SchemaName, &endDate, string(planType)); err != nil {
			return err
		}

		pending, _ := s.companyPaymentRepo.GetPendingPaymentByExternalReference(ctx, details.ExternalReference)

		var paymentToSave *model.CompanyPayment

		if pending != nil {
			paymentToSave = pending
		} else {
			paymentEntity := entity.NewEntity()
			domPay := &companyentity.CompanyPayment{
				Entity:            paymentEntity,
				CompanyID:         companyID,
				Provider:          mercadoPagoProvider,
				Status:            companyentity.PaymentStatusPending,
				Currency:          "BRL",
				Amount:            amount,
				Months:            months,
				PlanType:          planType,
				ProviderPaymentID: details.ID,
				ExternalReference: details.ExternalReference,
				IsMandatory:       false,
			}
			paymentToSave = &model.CompanyPayment{}
			paymentToSave.FromDomain(domPay)
		}

		rawPayload, _ := json.Marshal(dto)
		paymentToSave.Status = details.Status
		paymentToSave.ProviderPaymentID = details.ID
		paymentToSave.PaidAt = &paidAt
		paymentToSave.RawPayload = rawPayload

		if pending != nil {
			if err := s.companyPaymentRepo.UpdateCompanyPayment(ctx, paymentToSave); err != nil {
				return err
			}
		} else {
			if err := s.companyPaymentRepo.CreateCompanyPayment(ctx, paymentToSave); err != nil {
				return err
			}
		}

		sub := companyentity.NewCompanySubscription(companyModel.ID, planType, startDate, endDate)
		sub.PaymentID = &paymentToSave.ID

		subModel := &model.CompanySubscription{}
		subModel.FromDomain(sub)
		if err := s.companyRepo.CreateSubscription(ctx, subModel); err != nil {
			return err
		}

		return nil
	}

	// 1. Find the pending Payment by External Reference (Legacy/One-time Flow)
	// ExternalReference comes from MP, should be UUID
	paymentID, err := uuid.Parse(details.ExternalReference)
	if err != nil {
		// If it's not a UUID and didn't match SUB, it's unknown/invalid
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
	// We check metadata first (new flow) or fallback to CompanyPayment.Months (legacy/redundant)
	planType := details.Metadata.PlanType
	months := details.Metadata.Months
	isUpcoming := details.Metadata.IsUpcoming

	if planType != "" && months > 0 {
		var startDate, endDate time.Time

		if isUpcoming {
			// SCHEDULED SUBSCRIPTION: Start after current subscription expires
			if companyModel.SubscriptionExpiresAt == nil {
				return fmt.Errorf("cannot create scheduled subscription without active subscription")
			}
			startDate = companyModel.SubscriptionExpiresAt.AddDate(0, 0, 1) // Start day after expiration
			endDate = startDate.AddDate(0, months, 0)
		} else {
			// IMMEDIATE SUBSCRIPTION: Extend from current expiration or start now
			base := paidAt
			if companyModel.SubscriptionExpiresAt != nil && companyModel.SubscriptionExpiresAt.After(paidAt) {
				base = *companyModel.SubscriptionExpiresAt
			}
			startDate = base
			endDate = base.AddDate(0, months, 0)
		}

		// Create Subscription Record
		sub := companyentity.NewCompanySubscription(companyModel.ID, companyentity.PlanType(planType), startDate, endDate)
		// Link Payment
		sub.PaymentID = &paymentModel.ID

		subModel := &model.CompanySubscription{}
		subModel.FromDomain(sub)
		if err := s.companyRepo.CreateSubscription(ctx, subModel); err != nil {
			return err
		}

		// Update Company Current Plan Snapshot ONLY if not upcoming
		if !isUpcoming {
			if err := s.companyRepo.UpdateCompanySubscription(ctx, companyModel.ID, companyModel.SchemaName, &endDate, planType); err != nil {
				return err
			}
		}
		// If isUpcoming, we DON'T update company's current plan - it stays as-is until the scheduled time
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

	// 3. Process Upgrade Costs
	costs, err := s.costRepo.GetByPaymentID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get associated costs: %w", err)
	}

	for _, cost := range costs {
		cost.Status = "PAID"
		if err := s.costRepo.Update(ctx, cost); err != nil {
			return fmt.Errorf("failed to update cost status: %w", err)
		}

		// Handle Plan Upgrade Logic
		if isUpgradeCost(cost.CostType) {
			targetPlan := getPlanFromCostType(cost.CostType)
			if targetPlan != "" {
				// 1. Update Company Plan locally
				// Only update the PlanType, keep expiration date (it's proration for same cycle)
				if err := s.companyRepo.UpdateCompanySubscription(ctx, companyModel.ID, companyModel.SchemaName, companyModel.SubscriptionExpiresAt, string(targetPlan)); err != nil {
					return fmt.Errorf("failed to update company plan: %w", err)
				}

				// 2. Update Subscription Amount in Mercado Pago (for NEXT renewal)
				// We need to find the active subscription (Preapproval ID)
				// It's not in THIS payment (this is the upgrade payment), it's in the previous subscription payment
				// We can look up by company pending payment or just search for active subscription logic?
				// Better: We can check if companyModel has a "PreapprovalID" field? No.
				// We search for the *active subscription payment* for this company.

				// External reference format: SUB:<CompanyID>:<PlanType>:<Months>
				// But PlanType changes... so we search by prefix "SUB:<CompanyID>:"
				subRefPrefix := fmt.Sprintf("SUB:%s:", companyModel.ID.String())
				subPayment, err := s.companyPaymentRepo.GetLastPaymentByExternalReferencePrefix(ctx, subRefPrefix)

				if err == nil && subPayment != nil && subPayment.PreapprovalID != "" {
					newPrice := getPlanPrice(targetPlan)

					// 2a. Determine if it's a Full Renewal (Metadata check)
					// The metadata comes from the PAYMENT (upgrade payment), which is `details`.
					// IsFullRenewal is in details.Metadata.IsFullRenewal (if structure allows) or we fetch it from cost?
					// Cost doesn't have metadata easily accessible here without fetching.
					// But we passed metadata to Mercado Pago Payment! `details.Metadata`

					// Assuming details.Metadata is a map or struct that has our custom fields.
					// The SDK usually puts custom fields in Metadata.

					// We need to implement UpdateSubscription (generic) or UpdateSubscriptionDate
					// Let's use UpdateSubscriptionAmount for now, and if FullRenewal, we need to update date too.

					// Check "is_full_renewal" in metadata
					isFullRenewal := details.Metadata.IsFullRenewal

					if isFullRenewal {
						// It's a renewal!
						// 1. Calculate new NextPaymentDate
						// Current expiration seems to be NOW (because < 1 day).
						// Or better: Current Expiration + Period.
						// We updated companyModel locally in previous step?
						// "if isUpgradeCost ... UpdateCompanySubscription ... "
						// The local update logic needs to know about renewal too to update date locally!
						// Wait, previous logic was:
						// "if err := s.companyRepo.UpdateCompanySubscription(..., companyModel.SubscriptionExpiresAt, ...)"
						// It KEPT the old date. This is WRONG for renewal.

						// FIX: We need to know if it's renewal BEFORE updating local subscription.
						// But access to metadata is easy here.

						// Let's fix the local update first.

						// Re-read company to get fresh dates? No, companyModel is good.
						// If IsFullRenewal, add period.
						var newLocalExpire *time.Time
						if isFullRenewal {
							months := 1
							if subPayment != nil && subPayment.Months > 0 {
								months = subPayment.Months
							}

							// Add period to current expiration (or Now if expired?)
							// If < 1 day, it is essentially Now or tomorrow.
							// Safer: Add period to time.Now() or keep drift?
							// Logic in CreateCheckout was "until SubscriptionExpiresAt".
							// So we should respect SubscriptionExpiresAt + Period.
							if companyModel.SubscriptionExpiresAt != nil {
								t := companyModel.SubscriptionExpiresAt.AddDate(0, months, 0)
								newLocalExpire = &t
								// We don't have periodicity here easily. We know PlanType.
								// We assume same periodicity? Or we need to fetch plan details?
								// "Upgrade" usually implies keeping periodicity or changing it?
								// The current flow assumes keeping periodicity implicitly?
								// The upgrade simulation didn't ask for periodicity.
								// So we assume same periodicity as current.

								// We need to know current periodicity.
								// company.Periodicity?
								// Let's assume Monthly for now or fetch it.
								// Company model has Periodicity?
								// Let's look at company model...
								// If we can't find it, we fallback to Monthly (1 month).

								// subPayment logic merged above
							}
						} else {
							newLocalExpire = companyModel.SubscriptionExpiresAt
						}

						if err := s.companyRepo.UpdateCompanySubscription(ctx, companyModel.ID, companyModel.SchemaName, newLocalExpire, string(targetPlan)); err != nil {
							return fmt.Errorf("failed to update company plan: %w", err)
						}

						if isFullRenewal && newLocalExpire != nil {
							// Update MP Date
							// We need a method for this: UpdateSubscriptionDate
							// We'll implement UpdateSubscriptionAmountAndDate or generic.
							// For now, let's call UpdateSubscriptionAmount (we need to update amount anyway for next cycle).

							// And ALSO update date.
							if err := s.mpService.UpdateSubscriptionDate(ctx, subPayment.PreapprovalID, *newLocalExpire); err != nil {
								// Log error
								fmt.Printf("Error updating subscription date: %v\n", err)
							}
						}
					} else {
						// Normal Proration Upgrade
						if err := s.companyRepo.UpdateCompanySubscription(ctx, companyModel.ID, companyModel.SchemaName, companyModel.SubscriptionExpiresAt, string(targetPlan)); err != nil {
							return fmt.Errorf("failed to update company plan: %w", err)
						}
					}

					// Always Update Amount for next cycle
					if err := s.mpService.UpdateSubscriptionAmount(ctx, subPayment.PreapprovalID, newPrice); err != nil {
						return fmt.Errorf("failed to update MP subscription amount: %w", err)
					}

					// Also update the subscription payment external reference
					subPayment.PlanType = string(targetPlan)
					parts := strings.Split(subPayment.ExternalReference, ":")
					if len(parts) >= 4 {
						newRef := fmt.Sprintf("SUB:%s:%s:%s", parts[1], targetPlan, parts[3])
						subPayment.ExternalReference = newRef
						_ = s.companyPaymentRepo.UpdateCompanyPayment(ctx, subPayment)
					}
				}
			}
		}
	}

	// 4. Attempt to unblock company if there are no more OVERDUE (> 5 days) mandatory payments
	if companyModel.IsBlocked {
		// Use same cutoff logic as Scheduler: payments overdue by more than 5 days
		cutoffDate := time.Now().UTC().AddDate(0, 0, -5)
		overdue, err := s.companyPaymentRepo.ListOverduePaymentsByCompany(ctx, companyModel.ID, cutoffDate)
		if err == nil && len(overdue) == 0 {
			_ = s.companyRepo.UpdateBlockStatus(ctx, companyModel.ID, false)
		}
	}

	return nil
}

func isUpgradeCost(costType string) bool {
	return costType == "upgrade_basic" ||
		costType == "upgrade_intermediate" ||
		costType == "upgrade_advanced"
}

func getPlanFromCostType(costType string) domainbilling.PlanType {
	switch costType {
	case "upgrade_intermediate":
		return domainbilling.PlanIntermediate
	case "upgrade_advanced":
		return domainbilling.PlanAdvanced
	case "upgrade_basic":
		return domainbilling.PlanBasic
	default:
		return ""
	}
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
		ExpiresAt:         func() *time.Time { t := time.Now().UTC().AddDate(0, 0, 5); return &t }(),
		IsMandatory:       false,
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

	company, err := uc.companyRepo.GetCompanyOnlyByID(ctx, companyID)
	if err != nil {
		return err
	}

	// Determine expiration date
	now := time.Now().UTC()
	dueDay := company.MonthlyPaymentDueDay
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
		CompanyID:         company.ID.String(),
		Schema:            company.SchemaName,
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
		CompanyID:         company.ID,
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
	case domainbilling.PlanAdvanced:
		return "Avançado"
	default:
		return string(p)
	}
}

// CancelSubscription cancels the company's active subscription (Preapproval).
func (uc *CheckoutUseCase) CancelSubscription(ctx context.Context, companyID uuid.UUID) error {
	// Get company to find the active subscription
	company, err := uc.companyRepo.GetCompanyOnlyByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("failed to get company: %w", err)
	}

	// Find the pending subscription payment with Preapproval ID
	// External reference format: SUB:<CompanyID>:<PlanType>:<Months>
	externalRef := fmt.Sprintf("SUB:%s:", company.ID.String())

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

// CalculateUpgradeProration calculates the credit and new cost for upgrading a plan
func (uc *CheckoutUseCase) CalculateUpgradeProration(ctx context.Context, companyID uuid.UUID, targetPlan domainbilling.PlanType) (*billingdto.UpgradeSimulationDTO, error) {
	company, err := uc.companyRepo.GetCompanyOnlyByID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	currentPlan := domainbilling.PlanType(company.CurrentPlan)
	if currentPlan == targetPlan {
		return nil, errors.New("target plan is same as current plan")
	}

	// Get current subscription to determine periodicity (months)
	subRefPrefix := fmt.Sprintf("SUB:%s:", companyID.String())
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
	if company.SubscriptionExpiresAt != nil && company.SubscriptionExpiresAt.After(time.Now().UTC()) {
		daysRemaining = int(time.Until(*company.SubscriptionExpiresAt).Hours() / 24)
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
func (uc *CheckoutUseCase) CreateUpgradeCheckout(ctx context.Context, companyID uuid.UUID, targetPlan domainbilling.PlanType) (*billingdto.CheckoutResponseDTO, error) {
	sim, err := uc.CalculateUpgradeProration(ctx, companyID, targetPlan)
	if err != nil {
		return nil, err
	}

	company, err := uc.companyRepo.GetCompanyOnlyByID(ctx, companyID)
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
		CompanyID:         company.ID.String(),
		Schema:            company.SchemaName,
		Item:              checkoutItem,
		ExternalReference: paymentEntity.ID.String(),
		Metadata: map[string]interface{}{
			"upgrade_target_plan": string(targetPlan),
			"is_full_renewal":     sim.IsFullRenewal,
		},
	}

	pref, err := uc.mpService.CreateCheckoutPreference(ctx, mpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create preference: %w", err)
	}

	// Create Pending Payment
	payment := &companyentity.CompanyPayment{
		Entity:            paymentEntity,
		CompanyID:         company.ID,
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

	fmt.Printf("DEBUG: Creating upgrade payment for company %s, amount: %v, plan: %s\n", company.ID, amount, targetPlan)

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

func getPlanPrice(p domainbilling.PlanType) float64 {
	switch p {
	case domainbilling.PlanIntermediate:
		return getEnvFloat("PRICE_INTERMEDIATE", 119.90)
	case domainbilling.PlanAdvanced:
		return getEnvFloat("PRICE_ADVANCED", 129.90)
	case domainbilling.PlanBasic:
		return getEnvFloat("PRICE_BASIC", 99.90)
	default:
		return 0
	}
}

// getPlanPriceWithDiscount returns the plan price with discount applied based on periodicity
func getPlanPriceWithDiscount(p domainbilling.PlanType, months int) float64 {
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
