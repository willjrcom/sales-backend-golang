package handlerimpl

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	billingdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/checkout"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	"github.com/willjrcom/sales-backend-go/internal/infra/scheduler"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
	billingusecases "github.com/willjrcom/sales-backend-go/internal/usecases/checkout"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerCompanyImpl struct {
	s           *companyusecases.Service
	checkoutUC  *billingusecases.CheckoutUseCase
	costService *companyusecases.UsageCostService
	scheduler   *scheduler.DailyScheduler
}

func NewHandlerCompany(companyService *companyusecases.Service, checkoutUC *billingusecases.CheckoutUseCase, costService *companyusecases.UsageCostService, scheduler *scheduler.DailyScheduler) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerCompanyImpl{
		s:           companyService,
		checkoutUC:  checkoutUC,
		costService: costService,
		scheduler:   scheduler,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerNewCompany)
		c.Put("/update", h.handlerUpdateCompany)
		c.Get("/", h.handlerGetCompany)

		// Users
		c.Post("/add/user", h.handlerAddUserToCompany)
		c.Post("/remove/user", h.handlerRemoveUserFromCompany)
		c.Get("/user", h.handlerGetCompanyUsers)

		// Checkout
		c.Post("/payment/cancel/{paymentID}", h.handlerCancelPayment)
		c.Get("/payment", h.handlerListCompanyPayments)
		c.Get("/cost/monthly", h.handlerGetMonthlyCosts)
		c.Post("/cost/register", h.handlerCreateCost)

		// Scheduler
		c.Post("/billing/scheduler/trigger", h.handlerTriggerMonthlyBilling)

		// Subscription
		c.Post("/subscription/cancel", h.handlerCancelSubscription)
		c.Post("/subscription/checkout", h.handlerCheckoutCreateSubscription)
		c.Post("/subscription/checkout/upgrade", h.handlerCreateUpgradeCheckout)
		c.Get("/subscription/simulate/upgrade", h.handlerSimulateUpgrade)
		c.Get("/subscription/status", h.handlerGetSubscriptionStatus)

		// Webhook
		c.Post("/payments/mercadopago/webhook", h.handlerMercadoPagoWebhook)
	})

	return handler.NewHandler("/company", c, "/company/payments/mercadopago/webhook")
}

func (h *handlerCompanyImpl) handlerNewCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoCompany := &companydto.CompanyCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoCompany); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	companySchema, err := h.s.NewCompany(ctx, dtoCompany)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, companySchema)
}

func (h *handlerCompanyImpl) handlerUpdateCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoCompany := &companydto.CompanyUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dtoCompany); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdateCompany(ctx, dtoCompany); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerCompanyImpl) handlerGetCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	company, err := h.s.GetCompany(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, company)
}

func (h *handlerCompanyImpl) handlerGetCompanyUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// parse pagination query params
	page, perPage := headerservice.GetPageAndPerPage(r, 0, 10)

	users, total, err := h.s.GetCompanyUsers(ctx, page, perPage)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	jsonpkg.ResponseJson(w, r, http.StatusOK, users)
}

func (h *handlerCompanyImpl) handlerListCompanyPayments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page, perPage := headerservice.GetPageAndPerPage(r, 0, 10)
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))

	payments, total, err := h.s.ListCompanyPayments(ctx, page, perPage, month, year)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	jsonpkg.ResponseJson(w, r, http.StatusOK, payments)
}

func (h *handlerCompanyImpl) handlerAddUserToCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &companydto.UserToCompanyDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.AddUserToCompany(ctx, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerCompanyImpl) handlerRemoveUserFromCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &companydto.UserToCompanyDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.RemoveUserFromCompany(ctx, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerCompanyImpl) handlerCheckoutCreateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto := &billingdto.CreateSubscriptionCheckoutDTO{}
	if err := jsonpkg.ParseBody(r, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	resp, err := h.checkoutUC.CreateSubscriptionCheckout(ctx, dto)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerCompanyImpl) handlerCreateCost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto := &companydto.CompanyUsageCostCreateDTO{}
	if err := jsonpkg.ParseBody(r, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.costService.RegisterUsageCost(ctx, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerCompanyImpl) handlerGetMonthlyCosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	page, perPage := headerservice.GetPageAndPerPage(r, 0, 1000)

	summary, err := h.costService.GetMonthlySummary(ctx, month, year, page, perPage)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, summary)
}

