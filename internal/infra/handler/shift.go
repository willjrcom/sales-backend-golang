package handlerimpl

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	shiftdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/shift"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
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
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.s.OpenShift(ctx, dtoShift)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, id)
}

func (h *handlerShiftImpl) handlerCloseShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoShift := &shiftdto.ShiftUpdateCloseDTO{}
	if err := jsonpkg.ParseBody(r, dtoShift); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.CloseShift(ctx, dtoShift); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerShiftImpl) handlerGetShiftByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	shift, err := h.s.GetShiftByID(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, shift)
}

func (h *handlerShiftImpl) handlerGetCurrentShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	shift, err := h.s.GetCurrentShift(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, shift)
}

func (h *handlerShiftImpl) handlerGetAllShifts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// parse pagination query params
	page, perPage := headerservice.GetPageAndPerPage(r, 0, 10)

	if shifts, err := h.s.GetAllShifts(ctx, page, perPage); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, shifts)
	}
}

func (h *handlerShiftImpl) handlerAddRedeem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoRedeem := &shiftdto.ShiftRedeemCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoRedeem); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.AddRedeem(ctx, dtoRedeem); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)

}
