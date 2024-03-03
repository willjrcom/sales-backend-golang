package handlerimpl

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
	categorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/category"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	categoryproductusecases "github.com/willjrcom/sales-backend-go/internal/usecases/category_product"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerCategoryProductImpl struct {
	s *categoryproductusecases.Service
}

func NewHandlerCategoryProduct(categoryService *categoryproductusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerCategoryProductImpl{
		s: categoryService,
	}

	route := "/category-product"

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterCategoryProduct)
		c.Patch("/update/{id}", h.handlerUpdateCategoryProduct)
		c.Delete("/delete/{id}", h.handlerDeleteCategoryProduct)
		c.Get("/{id}", h.handlerGetCategoryProduct)
		c.Get("/all", h.handlerGetAllCategories)
	})

	return handler.NewHandler(route, c)
}

func (h *handlerCategoryProductImpl) handlerRegisterCategoryProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	category := &categorydto.RegisterCategoryInput{}
	jsonpkg.ParseBody(r, category)

	id, err := h.s.RegisterCategory(ctx, category)

	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerCategoryProductImpl) handlerUpdateCategoryProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	category := &categorydto.UpdateCategoryInput{}
	jsonpkg.ParseBody(r, category)

	err := h.s.UpdateCategory(ctx, dtoId, category)

	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerCategoryProductImpl) handlerDeleteCategoryProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	err := h.s.DeleteCategoryById(ctx, dtoId)

	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerCategoryProductImpl) handlerGetCategoryProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println(ctx.Value(schemaentity.Schema("schema")))
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	category, err := h.s.GetCategoryById(ctx, dtoId)

	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: category})
}

func (h *handlerCategoryProductImpl) handlerGetAllCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categories, err := h.s.GetAllCategories(ctx)

	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: categories})
}