func (h *handlerCompanyImpl) handlerMercadoPagoWebhook(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("========== MERCADO PAGO WEBHOOK RECEIVED ==========\n")
	fmt.Printf("URL: %s\n", r.URL.String())

	dto := &companydto.MercadoPagoWebhookDTO{}
	// Parse body manually or use helper. The payload from MP matches the DTO
	if err := jsonpkg.ParseBody(r, dto); err != nil {
		// Log error but MP might retry
		fmt.Printf("ERROR: Failed to parse webhook body: %v\n", err)
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	fmt.Printf("Webhook Type: %s\n", dto.Type)
	fmt.Printf("Action: %s\n", dto.Action)
	if dto.Data.ID != "" {
		fmt.Printf("Data.ID: %s\n", dto.Data.ID)
	}
	fmt.Printf("Query data.id: %s\n", r.URL.Query().Get("data.id"))

	// Extract headers for signature validation details
	// x-signature: ...
	// x-request-id: ...
	dto.XSignature = r.Header.Get("x-signature")
	dto.XRequestID = r.Header.Get("x-request-id")
	dto.DataIDFromQuery = r.URL.Query().Get("data.id")

	fmt.Printf("Processing webhook...\n")
	if err := h.checkoutUC.HandleMercadoPagoWebhook(context.Background(), dto); err != nil {
		if err == billingusecases.ErrInvalidWebhookSecret {
			fmt.Printf("ERROR: Invalid webhook secret %v", err)
			jsonpkg.ResponseErrorJson(w, r, http.StatusUnauthorized, err)
			return
		}
		fmt.Printf("ERROR: Failed to handle webhook: %v\n", err)
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	fmt.Printf("========== WEBHOOK PROCESSED SUCCESSFULLY ==========\n")

	jsonpkg.ResponseJson(w, r, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *handlerCompanyImpl) handlerCancelPayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	paymentIDStr := chi.URLParam(r, "paymentID")
	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, fmt.Errorf("invalid payment id: %w", err))
		return
	}

	if err := h.s.CancelPayment(ctx, paymentID); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerCompanyImpl) handlerTriggerMonthlyBilling(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.scheduler.UpdateCompanyPlans(ctx)
	h.scheduler.ProcessCostsToPay(ctx)
	h.scheduler.CheckOverdueAccounts(ctx)
	h.scheduler.CheckExpiredOptionalPayments(ctx)
	jsonpkg.ResponseJson(w, r, http.StatusOK, map[string]string{"status": "triggered, daily batch started"})
}

func (h *handlerCompanyImpl) handlerGetSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	status, err := h.s.GetSubscriptionStatus(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	// Populate Available Plans
	plans := companyentity.GetAllPlans()
	dtoPlans := make([]companydto.PlanDTO, len(plans))

	// Map current plan order
	currentPlanKey := status.CurrentPlan
	var currentOrder int
	for _, p := range plans {
		if string(p.Key) == currentPlanKey {
			currentOrder = p.Order
			break
		}
	}

	for i, p := range plans {
		isUpgrade := p.Order > currentOrder
		var upgradePrice *float64

		// Calculate upgrade price if applicable (only if user has paid plan and it is an upgrade)
		if isUpgrade && currentOrder > 0 {
			// GetCompany above returns ID, we use it here
			if sim, err := h.checkoutUC.CalculateUpgradeProration(ctx, p.Key); err == nil {
				upgradePrice = &sim.UpgradeAmount
			}
		}

		dtoPlans[i] = companydto.PlanDTO{
			Key:          string(p.Key),
			Name:         p.Name,
			Price:        p.Price,
			Features:     p.Features,
			Order:        p.Order,
			IsCurrent:    string(p.Key) == currentPlanKey,
			IsUpgrade:    isUpgrade,
			UpgradePrice: upgradePrice,
		}
	}
	status.AvailablePlans = dtoPlans

	jsonpkg.ResponseJson(w, r, http.StatusOK, status)
}

func (h *handlerCompanyImpl) handlerCancelSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := h.checkoutUC.CancelSubscription(ctx); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, map[string]string{"status": "subscription cancelled"})
}

func (h *handlerCompanyImpl) handlerSimulateUpgrade(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	targetPlan := r.URL.Query().Get("target_plan")
	if targetPlan == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, fmt.Errorf("target_plan required"))
		return
	}

	simulation, err := h.checkoutUC.CalculateUpgradeProration(ctx, companyentity.PlanType(targetPlan))
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, simulation)
}

func (h *handlerCompanyImpl) handlerCreateUpgradeCheckout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto := &billingdto.UpgradeCheckoutDTO{}
	if err := jsonpkg.ParseBody(r, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	resp, err := h.checkoutUC.CreateUpgradeCheckout(ctx, companyentity.PlanType(dto.TargetPlan))
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}
