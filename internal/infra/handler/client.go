package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	clientdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/client"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	clientusecases "github.com/willjrcom/sales-backend-go/internal/usecases/client"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerClientImpl struct {
	ps *clientusecases.Service
}

func NewHandlerClient(clientService *clientusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerClientImpl{
		ps: clientService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterClient)
		c.Put("/update/{id}", h.handlerUpdateClient)
		c.Delete("/delete/{id}", h.handlerDeleteClient)
		c.Get("/{id}", h.handlerGetClient)
		c.Post("/by", h.handlerGetClientsBy)
		c.Post("/all", h.handlerGetAllClients)
	})

	return handler.NewHandler("/client", c)
}

func (h *handlerClientImpl) handlerRegisterClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	client := &clientdto.RegisterClientInput{}
	jsonpkg.ParseBody(r, client)

	if id, err := h.ps.RegisterClient(ctx, client); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerClientImpl) handlerUpdateClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	client := &clientdto.UpdateClientInput{}
	jsonpkg.ParseBody(r, client)

	if err := h.ps.UpdateClient(ctx, dtoId, client); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}

func (h *handlerClientImpl) handlerDeleteClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.ps.DeleteClient(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}

func (h *handlerClientImpl) handlerGetClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if client, err := h.ps.GetClientById(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: client})
	}
}

func (h *handlerClientImpl) handlerGetClientsBy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter := &clientdto.FilterClientInput{}
	jsonpkg.ParseBody(r, filter)

	if client, err := h.ps.GetClientsBy(ctx, filter); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: client})
	}
}

func (h *handlerClientImpl) handlerGetAllClients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if categories, err := h.ps.GetAllClients(ctx); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: categories})
	}
}
