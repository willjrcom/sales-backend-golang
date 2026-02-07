package billing

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	billingdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/checkout"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	mercadopagoservice "github.com/willjrcom/sales-backend-go/internal/infra/service/mercadopago"
)

type CheckoutUseCase struct {
	costRepo                model.CompanyUsageCostRepository
	companyRepo             model.CompanyRepository
	companyPaymentRepo      model.CompanyPaymentRepository
	companySubscriptionRepo model.CompanySubscriptionRepository
	mpService               *mercadopagoservice.Client
	userRepo                model.UserRepository
}

func NewCheckoutUseCase(
	costRepo model.CompanyUsageCostRepository,
	companyRepo model.CompanyRepository,
	companyPaymentRepo model.CompanyPaymentRepository,
	companySubscriptionRepo model.CompanySubscriptionRepository,
	mpService *mercadopagoservice.Client,
) *CheckoutUseCase {
	return &CheckoutUseCase{
		costRepo:                costRepo,
		companyRepo:             companyRepo,
		companyPaymentRepo:      companyPaymentRepo,
		companySubscriptionRepo: companySubscriptionRepo,
		mpService:               mpService,
	}
}

const mercadoPagoProvider = "mercado_pago"

var (
	ErrMercadoPagoDisabled  = errors.New("mercado pago integration disabled")
	ErrInvalidWebhookSecret = errors.New("invalid mercado pago webhook secret")
)

func (uc *CheckoutUseCase) AddDependencies(userRepo model.UserRepository) {
	uc.userRepo = userRepo
}

