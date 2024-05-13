package handlerimpl

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
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
		c.Post("/new", h.handlerRegisterProductCategory)
		c.Patch("/update/{id}", h.handlerUpdateProductCategory)
		c.Delete("/{id}", h.handlerDeleteProductCategory)
		c.Get("/{id}", h.handlerGetProductCategory)
		c.Get("/all", h.handlerGetAllCategories)
	})

	return handler.NewHandler(route, c)
}

func (h *handlerProductCategoryImpl) handlerRegisterProductCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoCategory := &productcategorydto.RegisterCategoryInput{}
	if err := jsonpkg.ParseBody(r, dtoCategory); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	id, err := h.s.RegisterCategory(ctx, dtoCategory)

	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerProductCategoryImpl) handlerUpdateProductCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	dtoCategory := &productcategorydto.UpdateCategoryInput{}
	if err := jsonpkg.ParseBody(r, dtoCategory); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.s.UpdateCategory(ctx, dtoId, dtoCategory); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerProductCategoryImpl) handlerDeleteProductCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteCategoryById(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerProductCategoryImpl) handlerGetProductCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println(ctx.Value(schemaentity.Schema("schema")))
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	category, err := h.s.GetCategoryById(ctx, dtoId)

	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: category})
}

func (h *handlerProductCategoryImpl) handlerGetAllCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categories, err := h.s.GetAllCategories(ctx)

	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: categories})
}
