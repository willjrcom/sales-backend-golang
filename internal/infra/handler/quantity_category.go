package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	quantitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/quantity_category"
	quantityusecases "github.com/willjrcom/sales-backend-go/internal/usecases/quantity_category"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerQuantityCategoryImpl struct {
	s *quantityusecases.Service
}

func NewHandlerQuantityCategory(quantityService *quantityusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerQuantityCategoryImpl{
		s: quantityService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterQuantity)
		c.Patch("/update/{id}", h.handlerUpdateQuantity)
		c.Delete("/{id}", h.handlerDeleteQuantity)
	})

	return handler.NewHandler("/category-product/quantity", c)
}

func (h *handlerQuantityCategoryImpl) handlerRegisterQuantity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	quantity := &quantitydto.RegisterQuantityInput{}
	jsonpkg.ParseBody(r, quantity)

	id, err := h.s.RegisterQuantity(ctx, quantity)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerQuantityCategoryImpl) handlerUpdateQuantity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	Quantity := &quantitydto.UpdateQuantityInput{}
	jsonpkg.ParseBody(r, Quantity)

	if err := h.s.UpdateQuantity(ctx, dtoId, Quantity); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerQuantityCategoryImpl) handlerDeleteQuantity(w http.ResponseWriter, r *http.Request) {
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
