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
		c.Get("/{id}", h.handlerGetGroupByID)
		c.Post("/all-by-status", h.handlerGetGroupsByStatus)
		c.Post("/by-order-id-and-status", h.handlerGetGroupsByOrderIDAndStatus)
		c.Post("/start/{id}", h.handlerStartGroupByID)
		c.Post("/ready/{id}", h.handlerReadyGroupByID)
		c.Post("/cancel/{id}", h.handlerCancelGroupByID)
		c.Delete("/{id}", h.handlerDeleteGroupByID)
		c.Post("/update/{id}/complement-item/{id-complement}", h.handlerAddComplementItem)
		c.Delete("/update/{id}/complement-item", h.handlerDeleteComplementItem)
	})

	return handler.NewHandler("/group-item", c)
}

func (h *handlerGroupItemImpl) handlerGetGroupByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	groupItem, err := h.s.GetGroupByID(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: groupItem})
}

func (h *handlerGroupItemImpl) handlerGetGroupsByStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoGroupItem := &groupitemdto.GroupItemByStatusInput{}
	if err := jsonpkg.ParseBody(r, dtoGroupItem); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	groups, err := h.s.GetGroupsByStatus(ctx, dtoGroupItem)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: groups})
}

func (h *handlerGroupItemImpl) handlerGetGroupsByOrderIDAndStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoGroupItem := &groupitemdto.GroupItemByOrderIDAndStatusInput{}
	if err := jsonpkg.ParseBody(r, dtoGroupItem); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	groups, err := h.s.GetGroupsByOrderIDAndStatus(ctx, dtoGroupItem)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: groups})
}

func (h *handlerGroupItemImpl) handlerStartGroupByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.StartGroupItem(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerGroupItemImpl) handlerReadyGroupByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.ReadyGroupItem(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerGroupItemImpl) handlerCancelGroupByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.CancelGroupItem(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerGroupItemImpl) handlerDeleteGroupByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteGroupItem(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerGroupItemImpl) handlerAddComplementItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	idComplement := chi.URLParam(r, "id-complement")

	if idComplement == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id complement is required"})
		return
	}

	dtoIdComplement := &entitydto.IdRequest{ID: uuid.MustParse(idComplement)}

	if err := h.s.AddComplementItem(ctx, dtoId, dtoIdComplement); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, nil)
}

func (h *handlerGroupItemImpl) handlerDeleteComplementItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteComplementItem(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
