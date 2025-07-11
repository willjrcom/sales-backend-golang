package handlerimpl

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	reportdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/report"
	reportusecases "github.com/willjrcom/sales-backend-go/internal/usecases/report"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

// NewHandlerReport returns a handler for report endpoints.
func NewHandlerReport(s *reportusecases.Service) *handler.Handler {
	r := chi.NewRouter()
	h := &handlerReportImpl{s: s}
	base := "/report"

	r.Post("/sales-total-by-day", h.handleSalesTotalByDay)
	r.Post("/revenue-cumulative-by-month", h.handleRevenueCumulativeByMonth)
	r.Post("/sales-by-hour", h.handleSalesByHour)
	r.Post("/sales-by-channel", h.handleSalesByChannel)
	r.Post("/avg-ticket-by-day", h.handleAvgTicketByDay)
	r.Post("/avg-ticket-by-channel", h.handleAvgTicketByChannel)
	r.Post("/products-sold-by-day", h.handleProductsSoldByDay)
	r.Post("/top-products", h.handleTopProducts)
	r.Post("/sales-by-category", h.handleSalesByCategory)
	r.Post("/clients-registered-by-day", h.handleClientsRegisteredByDay)
	r.Post("/new-vs-recurring-clients", h.handleNewVsRecurringClients)
	r.Get("/orders-by-status", h.handleOrdersByStatus)
	r.Get("/avg-process-step-duration", h.handleAvgProcessStepDurationByRule)
	r.Get("/cancellation-rate", h.handleCancellationRate)
	r.Get("/current-queue-length", h.handleCurrentQueueLength)
	r.Get("/avg-delivery-time-by-driver", h.handleAvgDeliveryTimeByDriver)
	r.Get("/deliveries-per-driver", h.handleDeliveriesPerDriver)
	r.Get("/orders-per-table", h.handleOrdersPerTable)
	r.Get("/avg-queue-duration", h.handleAvgQueueDuration)
	r.Get("/avg-process-duration-by-product", h.handleAvgProcessDurationByProduct)
	r.Get("/total-queue-time-by-group-item", h.handleTotalQueueTimeByGroupItem)
	// Sales by shift
	r.Post("/sales-by-shift", h.handleSalesByShift)
	// Top 10 mesas mais utilizadas
	r.Post("/top-tables", h.handleTopTables)
	r.Post("/payments-by-method", h.handlePaymentsByMethod)
	// Custom reports 26–33
	r.Post("/sales-by-place", h.handleSalesByPlace)
	r.Post("/sales-by-size", h.handleSalesBySize)
	r.Post("/additional-items-sold", h.handleAdditionalItemsSold)
	r.Post("/avg-pickup-time", h.handleAvgPickupTime)
	r.Post("/group-items-status", h.handleGroupItemsByStatus)
	r.Post("/deliveries-by-cep", h.handleDeliveriesByCep)
	r.Post("/processed-count-by-rule", h.handleProcessedCountByRule)
	r.Post("/employee-payments-report", h.handleEmployeePaymentsReport)
	// Daily sales report for a specific day
	r.Post("/daily-sales", h.handleDailySales)
	// Profitability reports
	r.Post("/product-profitability", h.handleProductProfitability)
	r.Post("/category-profitability", h.handleCategoryProfitability)
	r.Post("/low-profit-products", h.handleLowProfitProducts)
	r.Post("/overall-profitability", h.handleOverallProfitability)
	r.Post("/debug-products", h.handleDebugProducts) // Endpoint temporário para debug

	unprotected := []string{
		"/report/debug-products", // Temporariamente desprotegido para debug
	}
	return handler.NewHandler(base, r, unprotected...)
}

