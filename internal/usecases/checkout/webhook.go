package billing

import (
	"context"
	"encoding/json"
	"fmt"
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

	// Route based on payment type
	paymentType := details.Metadata.PaymentType

	// Fallback: If payment type is missing, check ExternalReference for legacy/recurring format
	if paymentType == "" && strings.HasPrefix(details.ExternalReference, "SUB:") {
		paymentType = string(mercadopagoservice.PaymentCheckoutTypeSubscription)
	}

	switch mercadopagoservice.PaymentCheckoutType(paymentType) {
	case mercadopagoservice.PaymentCheckoutTypeSubscription:
		return s.runSubscriptionWebhook(ctx, details, dto)

	case mercadopagoservice.PaymentCheckoutTypeSubscriptionUpgrade:
		return s.runSubscriptionUpgradeWebhook(ctx, details, dto)

	case mercadopagoservice.PaymentCheckoutTypeCost:
		return s.runCostWebhook(ctx, details, dto)

	default:
		return fmt.Errorf("unknown payment type: %s", details.Metadata.PaymentType)
	}
}

// runSubscriptionWebhook handles payments from Mercado Pago Preapprovals (recurrent subscriptions)

func (s *CheckoutUseCase) runSubscriptionWebhook(ctx context.Context, details *mercadopagoservice.PaymentDetails, dto *companydto.MercadoPagoWebhookDTO) error {
	subscriptionExternalRef, err := mercadopagoservice.ExtractSubscriptionExternalRef(details.ExternalReference)
	if err != nil {
		return err
	}

	if subscriptionExternalRef.CompanyID == "" {
		return fmt.Errorf("company id not found")
	}

	if subscriptionExternalRef.PlanType == "" {
		return fmt.Errorf("plan type not found")
	}

	if subscriptionExternalRef.Frequency == 0 {
		return fmt.Errorf("frequency not found")
	}

	existing, err := s.companyPaymentRepo.GetCompanyPaymentByProviderID(ctx, details.ID)
	if err == nil && existing != nil {
		fmt.Println("Payment already exists")
		return nil
	}

	paidAt := time.Now().UTC()
	if details.DateApproved != nil {
		paidAt = details.DateApproved.UTC()
	}

	amount := decimal.NewFromFloat(details.TransactionAmount)
	startDate := paidAt

	// Check for active subscription to determine start date (renewal vs new)
	activeSub, _ := s.companySubscriptionRepo.GetActiveSubscription(ctx, uuid.MustParse(subscriptionExternalRef.CompanyID))
	if activeSub != nil && activeSub.PlanType != companyentity.PlanFree {
		if activeSub.EndDate.After(paidAt) {
			startDate = activeSub.EndDate
		}
	}

	// Calculate end date based on start date
	endDate := startDate.AddDate(0, subscriptionExternalRef.Frequency, 0)

	// Removed UpdateCompanySubscription call as fields are moved to Subscription entity

	pending, _ := s.companyPaymentRepo.GetPendingPaymentByExternalReference(ctx, details.ExternalReference)
	// First payment will exists
	if pending != nil {
		fmt.Printf("Creating first subscription by company id %s", subscriptionExternalRef.CompanyID)
		paymentToSave := pending

		rawPayload, _ := json.Marshal(dto)
		paymentToSave.Status = details.Status
		paymentToSave.ProviderPaymentID = details.ID
		paymentToSave.PaidAt = &paidAt
		paymentToSave.RawPayload = rawPayload

		if err := s.companyPaymentRepo.UpdateCompanyPayment(ctx, paymentToSave); err != nil {
			return err
		}

		sub := companyentity.NewCompanySubscription(uuid.MustParse(subscriptionExternalRef.CompanyID), companyentity.PlanType(subscriptionExternalRef.PlanType), startDate, endDate)
		sub.PaymentID = &paymentToSave.ID

		subModel := &model.CompanySubscription{}
		subModel.FromDomain(sub)
		if err := s.companySubscriptionRepo.CreateSubscription(ctx, subModel); err != nil {
			return err
		}

		// Update active free subscription plan type to finish the free trial
		if activeSub != nil && activeSub.PlanType == companyentity.PlanFree {
			activeSub.EndDate = startDate
			if err := s.companySubscriptionRepo.UpdateSubscription(ctx, activeSub); err != nil {
				return fmt.Errorf("failed to update current free subscription to new plan: %w", err)
			}
		}

		return nil
	}

	fmt.Printf("Creating recurrency subscription by company id %s", subscriptionExternalRef.CompanyID)

	// Recurrency payments must create
	paymentEntity := entity.NewEntity()
	domPay := &companyentity.CompanyPayment{
		Entity:            paymentEntity,
		CompanyID:         uuid.MustParse(subscriptionExternalRef.CompanyID),
		Provider:          mercadoPagoProvider,
		Status:            companyentity.PaymentStatusPending,
		Currency:          "BRL",
		Amount:            amount,
		Months:            subscriptionExternalRef.Frequency,
		PlanType:          companyentity.PlanType(subscriptionExternalRef.PlanType),
		ProviderPaymentID: details.ID,
		ExternalReference: details.ExternalReference,
		IsMandatory:       false,
	}
	paymentToSave := &model.CompanyPayment{}
	paymentToSave.FromDomain(domPay)

	rawPayload, _ := json.Marshal(dto)
	paymentToSave.Status = details.Status
	paymentToSave.ProviderPaymentID = details.ID
	paymentToSave.PaidAt = &paidAt
	paymentToSave.RawPayload = rawPayload

	if err := s.companyPaymentRepo.CreateCompanyPayment(ctx, paymentToSave); err != nil {
		return err
	}

	sub := companyentity.NewCompanySubscription(uuid.MustParse(subscriptionExternalRef.CompanyID), companyentity.PlanType(subscriptionExternalRef.PlanType), startDate, endDate)
	sub.PaymentID = &paymentToSave.ID

	subModel := &model.CompanySubscription{}
	subModel.FromDomain(sub)
	if err := s.companySubscriptionRepo.CreateSubscription(ctx, subModel); err != nil {
		return err
	}

	// Update active free subscription plan type to finish the free trial
	if activeSub != nil && activeSub.PlanType == companyentity.PlanFree {
		activeSub.EndDate = startDate
		if err := s.companySubscriptionRepo.UpdateSubscription(ctx, activeSub); err != nil {
			return fmt.Errorf("failed to update current free subscription to new plan: %w", err)
		}
	}
	return nil
}

