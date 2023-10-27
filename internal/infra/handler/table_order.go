package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
	tableorderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/table_order"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerTableOrderImpl struct {
	s *tableorderusecases.Service
}

func NewHandlerTableOrder(orderService *tableorderusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerTableOrderImpl{
		s: orderService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterTableOrder)
		c.Get("/{id}", h.handlerGetTableById)
		c.Get("/all", h.handlerGetAllTables)
	})

	return handler.NewHandler("/table-order", c)
}

func (h *handlerTableOrderImpl) handlerRegisterTableOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	table := &orderdto.CreateTableOrderInput{}
	jsonpkg.ParseBody(r, table)

	if id, err := h.s.CreateTableOrder(ctx, table); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerTableOrderImpl) handlerGetTableById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if id, err := h.s.GetTableById(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerTableOrderImpl) handlerGetAllTables(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if orders, err := h.s.GetAllTables(ctx); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: orders})
	}
}
