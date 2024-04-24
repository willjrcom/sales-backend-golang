package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	tableorderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/table_order"
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
		c.Post("/update/change-table/{id}", h.handlerChangeTable)
		c.Post("/update/finish/{id}", h.handlerFinishTableOrder)
		c.Delete("/{id}", h.handlerDeleteTableOrderById)
		c.Get("/{id}", h.handlerGetTableOrderById)
		c.Get("/all", h.handlerGetAllTables)
	})

	return handler.NewHandler("/table-order", c)
}

func (h *handlerTableOrderImpl) handlerRegisterTableOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	table := &tableorderdto.CreateTableOrderInput{}
	jsonpkg.ParseBody(r, table)

	id, err := h.s.CreateTableOrder(ctx, table)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerTableOrderImpl) handlerDeleteTableOrderById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteTableOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerTableOrderImpl) handlerChangeTable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	table := &tableorderdto.UpdateTableOrderInput{}
	jsonpkg.ParseBody(r, table)

	if err := h.s.ChangeTable(ctx, dtoId, table); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerTableOrderImpl) handlerFinishTableOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.FinishTableOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerTableOrderImpl) handlerGetTableOrderById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	table, err := h.s.GetTableById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: table})
}

func (h *handlerTableOrderImpl) handlerGetAllTables(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orders, err := h.s.GetAllTables(ctx)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: orders})
}
