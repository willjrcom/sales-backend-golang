package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderpickupdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_pickup"
	orderpickupusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_pickup"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerOrderPickupImpl struct {
	orderpickupusecases.IService
}

func NewHandlerOrderPickup(orderService orderpickupusecases.IService) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerOrderPickupImpl{orderService}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterOrderPickup)
		c.Get("/{id}", h.handlerGetPickupById)
		c.Get("/all", h.handlerGetAllPickups)
		c.Post("/update/pend/{id}", h.handlerPendingOrder)
		c.Post("/update/ready/{id}", h.handlerReadyOrder)
	})

	return handler.NewHandler("/pickup-order", c)
}

func (h *handlerOrderPickupImpl) handlerRegisterOrderPickup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoPickup := &orderpickupdto.CreateOrderPickupInput{}
	if err := jsonpkg.ParseBody(r, dtoPickup); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	id, err := h.IService.CreateOrderPickup(ctx, dtoPickup)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerOrderPickupImpl) handlerGetPickupById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	Pickup, err := h.IService.GetPickupById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: Pickup})
}

func (h *handlerOrderPickupImpl) handlerGetAllPickups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orders, err := h.IService.GetAllPickups(ctx)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: orders})
}

func (h *handlerOrderPickupImpl) handlerPendingOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.IService.PendingOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderPickupImpl) handlerReadyOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.IService.ReadyOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