func (uc *CheckoutUseCase) CreateSubscriptionCheckout(ctx context.Context, req *billingdto.CreateSubscriptionCheckoutDTO) (*billingdto.CheckoutResponseDTO, error) {
	companyModel, err := uc.companyRepo.GetCompany(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	if companyModel.Email == "" {
		return nil, fmt.Errorf("company email is required for subscription")
	}

	// Fetch active/upcoming subscriptions
	activeSub, _ := uc.companySubscriptionRepo.GetActiveSubscription(ctx, companyModel.ID)
	if activeSub != nil && activeSub.PlanType != string(companyentity.PlanFree) {
		return nil, fmt.Errorf("company already has an active subscription")
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

	frequency := 1
	discount := 0.0
	frequencyType := "months"
	switch req.ToFrequency() {
	case companyentity.FrequencySemiannually:
		frequency = 6
		discount = 0.05
	case companyentity.FrequencyAnnually:
		frequency = 12
		discount = 0.10
	}

	// DEBUG: Force daily frequency for testing if env set
	if os.Getenv("DEBUG_FAST_SUBSCRIPTION") == "true" {
		frequency = 1
		frequencyType = "days"
	}

	// Calculate total amount per period (e.g. Total for 12 months)
	// BasePrice is monthly price.
	totalAmount := basePrice.Mul(decimal.NewFromInt(int64(frequency)))
	finalAmount := totalAmount.Mul(decimal.NewFromFloat(1.0 - discount)).Round(2)

	title := fmt.Sprintf("Assinatura Gfood Plano %s - %s", translatePlanType(req.ToPlanType()), translateFrequency(req.ToFrequency()))

	// Generate external reference
	externalRef := mercadopagoservice.NewSubscriptionExternalRef(companyModel.ID.String(), string(req.ToPlanType()), frequency)

	// DEBUG: Verify amount being sent
	fmt.Printf("DEBUG: Creating subscription. Frequency: %d, Total Amount: %s\n", frequency, finalAmount.String())

	userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)
	if !ok {
		return nil, errors.New("context user not found")
	}

	userIDUUID := uuid.MustParse(userID)
	userModel, err := uc.userRepo.GetUserByID(ctx, userIDUUID, false)
	if err != nil {
		return nil, err
	}

	// Create Mercado Pago Preapproval (recurring)
	subReq := &mercadopagoservice.SubscriptionRequest{
		Title:         title,
		Price:         finalAmount.InexactFloat64(),
		Frequency:     frequency,
		FrequencyType: frequencyType,
		ExternalRef:   externalRef,
		PayerEmail:    userModel.Email,
		BackURL:       uc.mpService.SuccessURL(),
	}

	// Create subscription
	subResp, err := uc.mpService.CreateSubscription(ctx, subReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// Create Inactive Subscription Record (Single record per contract)
	inactiveSub := companyentity.NewCompanySubscription(
		companyModel.ID,
		companyentity.PlanType(req.ToPlanType()),
		time.Now().UTC(), // StartDate placeholder
		time.Now().UTC(), // EndDate placeholder
	)
	inactiveSub.IsActive = false
	pid := subResp.ID
	inactiveSub.PreapprovalID = &pid
	inactiveSub.ExternalReference = &externalRef

	inactiveSubModel := &model.CompanySubscription{}
	inactiveSubModel.FromDomain(inactiveSub)

	if err := uc.companySubscriptionRepo.CreateSubscription(ctx, inactiveSubModel); err != nil {
		return nil, fmt.Errorf("failed to create inactive subscription: %w", err)
	}

	return &billingdto.CheckoutResponseDTO{
		CheckoutUrl: subResp.InitPoint,
	}, nil
}

func (uc *CheckoutUseCase) GenerateMonthlyCostPayment(ctx context.Context, companyID uuid.UUID) error {
	// 1. Fetch pending costs
	pendingCosts, err := uc.costRepo.GetPendingCosts(ctx, companyID)
	if err != nil {
		return fmt.Errorf("failed to fetch pending costs: %w", err)
	}

	if len(pendingCosts) == 0 {
		fmt.Println("No pending costs found")
		return nil // Nothing to pay
	}

	// 2. Calculate total amount
	totalAmount := decimal.Zero
	costIDs := make([]uuid.UUID, len(pendingCosts))
	for i, cost := range pendingCosts {
		totalAmount = totalAmount.Add(cost.Amount)
		costIDs[i] = cost.ID
	}

	totalAmount = totalAmount.Round(2)

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
		"COST-"+description, // ID
		"services",          // CategoryID
		description,
		"Pagamento de custos mensais acumulados",
		1,
		totalAmount.InexactFloat64(),
	)

	paymentEntity := entity.NewEntity()

	// Generate external reference
	// Format: COST:company_id:day:month:year:company_payment_id
	externalRef := mercadopagoservice.NewCostExternalRef(companyModel.ID.String(), paymentEntity.ID.String())

	mpReq := &mercadopagoservice.CheckoutRequest{
		CompanyID:         companyModel.ID.String(),
		Schema:            companyModel.SchemaName,
		PaymentType:       mercadopagoservice.PaymentCheckoutTypeCost,
		Item:              checkoutItem,
		Payer:             uc.getPayerFromCompany(companyModel), // Added Payer
		ExternalReference: externalRef,
	}

	pref, err := uc.mpService.CreateUniqueCheckout(ctx, mpReq)
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
		PaymentURL:        pref.InitPoint,
		ExternalReference: externalRef,
		ExpiresAt:         &expiresAt,
		IsMandatory:       true,
	}

	paymentModel := &model.CompanyPayment{}
	paymentModel.FromDomain(payment)

	if err := uc.companyPaymentRepo.CreateCompanyPayment(ctx, paymentModel); err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	// 5. Link costs to Payment & Set Status
	if err := uc.costRepo.UpdateCostsPaymentID(ctx, costIDs, payment.ID); err != nil {
		return fmt.Errorf("failed to link costs to payment: %w", err)
	}

	return nil
}

