package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
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
		c.Post("/get/{id}", h.handlerGetGroupByID)
		c.Delete("/delete/{id}", h.handlerDeleteGroupByID)
		c.Get("/all", h.handlerGetAllPendingGroups)
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

func (h *handlerGroupItemImpl) handlerGetAllPendingGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if groups, err := h.s.GetAllPendingGroups(ctx); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: groups})
	}
}
