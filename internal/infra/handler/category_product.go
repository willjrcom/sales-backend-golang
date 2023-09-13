package handlerimpl

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product"
	categoryproductusecases "github.com/willjrcom/sales-backend-go/internal/usecases/category_product"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerCategoryProductImpl struct {
	pcs *categoryproductusecases.Service
}

func NewHandlerCategoryProduct(categoryService *categoryproductusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerCategoryProductImpl{
		pcs: categoryService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterCategoryProduct)
		c.Put("/update/{id}", h.handlerUpdateCategoryProduct)
		c.Delete("/delete/{id}", h.handlerDeleteCategoryProduct)
		c.Get("/{id}", h.handlerGetCategoryProduct)
		c.Get("/allproducts", h.handlerGetAllCategoryProducts)
		c.Get("/allsizes", h.handlerGetAllCategorySizes)
	})

	return handler.NewHandler("/category", c)
}

func (h *handlerCategoryProductImpl) handlerRegisterCategoryProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	category := &productdto.RegisterCategoryInput{}
	jsonpkg.ParseBody(r, category)

	id, err := h.pcs.RegisterCategory(ctx, category)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("new product: " + id.String()))
}

func (h *handlerCategoryProductImpl) handlerUpdateCategoryProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	category := &productdto.UpdateCategoryInput{}
	jsonpkg.ParseBody(r, category)

	err := h.pcs.UpdateCategory(ctx, dtoId, category)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("update product name"))
}

func (h *handlerCategoryProductImpl) handlerDeleteCategoryProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	err := h.pcs.DeleteCategoryById(ctx, dtoId)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("delete category"))
}

func (h *handlerCategoryProductImpl) handlerGetCategoryProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	category, err := h.pcs.GetCategoryById(ctx, dtoId)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	text, err := json.Marshal(category)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(text)
}

func (h *handlerCategoryProductImpl) handlerGetAllCategoryProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categories, err := h.pcs.GetAllCategoryProducts(ctx)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	text, err := json.Marshal(categories)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(text)
}

func (h *handlerCategoryProductImpl) handlerGetAllCategorySizes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	categories, err := h.pcs.GetAllCategorySizes(ctx)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	text, err := json.Marshal(categories)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(text)
}
