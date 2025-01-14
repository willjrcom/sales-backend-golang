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
		c.Put("/close", h.handlerCloseShift)
		c.Get("/{id}", h.handlerGetShiftByID)
		c.Get("/current", h.handlerGetCurrentShift)
		c.Get("/all", h.handlerGetAllShifts)
		c.Put("/redeem/add", h.handlerAddRedeem)
	})

	return handler.NewHandler("/shift", c)
}

func (h *handlerShiftImpl) handlerOpenShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoShift := &shiftdto.ShiftUpdateOpenDTO{}
	if err := jsonpkg.ParseBody(r, dtoShift); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	id, err := h.s.OpenShift(ctx, dtoShift)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerShiftImpl) handlerCloseShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoShift := &shiftdto.ShiftUpdateCloseDTO{}
	if err := jsonpkg.ParseBody(r, dtoShift); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.s.CloseShift(ctx, dtoShift); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerShiftImpl) handlerGetShiftByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	shift, err := h.s.GetShiftByID(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: shift})
}

func (h *handlerShiftImpl) handlerGetCurrentShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	shift, err := h.s.GetCurrentShift(ctx)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: shift})
}

func (h *handlerShiftImpl) handlerGetAllShifts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if shifts, err := h.s.GetAllShifts(ctx); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: shifts})
	}
}

func (h *handlerShiftImpl) handlerAddRedeem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoRedeem := &shiftdto.ShiftRedeemCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoRedeem); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.s.AddRedeem(ctx, dtoRedeem); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)

}
