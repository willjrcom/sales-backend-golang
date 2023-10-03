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
	pcs *orderusecases.Service
}

func NewHandlerOrder(orderService *orderusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerOrderImpl{
		pcs: orderService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterOrder)
		c.Get("/{id}", h.handlerGetOrderById)
		c.Get("/all", h.handlerGetAllOrders)
	})

	c.With().Group(func(c chi.Router) {
		c.Post("/delivery/new/{id}", h.handlerRegisterDeliveryOrder)
	})

	c.With().Group(func(c chi.Router) {
		c.Post("/table/new", h.handlerRegisterTableOrder)
	})

	return handler.NewHandler("/order", c)
}

func (h *handlerOrderImpl) handlerRegisterOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if id, err := h.pcs.CreateDefaultOrder(ctx); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerOrderImpl) handlerGetOrderById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if id, err := h.pcs.GetOrderById(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerOrderImpl) handlerGetAllOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if orders, err := h.pcs.GetAllOrders(ctx); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: orders})
	}
}

func (h *handlerOrderImpl) handlerRegisterDeliveryOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	delivery := &orderdto.CreateDeliveryOrderInput{}
	jsonpkg.ParseBody(r, delivery)

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if id, err := h.pcs.CreateDeliveryOrder(ctx, dtoId, delivery); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerOrderImpl) handlerRegisterTableOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	table := &orderdto.CreateTableOrderInput{}
	jsonpkg.ParseBody(r, table)

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if id, err := h.pcs.CreateTableOrder(ctx, dtoId, table); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}
