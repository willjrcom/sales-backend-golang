package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	sizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/size_category"
	sizeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/size_category"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerSizeCategoryImpl struct {
	s *sizeusecases.Service
}

func NewHandlerSizeCategory(sizeService *sizeusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerSizeCategoryImpl{
		s: sizeService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterSize)
		c.Patch("/update/{id}", h.handlerUpdateSize)
		c.Delete("/{id}", h.handlerDeleteSize)
	})

	return handler.NewHandler("/category-product/size", c)
}

func (h *handlerSizeCategoryImpl) handlerRegisterSize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoSize := &sizedto.RegisterSizeInput{}
	if err := jsonpkg.ParseBody(r, dtoSize); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	id, err := h.s.RegisterSize(ctx, dtoSize)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerSizeCategoryImpl) handlerUpdateSize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	dtoSize := &sizedto.UpdateSizeInput{}
	if err := jsonpkg.ParseBody(r, dtoSize); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.s.UpdateSize(ctx, dtoId, dtoSize); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerSizeCategoryImpl) handlerDeleteSize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteSize(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
