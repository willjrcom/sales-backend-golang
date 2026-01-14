package handlerimpl

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerCompanyUsageCostImpl struct {
	usageCostService *companyusecases.UsageCostService
	companyService   *companyusecases.Service
}

func NewHandlerCompanyUsageCost(usageCostService *companyusecases.UsageCostService, companyService *companyusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerCompanyUsageCostImpl{
		usageCostService: usageCostService,
		companyService:   companyService,
	}

	c.With().Group(func(c chi.Router) {
		c.Get("/costs/monthly", h.handlerGetMonthlySummary)
		c.Get("/costs/breakdown", h.handlerGetCostBreakdown)
		c.Get("/next-invoice", h.handlerGetNextInvoicePreview)
		c.Post("/fiscal/enable", h.handlerEnableFiscalInvoice)
		c.Post("/fiscal/disable", h.handlerDisableFiscalInvoice)
	})

	return handler.NewHandler("/company-usage-cost", c)
}

// handlerGetMonthlySummary godoc
// @Summary Get monthly cost summary
// @Description Get summary of all costs for a specific month
// @Tags Company Costs
// @Accept json
// @Produce json
// @Param month query int false "Month (1-12, default: current month)"
// @Param year query int false "Year (default: current year)"
// @Success 200 {object} companydto.MonthlyCostSummaryDTO
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/company/costs/monthly [get]
func (h *handlerCompanyUsageCostImpl) handlerGetMonthlySummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monthStr := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	// Default to current month/year if not provided
	if month == 0 || year == 0 {
		now := time.Now()
		if month == 0 {
			month = int(now.Month())
		}
		if year == 0 {
			year = now.Year()
		}
	}

	// Validate
	if month < 1 || month > 12 {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, http.ErrBodyNotAllowed)
		return
	}
	if year < 2020 {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, http.ErrBodyNotAllowed)
		return
	}

	summary, err := h.usageCostService.GetMonthlySummary(ctx, month, year)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, summary)
}

// handlerGetCostBreakdown godoc
// @Summary Get detailed cost breakdown
// @Description Get detailed breakdown of all costs for a specific month
// @Tags Company Costs
// @Accept json
// @Produce json
// @Param month query int false "Month (1-12, default: current month)"
// @Param year query int false "Year (default: current year)"
// @Success 200 {object} companydto.CostBreakdownDTO
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/company/costs/breakdown [get]
func (h *handlerCompanyUsageCostImpl) handlerGetCostBreakdown(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monthStr := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	// Default to current month/year if not provided
	if month == 0 || year == 0 {
		now := time.Now()
		if month == 0 {
			month = int(now.Month())
		}
		if year == 0 {
			year = now.Year()
		}
	}

	// Validate
	if month < 1 || month > 12 {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, http.ErrBodyNotAllowed)
		return
	}
	if year < 2020 {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, http.ErrBodyNotAllowed)
		return
	}

	breakdown, err := h.usageCostService.GetCostBreakdown(ctx, month, year)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, breakdown)
}

// handlerGetNextInvoicePreview godoc
// @Summary Get next invoice preview
// @Description Get preview of the next billing invoice with enabled services and estimated costs
// @Tags Company Costs
// @Accept json
// @Produce json
// @Success 200 {object} companydto.NextInvoicePreviewDTO
// @Failure 500 {object} error
// @Router /api/company-usage-cost/next-invoice [get]
func (h *handlerCompanyUsageCostImpl) handlerGetNextInvoicePreview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	preview, err := h.usageCostService.GetNextInvoicePreview(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, preview)
}

// handlerEnableFiscalInvoice godoc

// @Summary Enable fiscal invoice functionality
// @Description Enable NFC-e emission for the company
// @Tags Company Fiscal
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/company/fiscal/enable [post]
func (h *handlerCompanyUsageCostImpl) handlerEnableFiscalInvoice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get company
	company, err := h.companyService.GetCompany(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	// Validate required fiscal data
	if company.InscricaoEstadual == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, fmt.Errorf("inscricao_estadual is required to enable fiscal invoices"))
		return
	}
	if company.RegimeTributario == 0 {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, fmt.Errorf("regime_tributario is required to enable fiscal invoices"))
		return
	}

	// Enable fiscal
	updateDTO := &companydto.CompanyUpdateDTO{
		FiscalEnabled: func() *bool { b := true; return &b }(),
	}

	if err := h.companyService.UpdateCompany(ctx, updateDTO); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	// Register monthly fiscal subscription fee (R$ 20.00)
	// This will only charge once per month even if enabled/disabled multiple times
	if err := h.usageCostService.RegisterFiscalSubscriptionFee(ctx, company.ID); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, map[string]interface{}{
		"message": "Fiscal invoice functionality enabled",
		"enabled": true,
	})
}

// handlerDisableFiscalInvoice godoc
// @Summary Disable fiscal invoice functionality
// @Description Disable NFC-e emission for the company
// @Tags Company Fiscal
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} error
// @Router /api/company/fiscal/disable [post]
func (h *handlerCompanyUsageCostImpl) handlerDisableFiscalInvoice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Disable fiscal
	updateDTO := &companydto.CompanyUpdateDTO{
		FiscalEnabled: func() *bool { b := false; return &b }(),
	}

	if err := h.companyService.UpdateCompany(ctx, updateDTO); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, map[string]interface{}{
		"message": "Fiscal invoice functionality disabled",
		"enabled": false,
	})
}
