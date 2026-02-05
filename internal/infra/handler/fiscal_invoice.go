package handlerimpl

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	fiscalinvoicedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/fiscal_invoice"
	fiscalinvoiceusecases "github.com/willjrcom/sales-backend-go/internal/usecases/fiscal_invoice"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerFiscalInvoiceImpl struct {
	service *fiscalinvoiceusecases.Service
}

func NewHandlerFiscalInvoice(service *fiscalinvoiceusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerFiscalInvoiceImpl{
		service: service,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/nfce/emit", h.handlerEmitNFCe)
		c.Get("/nfce/{id}", h.handlerSearchNFCe)
		c.Get("/nfce", h.handlerListNFCe)
		c.Post("/nfce/{id}/cancel", h.handlerCancelNFCe)
	})

	return handler.NewHandler("/fiscal", c)
}

// handlerEmitNFCe godoc
// @Summary Emit NFC-e for order
// @Description Emit electronic fiscal coupon for an order
// @Tags Fiscal Invoice
// @Accept json
// @Produce json
// @Param request body fiscalinvoicedto.EmitNFCeRequestDTO true "Order ID"
// @Success 200 {object} fiscalinvoicedto.FiscalInvoiceDTO
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /api/fiscal/nfce/emit [post]
func (h *handlerFiscalInvoiceImpl) handlerEmitNFCe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto := &fiscalinvoicedto.EmitNFCeRequestDTO{}
	if err := jsonpkg.ParseBody(r, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	invoice, err := h.service.EmitNFCeOrder(ctx, dto.OrderID)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case fiscalinvoiceusecases.ErrFiscalNotEnabled:
			status = http.StatusForbidden
		case fiscalinvoiceusecases.ErrMissingFiscalData:
			status = http.StatusBadRequest
		case fiscalinvoiceusecases.ErrInvoiceAlreadyExists:
			status = http.StatusConflict
		}
		jsonpkg.ResponseErrorJson(w, r, status, err)
		return
	}

	response := &fiscalinvoicedto.FiscalInvoiceDTO{}
	response.FromDomain(invoice)
	jsonpkg.ResponseJson(w, r, http.StatusOK, response)
}

// handlerSearchNFCe godoc
// @Summary Query NFC-e
// @Description Get fiscal invoice details
// @Tags Fiscal Invoice
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} fiscalinvoicedto.FiscalInvoiceDTO
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /api/fiscal/nfce/{id} [get]
func (h *handlerFiscalInvoiceImpl) handlerSearchNFCe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	invoice, err := h.service.SearchNFCe(ctx, id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == fiscalinvoiceusecases.ErrInvoiceNotFound {
			status = http.StatusNotFound
		}
		jsonpkg.ResponseErrorJson(w, r, status, err)
		return
	}

	response := &fiscalinvoicedto.FiscalInvoiceDTO{}
	response.FromDomain(invoice)
	jsonpkg.ResponseJson(w, r, http.StatusOK, response)
}

// handlerListNFCe godoc
// @Summary List NFC-e
// @Description List fiscal invoices for the company
// @Tags Fiscal Invoice
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {array} fiscalinvoicedto.FiscalInvoiceDTO
// @Failure 500 {object} error
// @Router /api/fiscal/nfce [get]
func (h *handlerFiscalInvoiceImpl) handlerListNFCe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")

	page, _ := strconv.Atoi(pageStr)
	perPage, _ := strconv.Atoi(perPageStr)

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	invoices, total, err := h.service.ListInvoices(ctx, page, perPage)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	dtos := make([]fiscalinvoicedto.FiscalInvoiceDTO, len(invoices))
	for i, invoice := range invoices {
		dto := fiscalinvoicedto.FiscalInvoiceDTO{}
		dto.FromDomain(invoice)
		dtos[i] = dto
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	jsonpkg.ResponseJson(w, r, http.StatusOK, dtos)
}

// handlerCancelNFCe godoc
// @Summary Cancel NFC-e
// @Description Cancel an authorized fiscal invoice
// @Tags Fiscal Invoice
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Param request body fiscalinvoicedto.CancelNFCeRequestDTO true "Cancellation reason"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} error
// @Failure 404 {object} error
// @Failure 500 {object} error
// @Router /api/fiscal/nfce/{id}/cancel [post]
func (h *handlerFiscalInvoiceImpl) handlerCancelNFCe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	dto := &fiscalinvoicedto.CancelNFCeRequestDTO{}
	if err := jsonpkg.ParseBody(r, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.service.CancelNFCe(ctx, id, dto.Justification); err != nil {
		status := http.StatusInternalServerError
		switch err {
		case fiscalinvoiceusecases.ErrInvoiceNotFound:
			status = http.StatusNotFound
		case fiscalinvoiceusecases.ErrCannotCancelInvoice:
			status = http.StatusBadRequest
		}
		jsonpkg.ResponseErrorJson(w, r, status, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, map[string]interface{}{
		"message": "NFC-e cancelada com sucesso",
	})
}