// handleDailySales handles daily sales report generation.
func (h *handlerReportImpl) handleDailySales(w http.ResponseWriter, r *http.Request) {
	var req reportdto.DailySalesRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.DailySales(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

type handlerReportImpl struct {
	s *reportusecases.Service
}

func parseBody(r *http.Request, v interface{}, w http.ResponseWriter) bool {
	if err := jsonpkg.ParseBody(r, v); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return false
	}
	return true
}

// Handler implementations:
func (h *handlerReportImpl) handleSalesTotalByDay(w http.ResponseWriter, r *http.Request) {
	var req reportdto.SalesTotalByDayRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.SalesTotalByDay(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleRevenueCumulativeByMonth(w http.ResponseWriter, r *http.Request) {
	var req reportdto.RevenueCumulativeByMonthRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.RevenueCumulativeByMonth(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleSalesByHour(w http.ResponseWriter, r *http.Request) {
	var req reportdto.SalesByHourRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.SalesByHour(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleSalesByChannel(w http.ResponseWriter, r *http.Request) {
	var req reportdto.SalesByChannelRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.SalesByChannel(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleAvgTicketByDay(w http.ResponseWriter, r *http.Request) {
	var req reportdto.AvgTicketByDayRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.AvgTicketByDay(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleAvgTicketByChannel(w http.ResponseWriter, r *http.Request) {
	var req reportdto.AvgTicketByChannelRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.AvgTicketByChannel(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleProductsSoldByDay(w http.ResponseWriter, r *http.Request) {
	var req reportdto.ProductsSoldByDayRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.ProductsSoldByDay(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleTopProducts(w http.ResponseWriter, r *http.Request) {
	var req reportdto.TopProductsRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.TopProducts(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleSalesByCategory(w http.ResponseWriter, r *http.Request) {
	var req reportdto.SalesByCategoryRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.SalesByCategory(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleClientsRegisteredByDay(w http.ResponseWriter, r *http.Request) {
	var req reportdto.ClientsRegisteredByDayRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.ClientsRegisteredByDay(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleNewVsRecurringClients(w http.ResponseWriter, r *http.Request) {
	var req reportdto.NewVsRecurringClientsRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.NewVsRecurringClients(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleOrdersByStatus(w http.ResponseWriter, r *http.Request) {
	schema := r.URL.Query().Get("schema")
	if schema == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("schema is required"))
		return
	}
	req := reportdto.OrdersByStatusRequest{}
	resp, err := h.s.OrdersByStatus(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleAvgProcessStepDurationByRule(w http.ResponseWriter, r *http.Request) {
	schema := r.URL.Query().Get("schema")
	if schema == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("schema is required"))
		return
	}
	req := reportdto.AvgProcessStepDurationRequest{}
	resp, err := h.s.AvgProcessStepDurationByRule(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleCancellationRate(w http.ResponseWriter, r *http.Request) {
	schema := r.URL.Query().Get("schema")
	if schema == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("schema is required"))
		return
	}
	req := reportdto.CancellationRateRequest{}
	resp, err := h.s.CancellationRate(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleCurrentQueueLength(w http.ResponseWriter, r *http.Request) {
	schema := r.URL.Query().Get("schema")
	if schema == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("schema is required"))
		return
	}
	req := reportdto.CurrentQueueLengthRequest{}
	resp, err := h.s.CurrentQueueLength(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleAvgDeliveryTimeByDriver(w http.ResponseWriter, r *http.Request) {
	schema := r.URL.Query().Get("schema")
	if schema == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("schema is required"))
		return
	}
	req := reportdto.AvgDeliveryTimeByDriverRequest{}
	resp, err := h.s.AvgDeliveryTimeByDriver(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleDeliveriesPerDriver(w http.ResponseWriter, r *http.Request) {
	schema := r.URL.Query().Get("schema")
	if schema == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("schema is required"))
		return
	}
	req := reportdto.DeliveriesPerDriverRequest{}
	resp, err := h.s.DeliveriesPerDriver(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleOrdersPerTable(w http.ResponseWriter, r *http.Request) {
	schema := r.URL.Query().Get("schema")
	if schema == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("schema is required"))
		return
	}
	req := reportdto.OrdersPerTableRequest{}
	resp, err := h.s.OrdersPerTable(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

// handleAvgQueueDuration handles average queue duration report.
func (h *handlerReportImpl) handleAvgQueueDuration(w http.ResponseWriter, r *http.Request) {
	schema := r.URL.Query().Get("schema")
	if schema == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("schema is required"))
		return
	}
	req := reportdto.AvgQueueDurationRequest{}
	resp, err := h.s.AvgQueueDuration(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

// handleAvgProcessDurationByProduct handles average process duration by product report.
func (h *handlerReportImpl) handleAvgProcessDurationByProduct(w http.ResponseWriter, r *http.Request) {
	schema := r.URL.Query().Get("schema")
	if schema == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("schema is required"))
		return
	}
	req := reportdto.AvgProcessDurationByProductRequest{}
	resp, err := h.s.AvgProcessDurationByProduct(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

// handleTotalQueueTimeByGroupItem handles total queue time per group item report.
func (h *handlerReportImpl) handleTotalQueueTimeByGroupItem(w http.ResponseWriter, r *http.Request) {
	schema := r.URL.Query().Get("schema")
	if schema == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("schema is required"))
		return
	}
	req := reportdto.TotalQueueTimeByGroupItemRequest{}
	resp, err := h.s.TotalQueueTimeByGroupItem(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleSalesByShift(w http.ResponseWriter, r *http.Request) {
	var req reportdto.SalesByShiftRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.SalesByShift(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

// handleTopTables handles the report for top 10 most used tables.
func (h *handlerReportImpl) handleTopTables(w http.ResponseWriter, r *http.Request) {
	var req reportdto.TopTablesRequest
	if !parseBody(r, &req, w) {
		return
	}
	resp, err := h.s.TopTables(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handlePaymentsByMethod(w http.ResponseWriter, r *http.Request) {
	var req reportdto.PaymentsByMethodRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.PaymentsByMethod(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleEmployeePaymentsReport(w http.ResponseWriter, r *http.Request) {
	var req reportdto.EmployeePaymentsReportRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.EmployeePaymentsReport(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

// Custom reports handlers:
func (h *handlerReportImpl) handleSalesByPlace(w http.ResponseWriter, r *http.Request) {
	var req reportdto.SalesByPlaceRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.SalesByPlace(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleSalesBySize(w http.ResponseWriter, r *http.Request) {
	var req reportdto.SalesBySizeRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.SalesBySize(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleAdditionalItemsSold(w http.ResponseWriter, r *http.Request) {
	var req reportdto.AdditionalItemsRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.AdditionalItemsSold(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleAvgPickupTime(w http.ResponseWriter, r *http.Request) {
	var req reportdto.AvgPickupTimeRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.AvgPickupTime(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleGroupItemsByStatus(w http.ResponseWriter, r *http.Request) {
	var req reportdto.GroupItemsByStatusRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.GroupItemsByStatus(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleDeliveriesByCep(w http.ResponseWriter, r *http.Request) {
	var req reportdto.DeliveriesByCepRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.DeliveriesByCep(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleProcessedCountByRule(w http.ResponseWriter, r *http.Request) {
	var req reportdto.ProcessedCountByRuleRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.ProcessedCountByRule(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

// Profitability report handlers:
func (h *handlerReportImpl) handleProductProfitability(w http.ResponseWriter, r *http.Request) {
	var req reportdto.ProductProfitabilityRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.ProductProfitability(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleCategoryProfitability(w http.ResponseWriter, r *http.Request) {
	var req reportdto.CategoryProfitabilityRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.CategoryProfitability(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleLowProfitProducts(w http.ResponseWriter, r *http.Request) {
	var req reportdto.LowProfitProductsRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.LowProfitProducts(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleOverallProfitability(w http.ResponseWriter, r *http.Request) {
	var req reportdto.OverallProfitabilityRequest
	if !parseBody(r, &req, w) {
		return
	}

	resp, err := h.s.OverallProfitability(r.Context(), &req)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}

func (h *handlerReportImpl) handleDebugProducts(w http.ResponseWriter, r *http.Request) {
	resp, err := h.s.DebugProducts(r.Context())
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, resp)
}
