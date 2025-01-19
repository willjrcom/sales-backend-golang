package handlerimpl

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	ordertabledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_table"
	ordertableusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_table"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerOrderTableImpl struct {
	s *ordertableusecases.Service
}

func NewHandlerOrderTable(orderService *ordertableusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerOrderTableImpl{
		s: orderService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerCreateOrderTable)
		c.Post("/update/change-table/{id}", h.handlerChangeTable)
		c.Post("/update/close/{id}", h.handlerCloseOrderTable)
		c.Get("/{id}", h.handlerGetOrderTableById)
		c.Get("/all", h.handlerGetAllTables)
	})

	return handler.NewHandler("/order-table", c)
}

func (h *handlerOrderTableImpl) handlerCreateOrderTable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dtoTable := &ordertabledto.CreateOrderTableInput{}
	if err := jsonpkg.ParseBody(r, dtoTable); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.s.CreateOrderTable(ctx, dtoTable)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, id)
}

func (h *handlerOrderTableImpl) handlerChangeTable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoTable := &ordertabledto.OrderTableUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dtoTable); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.ChangeTable(ctx, dtoId, dtoTable); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderTableImpl) handlerCloseOrderTable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.CloseOrderTable(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderTableImpl) handlerGetOrderTableById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	table, err := h.s.GetTableById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, table)
}

func (h *handlerOrderTableImpl) handlerGetAllTables(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orders, err := h.s.GetAllTables(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, orders)
}
