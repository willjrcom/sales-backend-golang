package handlerimpl

import (
	"errors"
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
		c.Post("/new", h.handlerCreateTable)
		c.Delete("/{id}", h.handlerDeleteTableById)
		c.Patch("/update/{id}", h.handlerUpdateTableById)
		c.Get("/{id}", h.handlerGetTableById)
		c.Get("/all", h.handlerGetAllTables)
		c.Get("/all/unused", h.handlerGetUnusedTables)
	})

	return handler.NewHandler("/table", c)
}

func (h *handlerTableImpl) handlerCreateTable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoTable := &tabledto.TableCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoTable); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.s.CreateTable(ctx, dtoTable)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, id)
}

func (h *handlerTableImpl) handlerDeleteTableById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteTable(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerTableImpl) handlerGetTableById(w http.ResponseWriter, r *http.Request) {
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

func (h *handlerTableImpl) handlerUpdateTableById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoTable := &tabledto.TableUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dtoTable); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdateTable(ctx, dtoId, dtoTable); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerTableImpl) handlerGetAllTables(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tables, err := h.s.GetAllTables(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, tables)
}

func (h *handlerTableImpl) handlerGetUnusedTables(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tables, err := h.s.GetUnusedTables(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, tables)
}
