package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
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
		c.Get("/all", h.handlerGetAllProducts)
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
	//ctx := r.Context()
	w.Write([]byte("update product"))
}

func (h *handlerProductImpl) handlerDeleteProduct(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	w.Write([]byte("delete product"))
}

func (h *handlerProductImpl) handlerGetProduct(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	w.Write([]byte("get product"))
}

func (h *handlerProductImpl) handlerGetAllProducts(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	w.Write([]byte("get all product"))
}
