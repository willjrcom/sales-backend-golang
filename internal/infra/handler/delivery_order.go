package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
	deliveryorderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/delivery_order"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerDeliveryOrderImpl struct {
	s *deliveryorderusecases.Service
}

func NewHandlerDeliveryOrder(orderService *deliveryorderusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerDeliveryOrderImpl{
		s: orderService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterDeliveryOrder)
		c.Get("/{id}", h.handlerGetDeliveryById)
		c.Get("/all", h.handlerGetAllDeliveries)
		c.Put("/update/{id}/driver", h.handlerUpdateDriver)
		c.Put("/update/{id}/address", h.handlerUpdateDeliveryAddress)
	})

	return handler.NewHandler("/delivery-order", c)
}

func (h *handlerDeliveryOrderImpl) handlerRegisterDeliveryOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	delivery := &orderdto.CreateDeliveryOrderInput{}
	jsonpkg.ParseBody(r, delivery)

	if id, err := h.s.CreateDeliveryOrder(ctx, delivery); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerDeliveryOrderImpl) handlerGetDeliveryById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if id, err := h.s.GetDeliveryById(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerDeliveryOrderImpl) handlerGetAllDeliveries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if orders, err := h.s.GetAllDeliveries(ctx); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: orders})
	}
}

func (h *handlerDeliveryOrderImpl) handlerUpdateDeliveryAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	delivery := &orderdto.UpdateDeliveryOrder{}
	jsonpkg.ParseBody(r, delivery)

	if err := h.s.UpdateDeliveryAddress(ctx, dtoId, delivery); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}

func (h *handlerDeliveryOrderImpl) handlerUpdateDriver(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	delivery := &orderdto.UpdateDriverOrder{}
	jsonpkg.ParseBody(r, delivery)

	if err := h.s.UpdateDriver(ctx, dtoId, delivery); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}
