package handlerimpl

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdeliverydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_delivery"
	orderdeliveryusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_delivery"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerOrderDeliveryImpl struct {
	orderdeliveryusecases.IService
}

func NewHandlerOrderDelivery(orderService orderdeliveryusecases.IService) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerOrderDeliveryImpl{orderService}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerCreateOrderDelivery)
		c.Get("/{id}", h.handlerGetDeliveryById)
		c.Get("/all", h.handlerGetAllDeliveries)
		c.Post("/update/pend/{id}", h.handlerPendOrderDelivery)
		c.Post("/update/ship", h.handlerShipOrderDelivery)
		c.Post("/update/delivery/{id}", h.handlerDeliveryOrderDelivery)
		c.Put("/update/driver/{id}", h.handlerUpdateDriver)
		c.Put("/update/change/{id}", h.handlerUpdateChange)
		c.Put("/update/address/{id}", h.handlerUpdateDeliveryAddress)
	})

	return handler.NewHandler("/order-delivery", c)
}

func (h *handlerOrderDeliveryImpl) handlerCreateOrderDelivery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoDelivery := &orderdeliverydto.DeliveryOrderCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoDelivery); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.IService.CreateOrderDelivery(ctx, dtoDelivery)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, id)
}

func (h *handlerOrderDeliveryImpl) handlerGetDeliveryById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	delivery, err := h.IService.GetDeliveryById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, delivery)
}

func (h *handlerOrderDeliveryImpl) handlerGetAllDeliveries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orders, err := h.IService.GetAllDeliveries(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, orders)
}

func (h *handlerOrderDeliveryImpl) handlerPendOrderDelivery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.IService.PendOrderDelivery(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderDeliveryImpl) handlerShipOrderDelivery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoDelivery := &orderdeliverydto.DeliveryOrderUpdateShipDTO{}
	if err := jsonpkg.ParseBody(r, dtoDelivery); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.IService.ShipOrderDelivery(ctx, dtoDelivery); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderDeliveryImpl) handlerDeliveryOrderDelivery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.IService.OrderDelivery(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderDeliveryImpl) handlerUpdateDeliveryAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.IService.UpdateDeliveryAddress(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderDeliveryImpl) handlerUpdateDriver(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoDelivery := &orderdeliverydto.DeliveryOrderDriverUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dtoDelivery); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.IService.UpdateDeliveryDriver(ctx, dtoId, dtoDelivery); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderDeliveryImpl) handlerUpdateChange(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoDelivery := &orderdeliverydto.OrderChangeCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoDelivery); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.IService.UpdateDeliveryChange(ctx, dtoId, dtoDelivery); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
