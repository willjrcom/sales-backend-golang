package handlerimpl

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderpickupdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_pickup"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerOrderPickupImpl struct {
	orderusecases.IPickupService
}

func NewHandlerOrderPickup(orderService orderusecases.IPickupService) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerOrderPickupImpl{orderService}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerCreateOrderPickup)
		c.Get("/{id}", h.handlerGetPickupById)
		c.Get("/all", h.handlerGetAllPickups)
		c.Post("/update/pend/{id}", h.handlerPendingOrder)
		c.Post("/update/ready/{id}", h.handlerReadyOrder)
		c.Put("/update/name/{id}", h.handlerUpdateName)
	})

	return handler.NewHandler("/order-pickup", c)
}

func (h *handlerOrderPickupImpl) handlerCreateOrderPickup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoPickup := &orderpickupdto.OrderPickupCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoPickup); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.IPickupService.CreateOrderPickup(ctx, dtoPickup)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, id)
}

func (h *handlerOrderPickupImpl) handlerGetPickupById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	Pickup, err := h.IPickupService.GetPickupById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, Pickup)
}

func (h *handlerOrderPickupImpl) handlerGetAllPickups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orders, err := h.IPickupService.GetAllPickups(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, orders)
}

func (h *handlerOrderPickupImpl) handlerPendingOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.IPickupService.PendingOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderPickupImpl) handlerReadyOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.IPickupService.ReadyOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderPickupImpl) handlerUpdateName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoPickup := &orderpickupdto.UpdateOrderPickupInput{}
	if err := jsonpkg.ParseBody(r, dtoPickup); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.IPickupService.UpdateName(ctx, dtoId, dtoPickup); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
