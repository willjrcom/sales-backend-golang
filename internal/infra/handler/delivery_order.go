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
	pcs *deliveryorderusecases.Service
}

func NewHandlerDeliveryOrder(orderService *deliveryorderusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerDeliveryOrderImpl{
		pcs: orderService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterDeliveryOrder)
		c.Get("/{id}", h.handlerGetDeliveryById)
		c.Get("/all", h.handlerGetAllDeliveries)
	})

	return handler.NewHandler("/delivery-order", c)
}

func (h *handlerDeliveryOrderImpl) handlerRegisterDeliveryOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	delivery := &orderdto.CreateDeliveryOrderInput{}
	jsonpkg.ParseBody(r, delivery)

	if id, err := h.pcs.CreateDeliveryOrder(ctx, delivery); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}
func (h *handlerDeliveryOrderImpl) handlerGetDeliveryById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if id, err := h.pcs.GetDeliveryById(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerDeliveryOrderImpl) handlerGetAllDeliveries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if orders, err := h.pcs.GetAllDeliveries(ctx); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: orders})
	}
}
