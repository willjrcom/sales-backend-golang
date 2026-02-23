package handlerimpl

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	itemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/item"
	orderusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerItemImpl struct {
	s *orderusecases.ItemService
}

func NewHandlerItem(itemService *orderusecases.ItemService) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerItemImpl{
		s: itemService,
	}

	route := "/item"

	c.With().Group(func(c chi.Router) {
		c.Post("/add", h.handlerAddItem)
		c.Delete("/{id}", h.handlerDeleteItem)
		c.Post("/update/{id}/additional", h.handlerAddAdditionalItem)
		c.Delete("/delete/{id-additional}/additional", h.handlerDeleteAdditionalItem)
		c.Post("/update/{id}/removed-item", h.handlerAddRemovedItem)
		c.Delete("/delete/{id}/removed-item", h.handlerRemoveRemovedItem)
	})

	unprotectedRoutes := []string{}
	return handler.NewHandler(route, c, unprotectedRoutes...)
}

func (h *handlerItemImpl) handlerAddItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoAddItem := &itemdto.OrderItemCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoAddItem); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	ids, err := h.s.AddItemOrder(ctx, dtoAddItem)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, ids)
}

func (h *handlerItemImpl) handlerDeleteItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	groupItemDeleted, err := h.s.DeleteItemOrder(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, groupItemDeleted)
}

func (h *handlerItemImpl) handlerAddAdditionalItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoAddAdditionalItem := &itemdto.OrderAdditionalItemCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoAddAdditionalItem); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	additionalId, err := h.s.AddAdditionalItemOrder(ctx, dtoId, dtoAddAdditionalItem)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, additionalId)
}

func (h *handlerItemImpl) handlerDeleteAdditionalItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id-additional")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteAdditionalItemOrder(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerItemImpl) handlerAddRemovedItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoRemovedItem := &itemdto.RemovedItemDTO{}
	if err := jsonpkg.ParseBody(r, dtoRemovedItem); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.AddRemovedItem(ctx, dtoId, dtoRemovedItem); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerItemImpl) handlerRemoveRemovedItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoRemovedItem := &itemdto.RemovedItemDTO{}
	if err := jsonpkg.ParseBody(r, dtoRemovedItem); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.RemoveRemovedItem(ctx, dtoId, dtoRemovedItem); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
