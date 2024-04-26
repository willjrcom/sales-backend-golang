package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	pickuporderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/pickup_order"
	pickuporderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/pickup_order"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerPickupOrderImpl struct {
	pickuporderusecases.IService
}

func NewHandlerPickupOrder(orderService pickuporderusecases.IService) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerPickupOrderImpl{orderService}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterPickupOrder)
		c.Get("/{id}", h.handlerGetPickupById)
		c.Get("/all", h.handlerGetAllPickups)
		c.Post("/update/launch/{id}", h.handlerLaunchOrder)
		c.Post("/update/pickup/{id}", h.handlerPickupOrder)
	})

	return handler.NewHandler("/pickup-order", c)
}

func (h *handlerPickupOrderImpl) handlerRegisterPickupOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoPickup := &pickuporderdto.CreatePickupOrderInput{}
	if err := jsonpkg.ParseBody(r, dtoPickup); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	id, err := h.IService.CreatePickupOrder(ctx, dtoPickup)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerPickupOrderImpl) handlerGetPickupById(w http.ResponseWriter, r *http.Request) {
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

func (h *handlerPickupOrderImpl) handlerGetAllPickups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orders, err := h.IService.GetAllPickups(ctx)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: orders})
}

func (h *handlerPickupOrderImpl) handlerLaunchOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.IService.LaunchOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerPickupOrderImpl) handlerPickupOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.IService.PickupOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
