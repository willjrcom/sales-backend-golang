package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
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
		c.Post("/add", h.handlerAddItem)
		c.Delete("/delete/{id}", h.handlerDeleteItem)
		c.Post("/add/aditional/{id}", h.handlerAddAditionalItem)
	})

	return handler.NewHandler("/item", c)
}

func (h *handlerItemImpl) handlerAddItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	addItem := &itemdto.AddItemOrderInput{}
	jsonpkg.ParseBody(r, addItem)

	if ids, err := h.s.AddItemOrder(ctx, addItem); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: ids})
	}
}

func (h *handlerItemImpl) handlerDeleteItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteItemOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerItemImpl) handlerAddAditionalItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if id, err := h.s.AddAditionalItemOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}