// runCostWebhook handles payments for extra costs (NFC-e, etc)
func (s *CheckoutUseCase) runCostWebhook(ctx context.Context, details *mercadopagoservice.PaymentDetails, dto *companydto.MercadoPagoWebhookDTO) error {
	costExternalRef, err := mercadopagoservice.ExtractCostExternalRef(details.ExternalReference)
	if err != nil {
		return fmt.Errorf("invalid external reference (cost): %w", err)
	}

	paymentID, err := uuid.Parse(costExternalRef.PaymentID)
	if err != nil {
		return fmt.Errorf("invalid payment ID: %w", err)
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
func (s *CheckoutUseCase) runSubscriptionUpgradeWebhook(ctx context.Context, details *mercadopagoservice.PaymentDetails, dto *companydto.MercadoPagoWebhookDTO) error {
	// Parse payment ID from external reference
	subscriptionUpgradeExternalRef, err := mercadopagoservice.ExtractSubscriptionUpgradeExternalRef(details.ExternalReference)
	if err != nil {
		return fmt.Errorf("invalid external reference (subscription_upgrade): %w", err)
	}

	paymentID, err := uuid.Parse(subscriptionUpgradeExternalRef.PaymentID)
	if err != nil {
		return fmt.Errorf("invalid payment ID: %w", err)
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
	activeSub, err := s.companySubscriptionRepo.GetActiveSubscription(ctx, paymentModel.CompanyID)
	if err != nil || activeSub == nil {
		return fmt.Errorf("active subscription not found for upgrade")
	}

	// Get the payment linked to the active subscription to find the PreapprovalID
	if activeSub.PaymentID == nil {
		return fmt.Errorf("active subscription has no linked payment")
	}

	subPayment, err := s.companyPaymentRepo.GetCompanyPaymentByID(ctx, *activeSub.PaymentID)
	if err != nil {
		return fmt.Errorf("failed to get subscription payment: %w", err)
	}

	if subPayment.PreapprovalID == "" {
		return fmt.Errorf("subscription payment has no preapproval ID")
	}

	// Update Mercado Pago Subscription Amount
	if err := s.mpService.UpdateSubscriptionAmount(ctx, subPayment.PreapprovalID, subscriptionUpgradeExternalRef.NewAmount); err != nil {
		return fmt.Errorf("failed to update subscription amount: %w", err)
	}

	activeSub.PlanType = companyentity.PlanType(targetPlan)
	if err := s.companySubscriptionRepo.UpdateSubscription(ctx, activeSub); err != nil {
		return fmt.Errorf("failed to update subscription plan: %w", err)
	}

	fmt.Printf("Company %s plan upgraded successfully to %s\n", paymentModel.CompanyID, targetPlan)
	return nil
}
