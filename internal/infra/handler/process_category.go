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

type handlerProcessRuleCategoryImpl struct {
	s *processusecases.Service
}

func NewHandlerProcessRuleCategory(processRuleService *processusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerProcessRuleCategoryImpl{
		s: processRuleService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterProcessRule)
		c.Patch("/update/{id}", h.handlerUpdateProcessRule)
		c.Delete("/{id}", h.handlerDeleteProcessRule)
		c.Get("/{id}", h.handlerGetProcessRuleById)
	})

	return handler.NewHandler("/category-product/process", c)
}

func (h *handlerProcessRuleCategoryImpl) handlerRegisterProcessRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoProcessRule := &processdto.RegisterProcessRuleInput{}
	if err := jsonpkg.ParseBody(r, dtoProcessRule); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	id, err := h.s.RegisterProcessRule(ctx, dtoProcessRule)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerProcessRuleCategoryImpl) handlerUpdateProcessRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	dtoProcessRule := &processdto.UpdateProcessRuleInput{}
	if err := jsonpkg.ParseBody(r, dtoProcessRule); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.s.UpdateProcessRule(ctx, dtoId, dtoProcessRule); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerProcessRuleCategoryImpl) handlerDeleteProcessRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteProcessRule(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerProcessRuleCategoryImpl) handlerGetProcessRuleById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	processRule, err := h.s.GetProcessRuleById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: processRule})
}
