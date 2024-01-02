package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	processdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/process_category"
	processusecases "github.com/willjrcom/sales-backend-go/internal/usecases/process_category"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerProcessCategoryImpl struct {
	s *processusecases.Service
}

func NewHandlerProcessCategory(processService *processusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerProcessCategoryImpl{
		s: processService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterProcess)
		c.Patch("/update/{id}", h.handlerUpdateProcess)
		c.Delete("/delete/{id}", h.handlerDeleteProcess)
	})

	return handler.NewHandler("/category-product/process", c)
}

func (h *handlerProcessCategoryImpl) handlerRegisterProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	process := &processdto.RegisterProcessInput{}
	jsonpkg.ParseBody(r, process)

	if id, err := h.s.RegisterProcess(ctx, process); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerProcessCategoryImpl) handlerUpdateProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	Process := &processdto.UpdateProcessInput{}
	jsonpkg.ParseBody(r, Process)

	if err := h.s.UpdateProcess(ctx, dtoId, Process); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerProcessCategoryImpl) handlerDeleteProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteProcess(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
