package handlerimpl

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	advertisingdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/advertising"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	advertisingusecases "github.com/willjrcom/sales-backend-go/internal/usecases/advertising"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerAdvertisingImpl struct {
	s *advertisingusecases.AdvertisingService
}

func NewHandlerAdvertising(s *advertisingusecases.AdvertisingService) *handler.Handler {
	c := chi.NewRouter()
	h := &handlerAdvertisingImpl{s: s}

	c.With().Group(func(c chi.Router) {
		c.Post("/create", h.handlerCreateAdvertising)
		c.Put("/update/{id}", h.handlerUpdateAdvertising)
		c.Delete("/delete/{id}", h.handlerDeleteAdvertising)
		c.Get("/{id}", h.handlerGetAdvertisingByID)
		c.Get("/all", h.handlerGetAllAdvertisements)
		c.Get("/active", h.handlerGetActiveAdvertisements)
	})

	return handler.NewHandler("/advertising", c)
}

func (h *handlerAdvertisingImpl) handlerCreateAdvertising(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto := &advertisingdto.CreateAdvertisingDTO{}
	if err := jsonpkg.ParseBody(r, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.s.CreateAdvertising(ctx, dto)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, id)
}

func (h *handlerAdvertisingImpl) handlerUpdateAdvertising(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dto := &advertisingdto.UpdateAdvertisingDTO{}
	if err := jsonpkg.ParseBody(r, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdateAdvertising(ctx, dtoId, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerAdvertisingImpl) handlerDeleteAdvertising(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteAdvertising(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerAdvertisingImpl) handlerGetAdvertisingByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	adv, err := h.s.GetAdvertisingById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, adv)
}

func (h *handlerAdvertisingImpl) handlerGetAllAdvertisements(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	advs, err := h.s.GetAllAdvertisements(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, advs)
}
func (h *handlerAdvertisingImpl) handlerGetActiveAdvertisements(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	advs, err := h.s.GetActiveAdvertisements(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, advs)
}
