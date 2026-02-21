package handlerimpl

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order"
	ordertabledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_table"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerOrderImpl struct {
	s *orderusecases.OrderService
}

func NewHandlerOrder(orderService *orderusecases.OrderService) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerOrderImpl{
		s: orderService,
	}

	c.With().Group(func(c chi.Router) {
		c.Get("/{id}", h.handlerGetOrderById)
		c.Get("/all", h.handlerGetAllOrders)
		c.Get("/all/delivery", h.GetAllOrdersWithDelivery)
		c.Get("/all/pickup/ready", h.GetAllOrdersWithPickupReady)
		c.Get("/all/pickup/delivered", h.GetAllOrdersWithPickupDelivered)
		c.Get("/all/pickup/by-contact/{contact}", h.GetAllOrdersWithPickupByContact)
		c.Get("/all/delivery/by-client/{id}", h.handlerGetAllOrdersByClientID)
		c.Get("/all/table/by-table/{id}", h.handlerGetAllOrdersByTable)
		c.Put("/update/{id}/observation", h.handlerUpdateObservation)
		c.Put("/update/{id}/payment", h.handlerUpdatePaymentMethod)
		c.Post("/pending/{id}", h.handlerPendingOrder)
		c.Post("/ready/{id}", h.handlerReadyOrder)
		c.Post("/finish/{id}", h.handlerFinishOrder)
		c.Post("/cancel/{id}", h.handlerCancelOrder)
		c.Post("/archive/{id}", h.handlerArchiveOrder)
		c.Post("/unarchive/{id}", h.handlerUnarchiveOrder)
		c.Delete("/{id}", h.handlerDeleteOrder)
	})

	return handler.NewHandler("/order", c)
}

func (h *handlerOrderImpl) handlerGetOrderById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	order, err := h.s.GetOrderById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, order)
}

func (h *handlerOrderImpl) handlerGetAllOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orders, err := h.s.GetAllOrders(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, orders)
}

func (h *handlerOrderImpl) GetAllOrdersWithDelivery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page, perPage := headerservice.GetPageAndPerPage(r, 0, 20)

	orders, err := h.s.GetAllOrdersWithDelivery(ctx, page, perPage)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, orders)
}

func (h *handlerOrderImpl) GetAllOrdersWithPickupReady(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page, perPage := headerservice.GetPageAndPerPage(r, 0, 100)

	orders, err := h.s.GetAllOrdersWithPickup(ctx, orderentity.OrderPickupStatusReady, page, perPage)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, orders)
}

func (h *handlerOrderImpl) handlerGetAllOrdersByClientID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	orders, err := h.s.GetOrderIDFromOrderDeliveriesByClientId(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	if len(orders) == 0 {
		jsonpkg.ResponseJson(w, r, http.StatusOK, orders)
		return
	}

	delivery := orders[0].Delivery
	var contact string
	if delivery != nil && delivery.Client != nil && delivery.Client.Contact != nil {
		contact = delivery.Client.Contact.Number
	}

	if contact == "" {
		jsonpkg.ResponseJson(w, r, http.StatusOK, orders)
		return
	}

	pickupOrders, err := h.s.GetOrdersPickupByContact(ctx, contact)
	if len(pickupOrders) > 0 {
		orders = append(orders, pickupOrders...)
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, orders)
}

func (h *handlerOrderImpl) handlerGetAllOrdersByTable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	id := chi.URLParam(r, "id")
	contact := r.URL.Query().Get("contact")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoContact := &ordertabledto.OrderTableContactInput{
		TableID: uuid.MustParse(id),
		Contact: contact,
	}

	orders, err := h.s.GetOrdersTableByTableId(ctx, dtoContact)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, orders)
}

func (h *handlerOrderImpl) GetAllOrdersWithPickupDelivered(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page, perPage := headerservice.GetPageAndPerPage(r, 0, 10)

	orders, err := h.s.GetAllOrdersWithPickup(ctx, orderentity.OrderPickupStatusDelivered, page, perPage)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, orders)
}

func (h *handlerOrderImpl) GetAllOrdersWithPickupByContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contact := chi.URLParam(r, "contact")

	if contact == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("contact is required"))
		return
	}

	orders, err := h.s.GetOrdersPickupByContact(ctx, contact)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, orders)
}

func (h *handlerOrderImpl) handlerUpdateObservation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoObservation := &orderdto.OrderUpdateObservationDTO{}
	if err := jsonpkg.ParseBody(r, dtoObservation); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdateOrderObservation(ctx, dtoId, dtoObservation); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderImpl) handlerUpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoPayment := &orderdto.OrderPaymentCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoPayment); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.AddPayment(ctx, dtoId, dtoPayment); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderImpl) handlerPendingOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.PendingOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderImpl) handlerReadyOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.ReadyOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderImpl) handlerFinishOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.FinishOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderImpl) handlerCancelOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.CancelOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderImpl) handlerArchiveOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.ArchiveOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderImpl) handlerUnarchiveOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.UnarchiveOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerOrderImpl) handlerDeleteOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteOrderByID(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
