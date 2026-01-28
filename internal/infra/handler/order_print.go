package handlerimpl

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderprintusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_print"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

// handlerOrderPrintImpl implements print endpoints for orders.
type handlerOrderPrintImpl struct {
	s *orderprintusecases.Service
}

// NewHandlerOrderPrint returns a handler for printing individual orders.
func NewHandlerOrderPrint(svc *orderprintusecases.Service) *handler.Handler {
	r := chi.NewRouter()
	h := &handlerOrderPrintImpl{s: svc}
	r.With().Group(func(r chi.Router) {
		// Kitchen print: only items and complements
		r.Get("/kitchen/{id}", h.handlePrintGroupItemKitchen)
		// Full order print
		r.Get("/{id}", h.handlePrintOrder)
		// Daily report print
		r.Post("/daily", h.handlePrintByShift)
	})
	return handler.NewHandler("/order-print", r)
}

// handlePrintOrder handles GET /order-print/{id}
func (h *handlerOrderPrintImpl) handlePrintOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}
	format := r.URL.Query().Get("format")

	var data []byte
	var err error

	if format == "html" {
		data, err = h.s.PrintOrderHTML(ctx, dtoId)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else {
		data, err = h.s.PrintOrder(ctx, dtoId)
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

// handlePrintByShift handles POST /order-print/daily to generate daily sales report.
func (h *handlerOrderPrintImpl) handlePrintByShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	res, err := h.s.PrintDailyReport(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, res)
}

// handlePrintGroupItemKitchen handles GET /order-print/kitchen/{id} for kitchen tickets.
func (h *handlerOrderPrintImpl) handlePrintGroupItemKitchen(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}
	format := r.URL.Query().Get("format")

	var data []byte
	var err error

	if format == "html" {
		data, err = h.s.PrintGroupItemKitchenHTML(ctx, dtoId)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else {
		data, err = h.s.PrintGroupItemKitchen(ctx, dtoId)
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
