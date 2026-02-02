package billing

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	mercadopagoservice "github.com/willjrcom/sales-backend-go/internal/infra/service/mercadopago"
)

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
	// if !s.mpService.ValidateSignature(dto.XSignature, dto.XRequestID, dataIDForSignature) {
	// 	return ErrInvalidWebhookSecret
	// }

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
		return s.runCostWebhook(ctx, details, dto)
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

	activeSub, _, _ := s.companySubscriptionRepo.GetActiveAndUpcomingSubscriptions(ctx, company.ID)
	var subscriptionExpiresAt *time.Time
	if activeSub != nil {
		subscriptionExpiresAt = &activeSub.EndDate
	}

	if subscriptionExpiresAt != nil && subscriptionExpiresAt.After(paidAt) {
		newExpire := subscriptionExpiresAt.AddDate(0, months, 0)
		endDate = newExpire
	}

	// Removed UpdateCompanySubscription call as fields are moved to Subscription entity

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
		if err := s.companySubscriptionRepo.CreateSubscription(ctx, subModel); err != nil {
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
	if err := s.companySubscriptionRepo.CreateSubscription(ctx, subModel); err != nil {
		return err
	}

	return nil
}

// runCostWebhook handles payments for extra costs (NFC-e, etc)
func (s *CheckoutUseCase) runCostWebhook(ctx context.Context, details *mercadopagoservice.PaymentDetails, dto *companydto.MercadoPagoWebhookDTO) error {
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
	// Update active subscription plan
	activeSub, _, err := s.companySubscriptionRepo.GetActiveAndUpcomingSubscriptions(ctx, company.ID)
	if err == nil && activeSub != nil {
		activeSub.PlanType = companyentity.PlanType(targetPlan)
		if err := s.companySubscriptionRepo.UpdateSubscription(ctx, activeSub); err != nil {
			return fmt.Errorf("failed to update subscription plan: %w", err)
		}
	} else {
		return fmt.Errorf("no active subscription to upgrade")
	}

	fmt.Printf("Company %s plan upgraded successfully to %s\n", company.ID, targetPlan)
	return nil
}
