package handlerimpl

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	companycategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company_category"
	companycategoryusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company_category"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerCompanyCategoryImpl struct {
	s *companycategoryusecases.Service
}

func NewHandlerCompanyCategory(categoryService *companycategoryusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerCompanyCategoryImpl{
		s: categoryService,
	}

	route := "/company-category"

	c.With().Group(func(c chi.Router) {
		c.Post("/", h.handlerCreateCompanyCategory)
		c.Put("/{id}", h.handlerUpdateCompanyCategory)
		c.Delete("/{id}", h.handlerDeleteCompanyCategory)
		c.Get("/{id}", h.handlerGetCompanyCategory)
		c.Get("/", h.handlerGetAllCompanyCategories)
	})

	return handler.NewHandler(route, c)
}

func (h *handlerCompanyCategoryImpl) handlerCreateCompanyCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto := &companycategorydto.CreateCategoryDTO{}
	if err := jsonpkg.ParseBody(r, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	category, err := h.s.CreateCategory(ctx, dto.ToEntity())
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	response := companycategorydto.CompanyCategoryDTO{}
	response.FromDomain(category)
	jsonpkg.ResponseJson(w, r, http.StatusCreated, response)
}

func (h *handlerCompanyCategoryImpl) handlerUpdateCompanyCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	categoryID, err := uuid.Parse(id)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("invalid id format"))
		return
	}

	dto := &companycategorydto.UpdateCategoryDTO{}
	if err := jsonpkg.ParseBody(r, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	category, err := h.s.UpdateCategory(ctx, categoryID, dto.ToEntity())
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	response := companycategorydto.CompanyCategoryDTO{}
	response.FromDomain(category)
	jsonpkg.ResponseJson(w, r, http.StatusOK, response)
}

func (h *handlerCompanyCategoryImpl) handlerDeleteCompanyCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	categoryID, err := uuid.Parse(id)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("invalid id format"))
		return
	}

	if err := h.s.DeleteCategory(ctx, categoryID); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerCompanyCategoryImpl) handlerGetCompanyCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	categoryID, err := uuid.Parse(id)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("invalid id format"))
		return
	}

	category, err := h.s.GetCategory(ctx, categoryID)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	response := companycategorydto.CompanyCategoryDTO{}
	response.FromDomain(category)
	jsonpkg.ResponseJson(w, r, http.StatusOK, response)
}

func (h *handlerCompanyCategoryImpl) handlerGetAllCompanyCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := h.s.GetAllCompanyCategories(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	responses := make([]companycategorydto.CompanyCategoryDTO, len(categories))
	for i, cat := range categories {
		responses[i].FromDomain(&cat)
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, responses)
}
