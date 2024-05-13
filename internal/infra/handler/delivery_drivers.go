package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	deliverydriverdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/delivery_driver"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	deliverydriverusecases "github.com/willjrcom/sales-backend-go/internal/usecases/delivery_driver"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerDeliveryDriverImpl struct {
	s *deliverydriverusecases.Service
}

func NewHandlerDeliveryDriver(sizeService *deliverydriverusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerDeliveryDriverImpl{
		s: sizeService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerCreateDeliveryDriver)
		c.Patch("/update/{id}", h.handlerUpdateDeliveryDriver)
		c.Delete("/{id}", h.handlerDeleteDeliveryDriver)
		c.Get("/{id}", h.handlerGetDeliveryDriver)
		c.Get("/all", h.handlerGetAllDeliveryDrivers)
	})

	return handler.NewHandler("/delivery-driver", c)
}

func (h *handlerDeliveryDriverImpl) handlerCreateDeliveryDriver(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoDeliveryDriver := &deliverydriverdto.CreateDeliveryDriverInput{}
	if err := jsonpkg.ParseBody(r, dtoDeliveryDriver); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	id, err := h.s.CreateDeliveryDriver(ctx, dtoDeliveryDriver)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerDeliveryDriverImpl) handlerUpdateDeliveryDriver(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	dtoDeliveryDriver := &deliverydriverdto.UpdateDeliveryDriverInput{}
	if err := jsonpkg.ParseBody(r, dtoDeliveryDriver); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.s.UpdateDeliveryDriver(ctx, dtoId, dtoDeliveryDriver); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerDeliveryDriverImpl) handlerDeleteDeliveryDriver(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteDeliveryDriver(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerDeliveryDriverImpl) handlerGetDeliveryDriver(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	deliveryDriver, err := h.s.GetDeliveryDriverByID(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: deliveryDriver})
}

func (h *handlerDeliveryDriverImpl) handlerGetAllDeliveryDrivers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	deliveryDrivers, err := h.s.GetAllDeliveryDrivers(ctx)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: deliveryDrivers})
}
