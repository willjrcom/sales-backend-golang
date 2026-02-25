package handlerimpl

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	printmanagerusecases "github.com/willjrcom/sales-backend-go/internal/usecases/print_manager"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

// handlerOrderPrintImpl implements print endpoints for orders.
type handlerOrderPrintImpl struct {
	s *printmanagerusecases.Service
}

// NewHandlerPrintManager returns a handler for printing individual orders.
func NewHandlerPrintManager(svc *printmanagerusecases.Service) *handler.Handler {
	r := chi.NewRouter()
	h := &handlerOrderPrintImpl{s: svc}
	r.With().Group(func(r chi.Router) {
		// Request print kitchen
		r.Post("/group-item/{id}", h.handleRequestPrintGroupItem)
		// Kitchen print: only items and complements
		r.Get("/group-item/{id}", h.handleGetPrintGroupItem)

		// Request print order
		r.Post("/order/{id}", h.handleRequestPrintOrder)
		// Full order print
		r.Get("/order/{id}", h.handleGetPrintOrder)

		// Request print shift
		r.Post("/shift/{id}", h.handleRequestPrintShift)
		// Shift report print
		r.Get("/shift/{id}", h.handleGetPrintShift)
	})
	return handler.NewHandler("/print-manager", r)
}

// handleRequestPrintOrder handles POST /print-manager/order/{id}
func (h *handlerOrderPrintImpl) handleRequestPrintOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.RequestPrintOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

// handleGetPrintOrder handles GET /print-manager/order/{id}
func (h *handlerOrderPrintImpl) handleGetPrintOrder(w http.ResponseWriter, r *http.Request) {
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

// handleRequestPrintOrder handles POST /print-manager/group-item/{id}
func (h *handlerOrderPrintImpl) handleRequestPrintGroupItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.RequestPrintGroupItemKitchen(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

// handleGetPrintGroupItem handles GET /print-manager/group-item/{id} for kitchen tickets.
func (h *handlerOrderPrintImpl) handleGetPrintGroupItem(w http.ResponseWriter, r *http.Request) {
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

// handleRequestPrintShift handles POST /print-manager/shift/{id}
func (h *handlerOrderPrintImpl) handleRequestPrintShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.RequestPrintShift(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

// handleGetPrintShift handles GET /print-manager/shift/{id}
func (h *handlerOrderPrintImpl) handleGetPrintShift(w http.ResponseWriter, r *http.Request) {
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
		data, err = h.s.PrintShiftHTML(ctx, dtoId)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else {
		data, err = h.s.PrintShift(ctx, dtoId)
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
