package handlerimpl

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	placedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/place"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
	placeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/place"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerPlaceImpl struct {
	s *placeusecases.Service
}

func NewHandlerPlace(orderService *placeusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerPlaceImpl{
		s: orderService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerCreatePlace)
		c.Delete("/{id}", h.handlerDeletePlaceById)
		c.Patch("/update/{id}", h.handlerUpdatePlaceById)
		c.Get("/{id}", h.handlerGetPlaceById)
		c.Get("/all", h.handlerGetAllPlaces)
		c.Post("/table", h.handlerAddTableToPlace)
		c.Delete("/table/{id}", h.handlerRemoveTableFromPlace)
	})

	return handler.NewHandler("/place", c)
}

func (h *handlerPlaceImpl) handlerCreatePlace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoPlace := &placedto.CreatePlaceInput{}
	if err := jsonpkg.ParseBody(r, dtoPlace); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.s.CreatePlace(ctx, dtoPlace)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, id)
}

func (h *handlerPlaceImpl) handlerDeletePlaceById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeletePlace(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerPlaceImpl) handlerGetPlaceById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	place, err := h.s.GetPlaceById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, place)
}

func (h *handlerPlaceImpl) handlerUpdatePlaceById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoPlace := &placedto.PlaceUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dtoPlace); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdatePlace(ctx, dtoId, dtoPlace); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerPlaceImpl) handlerGetAllPlaces(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	isActive := true
	if isActiveParam := r.URL.Query().Get("is_active"); isActiveParam != "" {
		var err error
		isActive, err = strconv.ParseBool(isActiveParam)
		if err != nil {
			jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("invalid is_active parameter"))
			return
		}
	}

	page, perPage := headerservice.GetPageAndPerPage(r, 0, 100)

	places, count, err := h.s.GetAllPlaces(ctx, page, perPage, isActive)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("X-Total-Count", strconv.Itoa(count))
	jsonpkg.ResponseJson(w, r, http.StatusOK, places)
}

func (h *handlerPlaceImpl) handlerAddTableToPlace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto := &placedto.PlaceUpdateTableDTO{}
	if err := jsonpkg.ParseBody(r, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.AddTableToPlace(ctx, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerPlaceImpl) handlerRemoveTableFromPlace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.RemoveTableFromPlace(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
