package handlerimpl

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	billingdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/checkout"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
	billingusecases "github.com/willjrcom/sales-backend-go/internal/usecases/checkout"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerCompanyImpl struct {
	s           *companyusecases.Service
	checkoutUC  *billingusecases.CheckoutUseCase
	costService *companyusecases.UsageCostService
}

func NewHandlerCompany(companyService *companyusecases.Service, checkoutUC *billingusecases.CheckoutUseCase, costService *companyusecases.UsageCostService) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerCompanyImpl{
		s:           companyService,
		checkoutUC:  checkoutUC,
		costService: costService,
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
		c.Get("/payments", h.handlerListCompanyPayments)
		c.Get("/costs/monthly", h.handlerGetMonthlyCosts)
		c.Post("/costs/register", h.handlerCreateCost)
	})

	return handler.NewHandler("/company", c)
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

	page, perPage := headerservice.GetPageAndPerPage(r, 1, 10)
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

	summary, err := h.costService.GetMonthlySummary(ctx, month, year)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, summary)
}
