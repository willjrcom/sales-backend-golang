package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	deliveryorderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/delivery"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	deliveryorderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/delivery_order"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerDeliveryOrderImpl struct {
	deliveryorderusecases.IService
}

func NewHandlerDeliveryOrder(orderService deliveryorderusecases.IService) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerDeliveryOrderImpl{orderService}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterDeliveryOrder)
		c.Get("/{id}", h.handlerGetDeliveryById)
		c.Get("/all", h.handlerGetAllDeliveries)
		c.Post("/update/pend/{id}", h.handlerPendDeliveryOrder)
		c.Post("/update/ship/{id}", h.handlerShipDeliveryOrder)
		c.Post("/update/delivery/{id}", h.handlerDeliveryOrder)
		c.Put("/update/driver/{id}", h.handlerUpdateDriver)
		c.Put("/update/address/{id}", h.handlerUpdateDeliveryAddress)
	})

	return handler.NewHandler("/delivery-order", c)
}

func (h *handlerDeliveryOrderImpl) handlerRegisterDeliveryOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoDelivery := &deliveryorderdto.CreateDeliveryOrderInput{}
	if err := jsonpkg.ParseBody(r, dtoDelivery); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	id, err := h.IService.CreateDeliveryOrder(ctx, dtoDelivery)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerDeliveryOrderImpl) handlerGetDeliveryById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	delivery, err := h.IService.GetDeliveryById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: delivery})
}

func (h *handlerDeliveryOrderImpl) handlerGetAllDeliveries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orders, err := h.IService.GetAllDeliveries(ctx)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: orders})
}

func (h *handlerDeliveryOrderImpl) handlerPendDeliveryOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.IService.PendDeliveryOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerDeliveryOrderImpl) handlerShipDeliveryOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	dtoDelivery := &deliveryorderdto.UpdateDriverOrder{}
	if err := jsonpkg.ParseBody(r, dtoDelivery); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.IService.ShipDeliveryOrder(ctx, dtoId, dtoDelivery); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerDeliveryOrderImpl) handlerDeliveryOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.IService.DeliveryOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerDeliveryOrderImpl) handlerUpdateDeliveryAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.IService.UpdateDeliveryAddress(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerDeliveryOrderImpl) handlerUpdateDriver(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	dtoDelivery := &deliveryorderdto.UpdateDriverOrder{}
	if err := jsonpkg.ParseBody(r, dtoDelivery); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.IService.UpdateDeliveryDriver(ctx, dtoId, dtoDelivery); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
