package handlerimpl

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	processruledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/process_rule"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
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
		c.Get("/all-with-order-process", h.handlerGetAllProcessRulesWithOrderProcess)
	})

	return handler.NewHandler("/product-category/process-rule", c)
}

func (h *handlerProcessRuleCategoryImpl) handlerCreateProcessRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoProcessRule := &processruledto.ProcessRuleCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoProcessRule); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.s.CreateProcessRule(ctx, dtoProcessRule)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, id)
}

func (h *handlerProcessRuleCategoryImpl) handlerUpdateProcessRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoProcessRule := &processruledto.ProcessRuleUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dtoProcessRule); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdateProcessRule(ctx, dtoId, dtoProcessRule); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerProcessRuleCategoryImpl) handlerDeleteProcessRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteProcessRule(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerProcessRuleCategoryImpl) handlerGetProcessRuleById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	processRule, err := h.s.GetProcessRuleById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, processRule)
}

func (h *handlerProcessRuleCategoryImpl) handlerGetProcessRulesByCategoryID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	processRules, err := h.s.GetProcessRulesByCategoryId(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, processRules)
}

func (h *handlerProcessRuleCategoryImpl) handlerGetAllProcessRules(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse pagination
	page, perPage := headerservice.GetPageAndPerPage(r, 0, 100)

	// Parse is_active query parameter (default: true)
	isActive := true
	if isActiveParam := r.URL.Query().Get("is_active"); isActiveParam != "" {
		var err error
		isActive, err = strconv.ParseBool(isActiveParam)
		if err != nil {
			jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("invalid is_active parameter"))
			return
		}
	}

	processRules, total, err := h.s.GetAllProcessRules(ctx, page, perPage, isActive)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	jsonpkg.ResponseJson(w, r, http.StatusOK, processRules)
}

func (h *handlerProcessRuleCategoryImpl) handlerGetAllProcessRulesWithOrderProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	processRules, err := h.s.GetAllProcessRulesWithOrderProcess(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, processRules)
}
