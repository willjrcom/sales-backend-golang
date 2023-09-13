package handlerimpl

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	filterdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/filter"
	productdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product"
	productusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerProductImpl struct {
	ps *productusecases.Service
}

func NewHandlerProduct(productService *productusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerProductImpl{
		ps: productService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterProduct)
		c.Put("/update/{id}", h.handlerUpdateProduct)
		c.Delete("/delete/{id}", h.handlerDeleteProduct)
		c.Get("/{id}", h.handlerGetProduct)
		c.Post("/by", h.handlerGetProductBy)
		c.Post("/all", h.handlerGetAllProducts)
	})

	return handler.NewHandler("/product", c)
}

func (h *handlerProductImpl) handlerRegisterProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	product := &productdto.RegisterProductInput{}
	jsonpkg.ParseBody(r, product)

	id, err := h.ps.RegisterProduct(ctx, product)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("new product: " + id.String()))
}

func (h *handlerProductImpl) handlerUpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	product := &productdto.UpdateProductInput{}
	jsonpkg.ParseBody(r, product)

	err := h.ps.UpdateProduct(ctx, dtoId, product)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("new product:"))
}

func (h *handlerProductImpl) handlerDeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	err := h.ps.DeleteProductById(ctx, dtoId)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("delete product"))
}

func (h *handlerProductImpl) handlerGetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	product, err := h.ps.GetProductById(ctx, dtoId)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	text, err := json.Marshal(product)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(text)
}

func (h *handlerProductImpl) handlerGetProductBy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter := &productdto.FilterProductInput{}
	jsonpkg.ParseBody(r, filter)

	product, err := h.ps.GetProductBy(ctx, filter)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	text, err := json.Marshal(product)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(text)
}

func (h *handlerProductImpl) handlerGetAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	filter := &filterdto.Category{}
	jsonpkg.ParseBody(r, filter)

	categories, err := h.ps.GetAllProducts(ctx, filter)

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
