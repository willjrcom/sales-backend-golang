package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productcategoryproductdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category_product"
	productcategoryproductusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product_category_product"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerProductImpl struct {
	s *productcategoryproductusecases.Service
}

func NewHandlerProduct(productService *productcategoryproductusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerProductImpl{
		s: productService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterProduct)
		c.Patch("/update/{id}", h.handlerUpdateProduct)
		c.Delete("/{id}", h.handlerDeleteProduct)
		c.Get("/{id}", h.handlerGetProduct)
		c.Get("/code/{code}", h.handlerGetProductByCode)
		c.Get("/all", h.handlerGetAllProducts)
	})

	return handler.NewHandler("/product", c)
}

func (h *handlerProductImpl) handlerRegisterProduct(w http.ResponseWriter, r *http.Request) {
	// file, _, err := r.FormFile("image")
	// if err != nil && err.Error() != "http: no such file" {
	// 	jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
	// 	return
	// }

	// defer file.Close()

	ctx := r.Context()

	dtoProduct := &productcategoryproductdto.RegisterProductInput{}
	if err := jsonpkg.ParseBody(r, dtoProduct); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	// if file != nil {
	// 	product.Image = &file
	// }

	id, err := h.s.RegisterProduct(ctx, dtoProduct)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerProductImpl) handlerUpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	dtoProduct := &productcategoryproductdto.UpdateProductInput{}
	if err := jsonpkg.ParseBody(r, dtoProduct); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.s.UpdateProduct(ctx, dtoId, dtoProduct); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerProductImpl) handlerDeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteProductById(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerProductImpl) handlerGetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	product, err := h.s.GetProductById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: product})
}

func (h *handlerProductImpl) handlerGetProductByCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	code := chi.URLParam(r, "code")

	if code == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dto := &productcategoryproductdto.Keys{Code: code}

	product, err := h.s.GetProductByCode(ctx, dto)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: product})
}

func (h *handlerProductImpl) handlerGetAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := h.s.GetAllProducts(ctx)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: categories})
}
