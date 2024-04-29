package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	processdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/process"
	processusecases "github.com/willjrcom/sales-backend-go/internal/usecases/process"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerProcessImpl struct {
	s *processusecases.Service
}

func NewHandlerProcess(processService *processusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerProcessImpl{
		s: processService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerCreateProcess)
		c.Post("/start/{id}", h.handlerStartProcess)
		c.Post("/pause/{id}", h.handlerPauseProcess)
		c.Post("/continue/{id}", h.handlerContinueProcess)
		c.Post("/finish/{id}", h.handlerFinishProcess)
		c.Get("/{id}", h.handlerGetProcess)
		c.Get("/all", h.handlerGetAllProcesses)
	})

	return handler.NewHandler("/process", c)
}

func (h *handlerProcessImpl) handlerCreateProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoProcess := &processdto.CreateProcessInput{}
	if err := jsonpkg.ParseBody(r, dtoProcess); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	process, err := h.s.CreateProcess(ctx, dtoProcess)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: process})
}

func (h *handlerProcessImpl) handlerStartProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	dtoStart := &processdto.StartProcessInput{}
	if err := jsonpkg.ParseBody(r, &dtoStart); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	if err := h.s.StartProcess(ctx, dtoId, dtoStart); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerProcessImpl) handlerPauseProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.PauseProcess(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerProcessImpl) handlerContinueProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.ContinueProcess(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerProcessImpl) handlerFinishProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.FinishProcess(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerProcessImpl) handlerGetProcess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	process, err := h.s.GetProcessById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: process})
}

func (h *handlerProcessImpl) handlerGetAllProcesses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	processes, err := h.s.GetAllProcesses(ctx)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: processes})
}
