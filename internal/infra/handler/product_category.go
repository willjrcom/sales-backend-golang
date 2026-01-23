package handlerimpl

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
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
		c.Get("/all-map", h.handlerGetAllCategoriesMap)
		c.Get("/all-with-order-process", h.handlerGetAllCategoriesWithProcessRulesAndOrderProcess)
		c.Get("/{id}/complement-products", h.handlerGetComplementProducts)
		c.Get("/{id}/additional-products", h.handlerGetAdditionalProducts)
		c.Get("/{id}/default-products", h.handlerGetDefaultProducts)
		c.Get("/all/additionals", h.handlerGetAdditionalCategories)
		c.Get("/all/complements", h.handlerGetComplementCategories)
		c.Get("/all/default", h.handlerGetDefaultCategories)
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

	dtoIDs := &entitydto.IDRequest{}
	if err := jsonpkg.ParseBody(r, dtoIDs); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

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

	categories, total, err := h.s.GetAllCategories(ctx, dtoIDs, page, perPage, isActive)

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	jsonpkg.ResponseJson(w, r, http.StatusOK, categories)
}

func (h *handlerProductCategoryImpl) handlerGetAllCategoriesMap(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse is_active query parameter (default: true)
	isActive := true
	if isActiveParam := r.URL.Query().Get("is_active"); isActiveParam != "" {
		if val, err := strconv.ParseBool(isActiveParam); err == nil {
			isActive = val
		}
	}

	var isAdditional *bool
	if isAdditionalParam := r.URL.Query().Get("is_additional"); isAdditionalParam != "" {
		if val, err := strconv.ParseBool(isAdditionalParam); err == nil {
			isAdditional = &val
		}
	}

	var isComplement *bool
	if isComplementParam := r.URL.Query().Get("is_complement"); isComplementParam != "" {
		if val, err := strconv.ParseBool(isComplementParam); err == nil {
			isComplement = &val
		}
	}

	categories, err := h.s.GetAllCategoriesMap(ctx, isActive, isAdditional, isComplement)

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

func (h *handlerProductCategoryImpl) handlerGetComplementProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	products, err := h.s.GetComplementProducts(ctx, id)

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, products)
}

func (h *handlerProductCategoryImpl) handlerGetAdditionalProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	products, err := h.s.GetAdditionalProducts(ctx, id)

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, products)
}

func (h *handlerProductCategoryImpl) handlerGetDefaultProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	// Parse is_map query parameter (default: false)
	isMap := false
	if isMapParam := r.URL.Query().Get("is_map"); isMapParam != "" {
		var err error
		isMap, err = strconv.ParseBool(isMapParam)
		if err != nil {
			jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("invalid is_map parameter"))
			return
		}
	}

	if isMap {
		products, err := h.s.GetDefaultProductsMap(ctx, id)
		if err != nil {
			jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
			return
		}
		jsonpkg.ResponseJson(w, r, http.StatusOK, products)
	} else {
		products, err := h.s.GetDefaultProducts(ctx, id)
		if err != nil {
			jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
			return
		}
		jsonpkg.ResponseJson(w, r, http.StatusOK, products)
	}
}

func (h *handlerProductCategoryImpl) handlerGetComplementCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := h.s.GetComplementCategories(ctx)

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, categories)
}

func (h *handlerProductCategoryImpl) handlerGetAdditionalCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := h.s.GetAdditionalCategories(ctx)

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, categories)
}

func (h *handlerProductCategoryImpl) handlerGetDefaultCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := h.s.GetDefaultCategories(ctx)

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, categories)
}
