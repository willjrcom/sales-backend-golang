package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product"
	sizeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/size_category"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerSizeCategoryImpl struct {
	pcs *sizeusecases.Service
}

func NewHandlerSizeProduct(sizeService *sizeusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerSizeCategoryImpl{
		pcs: sizeService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterSize)
		c.Put("/update/{id}", h.handlerUpdateSize)
		c.Delete("/delete/{id}", h.handlerDeleteSize)
	})

	return handler.NewHandler("/size", c)
}

func (h *handlerSizeCategoryImpl) handlerRegisterSize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	size := &productdto.RegisterSizeInput{}
	jsonpkg.ParseBody(r, size)

	if id, err := h.pcs.RegisterSize(ctx, size); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerSizeCategoryImpl) handlerUpdateSize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	Size := &productdto.UpdateSizeInput{}
	jsonpkg.ParseBody(r, Size)

	if err := h.pcs.UpdateSize(ctx, dtoId, Size); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerSizeCategoryImpl) handlerDeleteSize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.pcs.DeleteSize(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
