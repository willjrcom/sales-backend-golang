package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerOrderImpl struct {
	s *orderusecases.Service
}

func NewHandlerOrder(orderService *orderusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerOrderImpl{
		s: orderService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterOrder)
		c.Get("/{id}", h.handlerGetOrderById)
		c.Get("/all", h.handlerGetAllOrders)
		c.Put("/update/{id}/observation", h.handlerUpdateObservation)
		c.Put("/update/{id}/payment", h.handlerUpdatePaymentMethod)
	})

	return handler.NewHandler("/order", c)
}

func (h *handlerOrderImpl) handlerRegisterOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if id, err := h.s.CreateDefaultOrder(ctx); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerOrderImpl) handlerGetOrderById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if order, err := h.s.GetOrderById(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: order})
	}
}

func (h *handlerOrderImpl) handlerGetAllOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if orders, err := h.s.GetAllOrders(ctx); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: orders})
	}
}

func (h *handlerOrderImpl) handlerUpdateObservation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	observation := &orderdto.UpdateObservationOrder{}
	jsonpkg.ParseBody(r, observation)

	if err := h.s.UpdateOrderObservation(ctx, dtoId, observation); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}

func (h *handlerOrderImpl) handlerUpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	payment := &orderdto.UpdatePaymentMethod{}
	jsonpkg.ParseBody(r, payment)

	if err := h.s.UpdatePaymentMethod(ctx, dtoId, payment); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}
