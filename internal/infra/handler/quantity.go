package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	quantitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/quantity"
	quantityusecases "github.com/willjrcom/sales-backend-go/internal/usecases/quantity"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerQuantityImpl struct {
	s *quantityusecases.Service
}

func NewHandlerQuantity(quantityService *quantityusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerQuantityImpl{
		s: quantityService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerCreateQuantity)
		c.Patch("/update/{id}", h.handlerUpdateQuantity)
		c.Delete("/{id}", h.handlerDeleteQuantity)
	})

	return handler.NewHandler("/product-category/quantity", c)
}

func (h *handlerQuantityImpl) handlerCreateQuantity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoQuantity := &quantitydto.CreateQuantityInput{}
	if err := jsonpkg.ParseBody(r, dtoQuantity); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	id, err := h.s.CreateQuantity(ctx, dtoQuantity)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerQuantityImpl) handlerUpdateQuantity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	dtoQuantity := &quantitydto.UpdateQuantityInput{}
	if err := jsonpkg.ParseBody(r, dtoQuantity); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.s.UpdateQuantity(ctx, dtoId, dtoQuantity); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerQuantityImpl) handlerDeleteQuantity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteQuantity(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
