package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	itemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/item"
	itemusecases "github.com/willjrcom/sales-backend-go/internal/usecases/item"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerItemImpl struct {
	s *itemusecases.Service
}

func NewHandlerItem(itemService *itemusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerItemImpl{
		s: itemService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/add", h.handlerRegisterItem)
	})

	return handler.NewHandler("/item", c)
}

func (h *handlerItemImpl) handlerRegisterItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	addItem := &itemdto.AddItemOrderInput{}
	jsonpkg.ParseBody(r, addItem)

	if id, err := h.s.AddItemOrder(ctx, addItem); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}