// CancelSubscription cancels the company's active subscription (Preapproval).
func (uc *CheckoutUseCase) CancelSubscription(ctx context.Context) error {
	// Get company to find the active subscription
	companyModel, err := uc.companyRepo.GetCompany(ctx)
	if err != nil {
		return fmt.Errorf("failed to get company: %w", err)
	}

	// Get active subscription
	activeSub, err := uc.companySubscriptionRepo.GetActiveSubscription(ctx, companyModel.ID)
	if err != nil {
		return fmt.Errorf("failed to get active subscription: %w", err)
	}
	if activeSub == nil {
		return fmt.Errorf("no active subscription found")
	}

	if activeSub.PreapprovalID == nil || *activeSub.PreapprovalID == "" {
		return fmt.Errorf("active subscription has no preapproval ID")
	}

	// Cancel the subscription (Preapproval) at Mercado Pago - stops future renewals
	if err := uc.mpService.CancelSubscription(ctx, *activeSub.PreapprovalID); err != nil {
		return fmt.Errorf("failed to cancel subscription at Mercado Pago: %w", err)
	}

	if activeSub.ExternalReference == nil || *activeSub.ExternalReference == "" {
		return fmt.Errorf("active subscription has no external reference")
	}

	// Mark the active subscription as cancelled directly
	if err := uc.companySubscriptionRepo.MarkSubscriptionAsCancelled(ctx, companyModel.ID, *activeSub.ExternalReference); err != nil {
		return fmt.Errorf("failed to mark subscription as cancelled: %w", err)
	}

	return nil
}

