package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	processruledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/process_rule"
	processruleusecases "github.com/willjrcom/sales-backend-go/internal/usecases/process_rule"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerProcessRuleCategoryImpl struct {
	s *processruleusecases.Service
}

func NewHandlerProcessRuleCategory(processRuleService *processruleusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerProcessRuleCategoryImpl{
		s: processRuleService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerCreateProcessRule)
		c.Patch("/update/{id}", h.handlerUpdateProcessRule)
		c.Delete("/{id}", h.handlerDeleteProcessRule)
		c.Get("/{id}", h.handlerGetProcessRuleById)
		c.Get("/by-category-id/{id}", h.handlerGetProcessRulesByCategoryID)
		c.Get("/all", h.handlerGetAllProcessRules)
	})

	return handler.NewHandler("/product-category/process-rule", c)
}

func (h *handlerProcessRuleCategoryImpl) handlerCreateProcessRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoProcessRule := &processruledto.CreateProcessRuleInput{}
	if err := jsonpkg.ParseBody(r, dtoProcessRule); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	id, err := h.s.CreateProcessRule(ctx, dtoProcessRule)
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

	dtoProcessRule := &processruledto.UpdateProcessRuleInput{}
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

func (h *handlerProcessRuleCategoryImpl) handlerGetProcessRulesByCategoryID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	processRules, err := h.s.GetProcessRulesByCategoryId(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: processRules})
}

func (h *handlerProcessRuleCategoryImpl) handlerGetAllProcessRules(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	processRules, err := h.s.GetAllProcessRules(ctx)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: processRules})
}
