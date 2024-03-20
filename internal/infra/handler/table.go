package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	tabledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/table"
	tableusecases "github.com/willjrcom/sales-backend-go/internal/usecases/table"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerTableImpl struct {
	s *tableusecases.Service
}

func NewHandlerTable(orderService *tableusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerTableImpl{
		s: orderService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterTable)
		c.Delete("/delete/{id}", h.handlerDeleteTableById)
		c.Get("/{id}", h.handlerGetTableById)
		c.Get("/all", h.handlerGetAllTables)
	})

	return handler.NewHandler("/table", c)
}

func (h *handlerTableImpl) handlerRegisterTable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	table := &tabledto.CreateTableInput{}
	jsonpkg.ParseBody(r, table)

	if id, err := h.s.CreateTable(ctx, table); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerTableImpl) handlerDeleteTableById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteTable(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerTableImpl) handlerGetTableById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if id, err := h.s.GetTableById(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerTableImpl) handlerGetAllTables(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if tables, err := h.s.GetAllTables(ctx); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: tables})
	}
}
