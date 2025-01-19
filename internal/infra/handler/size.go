package handlerimpl

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	sizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/size"
	sizeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/size"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerSizeImpl struct {
	s *sizeusecases.Service
}

func NewHandlerSize(sizeService *sizeusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerSizeImpl{
		s: sizeService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerCreateSize)
		c.Patch("/update/{id}", h.handlerUpdateSize)
		c.Delete("/{id}", h.handlerDeleteSize)
	})

	return handler.NewHandler("/product-category/size", c)
}

func (h *handlerSizeImpl) handlerCreateSize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoSize := &sizedto.SizeCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoSize); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.s.CreateSize(ctx, dtoSize)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, id)
}

func (h *handlerSizeImpl) handlerUpdateSize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoSize := &sizedto.SizeUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dtoSize); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdateSize(ctx, dtoId, dtoSize); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerSizeImpl) handlerDeleteSize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteSize(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
