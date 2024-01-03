package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	groupitemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/group_item"
	groupitemusecases "github.com/willjrcom/sales-backend-go/internal/usecases/group_item"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerGroupItemImpl struct {
	s *groupitemusecases.Service
}

func NewHandlerGroupItem(itemService *groupitemusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerGroupItemImpl{
		s: itemService,
	}

	c.With().Group(func(c chi.Router) {
		c.Get("/get/{id}", h.handlerGetGroupByID)
		c.Post("/all", h.handlerGetGroupsByStatus)
		c.Post("/by-order-id-and-status", h.handlerGetGroupsByOrderIDAndStatus)
		c.Post("/start/{id}", h.handlerStartGroupByID)
		c.Post("/ready/{id}", h.handlerReadyGroupByID)
		c.Post("/cancel/{id}", h.handlerCancelGroupByID)
		c.Delete("/delete/{id}", h.handlerDeleteGroupByID)
	})

	return handler.NewHandler("/group", c)
}

func (h *handlerGroupItemImpl) handlerGetGroupByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if groupItem, err := h.s.GetGroupByID(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: groupItem})
	}
}

func (h *handlerGroupItemImpl) handlerGetGroupsByStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dto := &groupitemdto.GroupItemByStatusInput{}

	jsonpkg.ParseBody(r, dto)

	if groups, err := h.s.GetGroupsByStatus(ctx, dto); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: groups})
	}
}

func (h *handlerGroupItemImpl) handlerGetGroupsByOrderIDAndStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dto := &groupitemdto.GroupItemByOrderIDAndStatusInput{}

	jsonpkg.ParseBody(r, dto)

	if groups, err := h.s.GetGroupsByOrderIDAndStatus(ctx, dto); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: groups})
	}
}

func (h *handlerGroupItemImpl) handlerStartGroupByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.StartGroupItem(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}

func (h *handlerGroupItemImpl) handlerReadyGroupByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.ReadyGroupItem(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}

func (h *handlerGroupItemImpl) handlerCancelGroupByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.CancelGroupItem(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}

func (h *handlerGroupItemImpl) handlerDeleteGroupByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteGroupItem(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}
