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

	route := "/item"

	c.With().Group(func(c chi.Router) {
		c.Post("/add", h.handlerAddItem)
		c.Delete("/{id}", h.handlerDeleteItem)
		c.Post("/update/{id}/additional", h.handlerAddAdditionalItem)
		c.Delete("/delete/{id-additional}/additional", h.handlerDeleteAdditionalItem)
	})

	unprotectedRoutes := []string{}
	return handler.NewHandler(route, c, unprotectedRoutes...)
}

func (h *handlerItemImpl) handlerAddItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoAddItem := &itemdto.AddItemOrderInput{}
	if err := jsonpkg.ParseBody(r, dtoAddItem); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	ids, err := h.s.AddItemOrder(ctx, dtoAddItem)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: ids})
}

func (h *handlerItemImpl) handlerDeleteItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteItemOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerItemImpl) handlerAddAdditionalItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	dtoAddAdditionalItem := &itemdto.AddAdditionalItemOrderInput{}
	if err := jsonpkg.ParseBody(r, dtoAddAdditionalItem); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	additionalId, err := h.s.AddAdditionalItemOrder(ctx, dtoId, dtoAddAdditionalItem)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: additionalId})
}

func (h *handlerItemImpl) handlerDeleteAdditionalItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id-additional")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteAdditionalItemOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
