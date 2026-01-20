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
	productusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type HandlerProductImpl struct {
	s *productusecases.Service
}

func NewHandlerProduct(productService *productusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &HandlerProductImpl{
		s: productService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerCreateProduct)
		c.Get("/all", h.handlerGetAllProducts)
		c.Get("/all/default", h.handlerGetDefaultProducts)
		c.Get("/all/by-category-id/{category_id}", h.handlerGetProductsByCategoryId)
		c.Get("/all-map", h.handlerGetAllProductsMap)
		c.Get("/code/{code}", h.handlerGetProductByCode)
		c.Patch("/update/{id}", h.handlerUpdateProduct)
		c.Delete("/{id}", h.handlerDeleteProduct)
		c.Get("/{id}", h.handlerGetProduct)
	})

	return handler.NewHandler("/product", c)
}

func (h *HandlerProductImpl) handlerCreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoProduct := &productcategorydto.ProductCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoProduct); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.s.CreateProduct(ctx, dtoProduct)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, id)
}

func (h *HandlerProductImpl) handlerUpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoProduct := &productcategorydto.ProductUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dtoProduct); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdateProduct(ctx, dtoId, dtoProduct); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *HandlerProductImpl) handlerDeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteProductById(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *HandlerProductImpl) handlerGetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	product, err := h.s.GetProductById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, product)
}

func (h *HandlerProductImpl) handlerGetProductByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")

	if code == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dto := &productcategorydto.Keys{Code: code}

	product, err := h.s.GetProductByCode(ctx, dto)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, product)
}

func (h *HandlerProductImpl) handlerGetAllProducts(w http.ResponseWriter, r *http.Request) {
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

	// Parse category_id query parameter (default: true)
	categoryID := r.URL.Query().Get("category_id")

	categories, total, err := h.s.GetAllProducts(ctx, page, perPage, isActive, categoryID)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	jsonpkg.ResponseJson(w, r, http.StatusOK, categories)
}

func (h *HandlerProductImpl) handlerGetDefaultProducts(w http.ResponseWriter, r *http.Request) {
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

	products, total, err := h.s.GetDefaultProducts(ctx, page, perPage, isActive)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	jsonpkg.ResponseJson(w, r, http.StatusOK, products)
}

func (h *HandlerProductImpl) handlerGetProductsByCategoryId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get category_id from URL path parameter
	categoryID := chi.URLParam(r, "category_id")

	if categoryID == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("category_id is required"))
		return
	}

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

	products, err := h.s.GetAllProductsMap(ctx, isActive, categoryID)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, products)
}

func (h *HandlerProductImpl) handlerGetAllProductsMap(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse is_active query parameter (default: true)
	isActive := true
	if isActiveParam := r.URL.Query().Get("is_active"); isActiveParam != "" {
		if val, err := strconv.ParseBool(isActiveParam); err == nil {
			isActive = val
		}
	}

	// Parse category_id query parameter (optional)
	categoryID := r.URL.Query().Get("category_id")

	products, err := h.s.GetAllProductsMap(ctx, isActive, categoryID)

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, products)
}
