package handlerimpl

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"

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
	scheduler   *scheduler.MonthlyBillingScheduler
}

func NewHandlerCompany(companyService *companyusecases.Service, checkoutUC *billingusecases.CheckoutUseCase, costService *companyusecases.UsageCostService, scheduler *scheduler.MonthlyBillingScheduler) *handler.Handler {
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
		c.Get("/users", h.handlerGetCompanyUsers)

		// Checkout
		c.Post("/checkout/subscription", h.handlerCheckoutCreateSubscription)
		c.Post("/checkout/costs", h.handlerCheckoutCosts)
		c.Post("/checkout/cancel/{paymentID}", h.handlerCancelPayment)
		c.Get("/payments", h.handlerListCompanyPayments)
		c.Get("/costs/monthly", h.handlerGetMonthlyCosts)
		c.Post("/costs/register", h.handlerCreateCost)
		c.Post("/payments/mercadopago/webhook", h.handlerMercadoPagoWebhook)
		c.Post("/billing/scheduler/trigger", h.handlerTriggerMonthlyBilling)
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

	id, err := h.s.GetCompany(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, id)
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
	payments, total, err := h.s.ListCompanyPayments(ctx, page, perPage)
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

	dto := &billingdto.CreateCheckoutDTO{}
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

func (h *handlerCompanyImpl) handlerCheckoutCosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	companyID, err := h.s.GetCompany(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	resp, err := h.checkoutUC.CreateCostCheckout(ctx, companyID.ID)
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

	if err := h.costService.RegisterUsageCost(ctx, dto.ToDomain()); err != nil {
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
	ctx := r.Context()

	dto := &companydto.MercadoPagoWebhookDTO{}
	// Parse body manually or use helper. The payload from MP matches the DTO
	if err := jsonpkg.ParseBody(r, dto); err != nil {
		// Log error but MP might retry
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	// Extract headers for signature validation details
	// x-signature: ...
	// x-request-id: ...
	dto.XSignature = r.Header.Get("x-signature")
	dto.XRequestID = r.Header.Get("x-request-id")
	dto.DataIDFromQuery = r.URL.Query().Get("data.id")

	if err := h.checkoutUC.HandleMercadoPagoWebhook(ctx, dto); err != nil {
		if err == billingusecases.ErrInvalidWebhookSecret {
			jsonpkg.ResponseErrorJson(w, r, http.StatusUnauthorized, err)
			return
		}
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

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

	if err := h.checkoutUC.CancelPayment(ctx, paymentID); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerCompanyImpl) handlerTriggerMonthlyBilling(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.scheduler.ProcessDailyBatch(ctx)
	jsonpkg.ResponseJson(w, r, http.StatusOK, map[string]string{"status": "triggered, daily batch started"})
}
