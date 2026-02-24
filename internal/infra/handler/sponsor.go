package handlerimpl

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	sponsordto "github.com/willjrcom/sales-backend-go/internal/infra/dto/sponsor"
	sponsorusecases "github.com/willjrcom/sales-backend-go/internal/usecases/sponsor"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerSponsorImpl struct {
	s *sponsorusecases.SponsorService
}

func NewHandlerSponsor(s *sponsorusecases.SponsorService) *handler.Handler {
	c := chi.NewRouter()
	h := &handlerSponsorImpl{s: s}

	c.With().Group(func(c chi.Router) {
		c.Post("/create", h.handlerCreateSponsor)
		c.Put("/update/{id}", h.handlerUpdateSponsor)
		c.Delete("/delete/{id}", h.handlerDeleteSponsor)
		c.Get("/{id}", h.handlerGetSponsorByID)
		c.Get("/all", h.handlerGetAllSponsors)
	})

	return handler.NewHandler("/sponsor", c)
}

func (h *handlerSponsorImpl) handlerCreateSponsor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoSponsor := &sponsordto.CreateSponsorDTO{}
	if err := jsonpkg.ParseBody(r, dtoSponsor); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.s.CreateSponsor(ctx, dtoSponsor)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, id)
}

func (h *handlerSponsorImpl) handlerUpdateSponsor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoSponsor := &sponsordto.UpdateSponsorDTO{}
	if err := jsonpkg.ParseBody(r, dtoSponsor); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdateSponsor(ctx, dtoId, dtoSponsor); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerSponsorImpl) handlerDeleteSponsor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteSponsor(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerSponsorImpl) handlerGetSponsorByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	sponsor, err := h.s.GetSponsorById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, sponsor)
}

func (h *handlerSponsorImpl) handlerGetAllSponsors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sponsors, err := h.s.GetAllSponsors(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, sponsors)
}
