package handlerimpl

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
	productcategoryusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product_category"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerProductCategoryImpl struct {
	s *productcategoryusecases.Service
}

func NewHandlerProductCategory(categoryService *productcategoryusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerProductCategoryImpl{
		s: categoryService,
	}

	route := "/product-category"

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerCreateProductCategory)
		c.Patch("/update/{id}", h.handlerUpdateProductCategory)
		c.Delete("/{id}", h.handlerDeleteProductCategory)
		c.Get("/{id}", h.handlerGetProductCategory)
		c.Get("/all", h.handlerGetAllCategories)
		c.Get("/all-with-order-process", h.handlerGetAllCategoriesWithProcessRulesAndOrderProcess)
	})

	return handler.NewHandler(route, c)
}

func (h *handlerProductCategoryImpl) handlerCreateProductCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoCategory := &productcategorydto.CategoryCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoCategory); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.s.CreateCategory(ctx, dtoCategory)

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, id)
}

func (h *handlerProductCategoryImpl) handlerUpdateProductCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoCategory := &productcategorydto.CategoryUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dtoCategory); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdateCategory(ctx, dtoId, dtoCategory); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerProductCategoryImpl) handlerDeleteProductCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteCategoryById(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerProductCategoryImpl) handlerGetProductCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	category, err := h.s.GetCategoryById(ctx, dtoId)

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, category)
}

func (h *handlerProductCategoryImpl) handlerGetAllCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categories, err := h.s.GetAllCategories(ctx)

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, categories)
}

func (h *handlerProductCategoryImpl) handlerGetAllCategoriesWithProcessRulesAndOrderProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	processRules, err := h.s.GetAllCategoriesWithProcessRulesAndOrderProcess(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, processRules)
}
