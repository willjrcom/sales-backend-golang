package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	shiftdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/shift"
	shiftusecases "github.com/willjrcom/sales-backend-go/internal/usecases/shift"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerShiftImpl struct {
	s *shiftusecases.Service
}

func NewHandlerShift(shiftService *shiftusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerShiftImpl{
		s: shiftService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/open", h.handlerOpenShift)
		c.Put("/close/{id}", h.handlerCloseShift)
		c.Get("/{id}", h.handlerGetShiftByID)
	})

	return handler.NewHandler("/shift", c)
}

func (h *handlerShiftImpl) handlerOpenShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto := &shiftdto.OpenShift{}
	jsonpkg.ParseBody(r, dto)

	if id, err := h.s.OpenShift(ctx, dto); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerShiftImpl) handlerCloseShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto := &shiftdto.CloseShift{}
	jsonpkg.ParseBody(r, dto)

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.CloseShift(ctx, dtoId, dto); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}

func (h *handlerShiftImpl) handlerGetShiftByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if id, err := h.s.GetShiftByID(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: id})
	}
}

// func (h *handlerShiftImpl) handlerGetAllShifts(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	if shifts, err := h.s.GetAllShifts(ctx); err != nil {
// 		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
// 	} else {
// 		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: shifts})
// 	}
// }