// CalculateUpgradeProration calculates the credit and new cost for upgrading a plan
func (uc *CheckoutUseCase) CalculateUpgradeProration(ctx context.Context, targetPlan companyentity.PlanType) (*billingdto.UpgradeSimulationDTO, error) {
	companyModel, err := uc.companyRepo.GetCompany(ctx)
	if err != nil {
		return nil, err
	}

	activeSub, _ := uc.companySubscriptionRepo.GetActiveSubscription(ctx, companyModel.ID)
	if activeSub == nil {
		return nil, errors.New("no active subscription found")
	}

	// Get current plan from active subscription
	currentPlan := activeSub.PlanType

	if currentPlan == string(targetPlan) {
		return nil, errors.New("target plan is same as current plan")
	}
	// Get current subscription to determine frequency (months)
	months := 1 // Default to monthly if no subscription found

	if activeSub.ExternalReference == nil {
		return nil, errors.New("active subscription has no external reference")
	}
	subExtRef, err := mercadopagoservice.ExtractSubscriptionExternalRef(*activeSub.ExternalReference)
	if err != nil {
		return nil, err
	}
	months = subExtRef.Frequency

	// Prices with discount applied based on frequency
	currentPrice := getPlanPriceWithDiscount(companyentity.PlanType(currentPlan), months)
	targetPrice := getPlanPriceWithDiscount(targetPlan, months)

	if targetPrice.LessThanOrEqual(currentPrice) {
		return nil, errors.New("target plan must be higher value (downgrade not supported via this flow)")
	}

	// Calculate remaining days
	daysRemaining := 0
	if activeSub.EndDate.After(time.Now().UTC()) {
		daysRemaining = int(time.Until(activeSub.EndDate).Hours() / 24)
	}

	// Default logic: calculate proration based on remaining days
	// targetPrice and currentPrice are now TOTAL amounts for the period (months)
	diffTotal := targetPrice.Sub(currentPrice)

	// Calculate daily difference: Total Difference / Total Days (approx 30 days per month)
	totalDays := int64(months * 30)
	dailyDiff := diffTotal.Div(decimal.NewFromInt(totalDays))

	upgradeAmount := dailyDiff.Mul(decimal.NewFromFloat(float64(daysRemaining))).InexactFloat64()

	// Early Renewal Logic OR Upgrading from Free Plan
	if daysRemaining < 1 || currentPrice.IsZero() {
		// User pays full price of new plan, and subscription is extended by 1 full period
		daysRemaining = 0
		upgradeAmount = targetPrice.InexactFloat64()
	}

	// Always round to 2 decimals for API compatibility
	upgradeAmount = decimal.NewFromFloat(upgradeAmount).Round(2).InexactFloat64()

	// For display simplicity, monthly cost is Total / Months
	monthlyDisplayCost := targetPrice.Div(decimal.NewFromInt(int64(months))).Round(2).InexactFloat64()

	return &billingdto.UpgradeSimulationDTO{
		TargetPlan:     string(targetPlan),
		OldPlan:        string(currentPlan),
		DaysRemaining:  daysRemaining,
		UpgradeAmount:  upgradeAmount,
		NewMonthlyCost: monthlyDisplayCost,
		Frequency:      months,
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
	// Create Checkout Item
	checkoutItem := mercadopagoservice.NewCheckoutItem(
		string(targetPlan), // ID
		"services",         // CategoryID
		description,
		"Upgrade de Plano",
		1,
		amount.InexactFloat64(),
	)

	paymentEntity := checkoutPaymentEntity()

	totalCost := decimal.NewFromFloat(sim.NewMonthlyCost).Mul(decimal.NewFromInt(int64(sim.Frequency))).Round(2).InexactFloat64()

	// Generate external reference
	externalRef := mercadopagoservice.NewSubscriptionUpgradeExternalRef(companyModel.ID.String(), sim.TargetPlan, sim.Frequency, totalCost, paymentEntity.ID.String())

	mpReq := &mercadopagoservice.CheckoutRequest{
		CompanyID:         companyModel.ID.String(),
		Schema:            companyModel.SchemaName,
		PaymentType:       mercadopagoservice.PaymentCheckoutTypeSubscriptionUpgrade,
		Item:              checkoutItem,
		Payer:             uc.getPayerFromCompany(companyModel), // Added Payer
		ExternalReference: externalRef,
		Metadata: map[string]interface{}{
			"upgrade_target_plan": string(targetPlan),
		},
	}

	pref, err := uc.mpService.CreateUniqueCheckout(ctx, mpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create preference: %w", err)
	}

	expiresAt := time.Now().UTC().AddDate(0, 0, 2)
	// Create Pending Payment
	payment := &companyentity.CompanyPayment{
		Entity:            paymentEntity,
		CompanyID:         companyModel.ID,
		Provider:          mercadoPagoProvider,
		Status:            companyentity.PaymentStatusPending,
		Currency:          "BRL",
		Amount:            amount,
		Months:            0,
		PaymentURL:        pref.InitPoint,
		ExternalReference: paymentEntity.ID.String(), // Payment ID
		ExpiresAt:         &expiresAt,
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

// getPlanPriceWithDiscount returns the plan price with discount applied based on frequency
func getPlanPriceWithDiscount(p companyentity.PlanType, months int) decimal.Decimal {
	basePrice := decimal.NewFromFloat(getPlanPrice(p))
	if basePrice.IsZero() {
		return decimal.Zero
	}

	// Apply discount based on frequency
	var discountPercent float64
	if months >= 12 {
		// Annual discount
		discountPercent = getEnvFloat("DISCOUNT_ANNUALLY_PERCENT", 10)
	} else if months >= 6 {
		// Semester discount
		discountPercent = getEnvFloat("DISCOUNT_SEMIANNUALLY_PERCENT", 5)
	}

	discountedMonthly := basePrice.Mul(decimal.NewFromFloat(1 - discountPercent/100))
	return discountedMonthly.Mul(decimal.NewFromInt(int64(months))).Round(2)
}

func (uc *CheckoutUseCase) getPayerFromCompany(company *model.Company) *mercadopagoservice.CheckoutPayer {
	payer := &mercadopagoservice.CheckoutPayer{
		Email: company.Email,
		Name:  company.BusinessName,
	}

	// Try to get phone from contacts
	if len(company.Contacts) > 0 {
		// Simple parser for phone: +5511999999999 or 11999999999
		// Provide basic AreaCode/Number split if possible, otherwise send as number
		phone := company.Contacts[0]
		// Remove non-numeric characters
		numericPhone := strings.Map(func(r rune) rune {
			if r >= '0' && r <= '9' {
				return r
			}
			return -1
		}, phone)

		if len(numericPhone) >= 10 {
			// Assume standard Brazil format with DDD (2 digits)
			payer.Phone.AreaCode = numericPhone[:2]
			payer.Phone.Number = numericPhone[2:]
		} else {
			payer.Phone.Number = numericPhone
		}
	}

	// Address
	if company.Address != nil {
		payer.Address.ZipCode = company.Address.Cep
		payer.Address.StreetName = company.Address.Street
		payer.Address.StreetNumber = company.Address.Number
		payer.Address.Neighborhood = company.Address.Neighborhood
		payer.Address.City = company.Address.City
		payer.Address.State = company.Address.UF
	}

	return payer
}

func checkoutPaymentEntity() entity.Entity {
	return entity.NewEntity()
}
