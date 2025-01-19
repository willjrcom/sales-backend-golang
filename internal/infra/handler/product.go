package handlerimpl

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
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
		c.Patch("/update/{id}", h.handlerUpdateProduct)
		c.Delete("/{id}", h.handlerDeleteProduct)
		c.Get("/{id}", h.handlerGetProduct)
		c.Get("/code/{code}", h.handlerGetProductByCode)
		c.Get("/all", h.handlerGetAllProducts)
	})

	return handler.NewHandler("/product", c)
}

func (h *HandlerProductImpl) handlerCreateProduct(w http.ResponseWriter, r *http.Request) {
	// file, _, err := r.FormFile("image")
	// if err != nil && err.Error() != "http: no such file" {
	// 	jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
	// 	return
	// }

	// defer file.Close()

	ctx := r.Context()

	dtoProduct := &productcategorydto.ProductCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoProduct); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	// if file != nil {
	// 	product.Image = &file
	// }

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

	categories, err := h.s.GetAllProducts(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, categories)
}
