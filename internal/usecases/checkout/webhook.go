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

	// Get data.id for signature validation (prefer query param, fallback to body)
	dataIDForSignature := dto.DataIDFromQuery
	if dataIDForSignature == "" {
		dataIDForSignature = dto.Data.ID
	}

	// Validate signature from Mercado Pago
	// if !s.mpService.ValidateSignature(dto.XSignature, dto.XRequestID, dataIDForSignature) {
	// 	return ErrInvalidWebhookSecret
	// }

	if dto.Type == companydto.MercadoPagoWebhookTypeSubscriptionPreapproval {
		return s.runSubscriptionPreapprovalWebhook(ctx, dto)
	}

	if dto.Type != companydto.MercadoPagoWebhookTypePayment {
		return fmt.Errorf("unknown webhook type: %s", dto.Type)
	}

	fmt.Printf("Payment webhook received: %s, type: %s\n", dto.Data.ID, dto.Type)

	mpPaymentID := dto.Data.ID
	details, err := s.mpService.GetPayment(ctx, mpPaymentID)
	if err != nil {
		return err
	}

	// if details.Status != "approved" {
	// 	return nil
	// }

	parts := strings.Split(details.ExternalReference, ":")
	if len(parts) == 0 {
		return fmt.Errorf("external reference not found")
	}

	prefix := parts[0]

	switch prefix {
	case "SUB_UP":
		return s.runSubscriptionUpgradeWebhook(ctx, details, dto)
	case "COST":
		return s.runCostWebhook(ctx, details, dto)
	case "SUB":
		return s.runSubscriptionPaymentWebhook(ctx, details, dto)
	default:
		return fmt.Errorf("unknown payment type: %s", prefix)
	}
}

func (s *CheckoutUseCase) runSubscriptionPreapprovalWebhook(ctx context.Context, dto *companydto.MercadoPagoWebhookDTO) error {
	preapprovalID := dto.Data.ID
	preapproval, err := s.mpService.GetPreapproval(ctx, preapprovalID)
	if err != nil {
		return fmt.Errorf("failed to get preapproval %s: %w", preapprovalID, err)
	}

	fmt.Printf("Subscription preapproval webhook received: %s status: %s\n", preapprovalID, preapproval.Status)

	subscriptionExternalRef, err := mercadopagoservice.ExtractSubscriptionExternalRef(preapproval.ExternalReference)
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

	if preapproval.Status == "cancelled" {
		fmt.Printf("Subscription cancelled for company %s\n", subscriptionExternalRef.CompanyID)
		if err := s.companySubscriptionRepo.MarkSubscriptionAsCancelled(ctx, uuid.MustParse(subscriptionExternalRef.CompanyID)); err != nil {
			return fmt.Errorf("failed to cancel subscription: %w", err)
		}

		fmt.Printf("Subscription cancelled for company %s\n", subscriptionExternalRef.CompanyID)
		return nil
	}

	if preapproval.Status == "authorized" {
		fmt.Printf("Subscription authorized for company %s\n", subscriptionExternalRef.CompanyID)
		if err := s.companySubscriptionRepo.MarkSubscriptionAsActive(ctx, uuid.MustParse(subscriptionExternalRef.CompanyID)); err != nil {
			return fmt.Errorf("failed to activate subscription: %w", err)
		}
		fmt.Printf("Subscription authorized for company %s\n", subscriptionExternalRef.CompanyID)
		return nil
	}

	// Unknown status
	if err := s.companySubscriptionRepo.UpdateSubscriptionStatus(ctx, uuid.MustParse(subscriptionExternalRef.CompanyID), preapproval.Status); err != nil {
		fmt.Printf("Failed to update subscription status: %s\n", err.Error())
		return err
	}

	fmt.Printf("Subscription status updated for company %s\n", subscriptionExternalRef.CompanyID)
	return nil
}

