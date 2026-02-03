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

	if details.Status != "approved" {
		return nil
	}

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

	if preapproval.Status != "authorized" {
		fmt.Println("Preapproval not authorized")
		return nil
	}

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

	// Check removed: We allow multiple payments for the same preapproval (recurrence)
	// and we don't store PreapprovalID in ProviderPaymentID anymore.

	// Handle cancellation
	if preapproval.Status == "cancelled" {
		fmt.Printf("Subscription cancelled for company %s\n", subscriptionExternalRef.CompanyID)
		if err := s.companySubscriptionRepo.MarkActiveSubscriptionAsCanceled(ctx, uuid.MustParse(subscriptionExternalRef.CompanyID)); err != nil {
			return fmt.Errorf("failed to cancel subscription: %w", err)
		}
		return nil
	}

	// Use DateCreated for subscriptions as approval date isn't always present in the same format
	paidAt := preapproval.DateCreated.UTC()

	amount := decimal.NewFromFloat(preapproval.AutoRecurring.TransactionAmount)
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

	existingPayments, _ := s.companyPaymentRepo.GetCompanyPaymentsByExternalReference(ctx, subscriptionExternalRef.ExternalReference)
	if len(existingPayments) > 0 {
		fmt.Printf("Updating existing subscription payment (First Payment) %s\n", existingPayments[0].ID)
		// We only update the PreapprovalID linkage if missing, but we DON'T touch payment status/dates here.
		// The Payment Webhook logic (runSubscriptionPaymentWebhook) is the owner of Status, PaidAt and ProviderPaymentID.
		if existingPayments[0].PreapprovalID == nil || *existingPayments[0].PreapprovalID == "" {
			existingPayment.PreapprovalID = &preapproval.ID
			_ = s.companyPaymentRepo.UpdateCompanyPayment(ctx, existingPayment)
		}

		sub := companyentity.NewCompanySubscription(uuid.MustParse(subscriptionExternalRef.CompanyID), companyentity.PlanType(subscriptionExternalRef.PlanType), startDate, endDate)
		sub.PaymentID = &existingPayment.ID

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
		Status:            companyentity.PaymentStatus(preapproval.Status),
		Currency:          "BRL",
		Amount:            amount,
		Months:            subscriptionExternalRef.Frequency,
		PlanType:          companyentity.PlanType(subscriptionExternalRef.PlanType),
		ProviderPaymentID: nil,             // No payment ID yet. Will be updated by payment webhook
		PreapprovalID:     &preapproval.ID, // Link to the same preapproval
		ExternalReference: preapproval.ExternalReference,
		IsMandatory:       false,
		PaidAt:            &paidAt,
	}
	paymentToSave := &model.CompanyPayment{}
	paymentToSave.FromDomain(domPay)

	rawPayload, _ := json.Marshal(dto)
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

// runSubscriptionPaymentWebhook handles the actual payment transaction for a subscription
func (s *CheckoutUseCase) runSubscriptionPaymentWebhook(ctx context.Context, details *mercadopagoservice.PaymentDetails, dto *companydto.MercadoPagoWebhookDTO) error {
	// Extract details from external reference
	subscriptionExternalRef, err := mercadopagoservice.ExtractSubscriptionExternalRef(details.ExternalReference)
	if err != nil {
		return fmt.Errorf("invalid external reference (subscription): %w", err)
	}

	if details.Order.ID == "" {
		return fmt.Errorf("order id not found")
	}

	// se ja existir, é porque o pagamento ja foi processado
	paymentModel, err := s.companyPaymentRepo.GetCompanyPaymentByExternalReferenceAndProviderID(ctx, details.ExternalReference, strconv.Itoa(details.ID))
	if paymentModel != nil {
		return fmt.Errorf("payment already processed")
	}

	// 2. Update with Payment Transaction details
	paidAt := time.Now().UTC()
	if !details.DateApproved.IsZero() {
		paidAt = details.DateApproved.UTC()
	}

	transactionID := strconv.Itoa(details.ID)
	rawPayload, _ := json.Marshal(dto)
	paymentModel = &model.CompanyPayment{
		Status:            details.Status,
		ProviderPaymentID: &transactionID,
		PaidAt:            &paidAt,
		RawPayload:        rawPayload,
	}

	// 3. Save updates
	if err := s.companyPaymentRepo.UpdateCompanyPayment(ctx, paymentModel); err != nil {
		return fmt.Errorf("failed to update subscription payment: %w", err)
	}

	fmt.Printf("Subscription payment processed: %s (Transaction ID: %s)\n", subscriptionExternalRef.CompanyID, transactionID)
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
	if !details.DateApproved.IsZero() {
		paidAt = details.DateApproved.UTC()
	}

	transactionID := strconv.Itoa(details.ID)
	rawPayload, _ := json.Marshal(dto)
	paymentModel.Status = details.Status
	paymentModel.ProviderPaymentID = &transactionID
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

	payments, err := s.companyPaymentRepo.ListOverduePaymentsByCompany(ctx, paymentModel.CompanyID, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("failed to get associated payments: %w", err)
	}

	if len(payments) > 0 {
		return nil
	}

	// Unblock company
	if err := s.companyRepo.UpdateBlockStatus(ctx, paymentModel.CompanyID, false); err != nil {
		return fmt.Errorf("failed to update company: %w", err)
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
	if !details.DateApproved.IsZero() {
		paidAt = details.DateApproved.UTC()
	}

	transactionID := strconv.Itoa(details.ID)
	rawPayload, _ := json.Marshal(dto)
	paymentModel.Status = details.Status
	paymentModel.ProviderPaymentID = &transactionID
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

	if subPayment.PreapprovalID == nil {
		return fmt.Errorf("subscription payment has no preapproval ID")
	}

	// Update Mercado Pago Subscription Amount with FULL new plan price
	if err := s.mpService.UpdateSubscriptionAmount(ctx, *subPayment.PreapprovalID, subscriptionUpgradeExternalRef.NewAmount); err != nil {
		return fmt.Errorf("failed to update subscription amount: %w", err)
	}

	activeSub.PlanType = companyentity.PlanType(targetPlan)
	if err := s.companySubscriptionRepo.UpdateSubscription(ctx, activeSub); err != nil {
		return fmt.Errorf("failed to update subscription plan: %w", err)
	}

	fmt.Printf("Company %s plan upgraded successfully to %s\n", paymentModel.CompanyID, targetPlan)
	return nil
}
