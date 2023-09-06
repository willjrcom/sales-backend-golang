package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	productusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product"
)

type handlerProductImpl struct {
	productService *productusecases.Service
}

func NewHandlerProduct(productService *productusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerProductImpl{
		productService: productService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterProduct)
	})

	return handler.NewHandler("/product", c)
}

func (h *handlerProductImpl) handlerRegisterProduct(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	w.Write([]byte("new product"))
}

func (h *handlerProductImpl) handlerUpdateProduct(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	w.Write([]byte("update product"))
}

func (h *handlerProductImpl) handlerDeleteProduct(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	w.Write([]byte("update product"))
}