// runSubscriptionPaymentWebhook handles the actual payment transaction for a subscription
func (s *CheckoutUseCase) runSubscriptionPaymentWebhook(ctx context.Context, details *mercadopagoservice.PaymentDetails, dto *companydto.MercadoPagoWebhookDTO) error {
	// Extract details from external reference
	subscriptionExternalRef, err := mercadopagoservice.ExtractSubscriptionExternalRef(details.ExternalReference)
	if err != nil {
		return fmt.Errorf("invalid external reference (subscription): %w", err)
	}

	// 1. Upsert Payment Record (History)
	transactionID := strconv.Itoa(details.ID)
	existingPayment, err := s.companyPaymentRepo.GetCompanyPaymentByExternalReferenceAndProviderID(ctx, details.ExternalReference, transactionID)

	paidAt := time.Now().UTC()
	if !details.DateApproved.IsZero() {
		paidAt = details.DateApproved.UTC()
	}
	rawPayload, _ := json.Marshal(dto)

	if existingPayment != nil {
		// Update existing payment
		fmt.Printf("Updating existing subscription payment: %s status: %s\n", existingPayment.ID, details.Status)
		existingPayment.Status = details.Status
		existingPayment.PaidAt = &paidAt
		existingPayment.RawPayload = rawPayload
		// PreapprovalID removed from CompanyPayment

		if err := s.companyPaymentRepo.UpdateCompanyPayment(ctx, existingPayment); err != nil {
			return fmt.Errorf("failed to update subscription payment: %w", err)
		}
	} else {
		// Create new payment (Recurrence or missed initial webhook)
		fmt.Printf("Creating new subscription payment (Recurrence): %s status: %s\n", transactionID, details.Status)
		amount := decimal.NewFromFloat(details.TransactionAmount)
		paymentEntity := entity.NewEntity()

		newPayment := &companyentity.CompanyPayment{
			Entity:            paymentEntity,
			CompanyID:         uuid.MustParse(subscriptionExternalRef.CompanyID),
			Provider:          mercadoPagoProvider,
			Status:            companyentity.PaymentStatus(details.Status),
			Currency:          details.CurrencyID,
			Amount:            amount,
			Months:            subscriptionExternalRef.Frequency,
			PaymentURL:        "", // No init point for recurrence
			ExternalReference: details.ExternalReference,
			ProviderPaymentID: &transactionID,
			IsMandatory:       false,
			PlanType:          companyentity.PlanType(subscriptionExternalRef.PlanType),
		}

		if details.Status == "approved" {
			newPayment.PaidAt = &paidAt
		}
		newPayment.RawPayload = rawPayload

		paymentModel := &model.CompanyPayment{}
		paymentModel.FromDomain(newPayment)

		if err := s.companyPaymentRepo.CreateCompanyPayment(ctx, paymentModel); err != nil {
			return fmt.Errorf("failed to create subscription payment: %w", err)
		}
	}

	// 3. Manage Subscription Contract (Activation/Renewal)
	// Only proceed if payment is approved
	if details.Status == "approved" {
		// Fetch the single subscription record for this contract
		subModel, err := s.companySubscriptionRepo.GetByExternalReference(ctx, details.ExternalReference)
		if err != nil {
			return fmt.Errorf("failed to get subscription: %w", err)
		}

		// Subscription found: Activate or Extend
		fmt.Printf("Updating active subscription %s for external reference %s\n", subModel.ID, details.ExternalReference)

		// If inactive (First payment logic), set start date
		if !subModel.IsActive {
			subModel.IsActive = true
			subModel.StartDate = paidAt
			subModel.EndDate = paidAt.AddDate(0, subscriptionExternalRef.Frequency, 0)
		} else {
			// Renewal: Extend end date
			// Simple logic: If expired, reset start/end from now. If valid, extend.
			if subModel.EndDate.Before(paidAt) {
				subModel.EndDate = paidAt.AddDate(0, subscriptionExternalRef.Frequency, 0)
			} else {
				subModel.EndDate = subModel.EndDate.AddDate(0, subscriptionExternalRef.Frequency, 0)
			}
		}

		if err := s.companySubscriptionRepo.UpdateSubscription(ctx, subModel); err != nil {
			return fmt.Errorf("failed to update subscription contract: %w", err)
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
	// Update payment status
	var paidAt *time.Time
	if details.Status == "approved" {
		t := time.Now().UTC()
		if !details.DateApproved.IsZero() {
			t = details.DateApproved.UTC()
		}
		paidAt = &t
	}

	transactionID := strconv.Itoa(details.ID)
	rawPayload, _ := json.Marshal(dto)
	paymentModel.Status = details.Status
	paymentModel.ProviderPaymentID = &transactionID
	paymentModel.PaidAt = paidAt
	paymentModel.RawPayload = rawPayload

	if err := s.companyPaymentRepo.UpdateCompanyPayment(ctx, paymentModel); err != nil {
		return err
	}

	// Mark linked costs as PAID only if approved
	if details.Status == "approved" {
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
	}

	// Unblock company only if approved
	if details.Status == "approved" {
		// Verify if there are other overdue payments before unblocking?
		payments, err := s.companyPaymentRepo.ListOverduePaymentsByCompany(ctx, paymentModel.CompanyID, time.Now().UTC())
		if err != nil {
			return fmt.Errorf("failed to get associated payments: %w", err)
		}

		if len(payments) == 0 {
			if err := s.companyRepo.UpdateBlockStatus(ctx, paymentModel.CompanyID, false); err != nil {
				return fmt.Errorf("failed to update company: %w", err)
			}
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
	// Update payment status
	var paidAt *time.Time
	if details.Status == "approved" {
		t := time.Now().UTC()
		if !details.DateApproved.IsZero() {
			t = details.DateApproved.UTC()
		}
		paidAt = &t
	}

	transactionID := strconv.Itoa(details.ID)
	rawPayload, _ := json.Marshal(dto)
	paymentModel.Status = details.Status
	paymentModel.ProviderPaymentID = &transactionID
	paymentModel.PaidAt = paidAt
	paymentModel.RawPayload = rawPayload

	if err := s.companyPaymentRepo.UpdateCompanyPayment(ctx, paymentModel); err != nil {
		return err
	}

	// Update company plan only if approved
	if details.Status == "approved" {
		// Get target plan from metadata
		targetPlan := details.Metadata.UpgradeTargetPlan
		if targetPlan == "" {
			return fmt.Errorf("upgrade_target_plan missing from metadata")
		}

		// Update active subscription plan
		activeSub, err := s.companySubscriptionRepo.GetActiveSubscription(ctx, paymentModel.CompanyID)
		if err != nil || activeSub == nil {
			return fmt.Errorf("active subscription not found for upgrade")
		}

		// Get PreapprovalID directly from active subscription
		if activeSub.PreapprovalID == nil {
			return fmt.Errorf("active subscription has no preapproval ID")
		}

		var frequencyMonth string
		switch subscriptionUpgradeExternalRef.Frequency {
		case 1:
			frequencyMonth = "MONTHLY"
		case 6:
			frequencyMonth = "SEMIANNUALLY"
		case 12:
			frequencyMonth = "ANNUALLY"
		}

		translatedPlanType := translatePlanType(companyentity.PlanType(subscriptionUpgradeExternalRef.PlanType))
		translatedFrequency := translateFrequency(companyentity.Frequency(frequencyMonth))
		title := fmt.Sprintf("Assinatura Gfood Plano %s - %s", translatedPlanType, translatedFrequency)

		// Update Mercado Pago Subscription Amount with FULL new plan price
		if err := s.mpService.UpdateSubscriptionAmount(ctx, *activeSub.PreapprovalID, title, subscriptionUpgradeExternalRef.NewAmount); err != nil {
			return fmt.Errorf("failed to update subscription amount: %w", err)
		}

		activeSub.PlanType = targetPlan
		if err := s.companySubscriptionRepo.UpdateSubscription(ctx, activeSub); err != nil {
			return fmt.Errorf("failed to update subscription plan: %w", err)
		}

		fmt.Printf("Company %s plan upgraded successfully to %s\n", paymentModel.CompanyID, targetPlan)
	}
	return nil
}
