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
	})

	return handler.NewHandler("/category-product/process", c)
}

func (h *handlerProcessRuleCategoryImpl) handlerRegisterProcessRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	processRule := &processdto.RegisterProcessRuleInput{}
	jsonpkg.ParseBody(r, processRule)

	id, err := h.s.RegisterProcessRule(ctx, processRule)
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

	ProcessRule := &processdto.UpdateProcessRuleInput{}
	jsonpkg.ParseBody(r, ProcessRule)

	if err := h.s.UpdateProcessRule(ctx, dtoId, ProcessRule); err != nil {
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
