package handlerimpl

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	stockdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/stock"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
	stockusecases "github.com/willjrcom/sales-backend-go/internal/usecases/stock"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerStockImpl struct {
	s *stockusecases.Service
}

func NewHandlerStock(stockService *stockusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerStockImpl{
		s: stockService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerCreateStock)
		c.Put("/update/{id}", h.handlerUpdateStock)
		c.Get("/{id}", h.handlerGetStockByID)
		c.Get("/{id}/with-product", h.handlerGetStockWithProduct)
		c.Get("/product/{product_id}", h.handlerGetStockByProductID)
		c.Get("/all", h.handlerGetAllStocks)
		c.Get("/all/with-product", h.handlerGetAllStocksWithProduct)
		// Movements
		c.Post("/{id}/movement/add", h.handlerAddStock)
		c.Post("/{id}/movement/remove", h.handlerRemoveStock)
		c.Post("/{id}/movement/adjust", h.handlerAdjustStock)
		c.Get("/movements/{stock_id}", h.handlerGetMovementsByStockID)

		// Alerts
		c.Get("/alerts", h.handlerGetAllAlerts)
		c.Get("/alerts/{id}", h.handlerGetAlertByID)
		c.Put("/alerts/{id}/resolve", h.handlerResolveAlert)
		c.Delete("/alerts/{id}", h.handlerDeleteAlert)

		// Reports
		c.Get("/report", h.handlerGetStockReport)
		c.Get("/low-stock", h.handlerGetLowStockProducts)
		c.Get("/out-of-stock", h.handlerGetOutOfStockProducts)
	})

	return handler.NewHandler("/stock", c)
}

func (h *handlerStockImpl) handlerCreateStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoStock := &stockdto.StockCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoStock); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	stock, err := h.s.CreateStock(ctx, dtoStock)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, stock)
}

func (h *handlerStockImpl) handlerUpdateStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoID := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoStock := &stockdto.StockUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dtoStock); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdateStock(ctx, dtoID, dtoStock); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerStockImpl) handlerGetStockByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoID := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	stock, err := h.s.GetStockByID(ctx, dtoID)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, stock)
}

func (h *handlerStockImpl) handlerGetStockWithProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoID := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	stock, err := h.s.GetStockWithProduct(ctx, dtoID)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, stock)
}

func (h *handlerStockImpl) handlerGetStockByProductID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	productID := chi.URLParam(r, "product_id")
	if productID == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("product_id is required"))
		return
	}

	stock, err := h.s.GetStockByProductID(ctx, productID)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, stock)
}

func (h *handlerStockImpl) handlerGetAllStocks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page, perPage := headerservice.GetPageAndPerPage(r, 0, 100)

	stocks, count, err := h.s.GetAllStocks(ctx, page, perPage)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(count))

	jsonpkg.ResponseJson(w, r, http.StatusOK, stocks)
}

func (h *handlerStockImpl) handlerGetAllStocksWithProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stocks, err := h.s.GetAllStocksWithProduct(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, stocks)
}

func (h *handlerStockImpl) handlerGetLowStockProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stocks, err := h.s.GetLowStockProducts(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, stocks)
}

func (h *handlerStockImpl) handlerGetOutOfStockProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stocks, err := h.s.GetOutOfStockProducts(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, stocks)
}

func (h *handlerStockImpl) handlerAddStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoMovement := &stockdto.StockMovementCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoMovement); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	movement, err := h.s.AddMovementStock(ctx, dtoId, dtoMovement)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, movement)
}

func (h *handlerStockImpl) handlerRemoveStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoMovement := &stockdto.StockMovementRemoveDTO{}
	if err := jsonpkg.ParseBody(r, dtoMovement); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	movement, err := h.s.RemoveMovementStock(ctx, dtoId, dtoMovement)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, movement)
}

func (h *handlerStockImpl) handlerAdjustStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	stockAdjust := &stockdto.StockMovementAdjustDTO{}
	if err := jsonpkg.ParseBody(r, &stockAdjust); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	movement, err := h.s.AdjustMovementStock(ctx, dtoId, stockAdjust)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, movement)
}

func (h *handlerStockImpl) handlerGetMovementsByStockID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stockID := chi.URLParam(r, "stock_id")
	if stockID == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("stock_id is required"))
		return
	}

	movements, err := h.s.GetMovementsByStockID(ctx, stockID)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, movements)
}

func (h *handlerStockImpl) handlerGetAllAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.s.GetAllAlerts(r.Context())
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, alerts)
}

func (h *handlerStockImpl) handlerGetAlertByID(w http.ResponseWriter, r *http.Request) {
	alertID := chi.URLParam(r, "id")
	if alertID == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("alert ID is required"))
		return
	}

	alert, err := h.s.GetAlertByID(r.Context(), alertID)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusNotFound, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, alert)
}

func (h *handlerStockImpl) handlerResolveAlert(w http.ResponseWriter, r *http.Request) {
	alertID := chi.URLParam(r, "id")
	if alertID == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("alert ID is required"))
		return
	}

	err := h.s.ResolveAlert(r.Context(), alertID)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, map[string]string{"message": "Alert resolved successfully"})
}

func (h *handlerStockImpl) handlerDeleteAlert(w http.ResponseWriter, r *http.Request) {
	alertID := chi.URLParam(r, "id")
	if alertID == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("alert ID is required"))
		return
	}

	err := h.s.DeleteAlert(r.Context(), alertID)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, map[string]string{"message": "Alert deleted successfully"})
}

func (h *handlerStockImpl) handlerGetStockReport(w http.ResponseWriter, r *http.Request) {
	page, perPage := headerservice.GetPageAndPerPage(r, 0, 100)
	report, count, err := h.s.GetStockReport(r.Context(), page, perPage)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(count))

	jsonpkg.ResponseJson(w, r, http.StatusOK, report)
}
