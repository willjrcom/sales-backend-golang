package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	placedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/place"
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
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	id, err := h.s.CreatePlace(ctx, dtoPlace)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerPlaceImpl) handlerDeletePlaceById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeletePlace(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerPlaceImpl) handlerGetPlaceById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	place, err := h.s.GetPlaceById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: place})
}

func (h *handlerPlaceImpl) handlerUpdatePlaceById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoPlace := &placedto.PlaceUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dtoPlace); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.s.UpdatePlace(ctx, dtoId, dtoPlace); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerPlaceImpl) handlerGetAllPlaces(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	places, err := h.s.GetAllPlaces(ctx)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: places})
}

func (h *handlerPlaceImpl) handlerAddTableToPlace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto := &placedto.PlaceUpdateTableDTO{}
	if err := jsonpkg.ParseBody(r, dto); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.s.AddTableToPlace(ctx, dto); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerPlaceImpl) handlerRemoveTableFromPlace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.RemoveTableFromPlace(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
