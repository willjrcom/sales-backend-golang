package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productcategorysizedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category_size"
	productcategorysizeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/product_category_size"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerSizeImpl struct {
	s *productcategorysizeusecases.Service
}

func NewHandlerSize(sizeService *productcategorysizeusecases.Service, path string) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerSizeImpl{
		s: sizeService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterSize)
		c.Patch("/update/{id}", h.handlerUpdateSize)
		c.Delete("/{id}", h.handlerDeleteSize)
	})

	return handler.NewHandler(path+"/size", c)
}

func (h *handlerSizeImpl) handlerRegisterSize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoSize := &productcategorysizedto.RegisterSizeInput{}
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

func (h *handlerSizeImpl) handlerUpdateSize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	dtoSize := &productcategorysizedto.UpdateSizeInput{}
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

func (h *handlerSizeImpl) handlerDeleteSize(w http.ResponseWriter, r *http.Request) {
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
