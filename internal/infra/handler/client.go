package handlerimpl

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	clientdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/client"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
	clientusecases "github.com/willjrcom/sales-backend-go/internal/usecases/client"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerClientImpl struct {
	s *clientusecases.Service
}

func NewHandlerClient(clientService *clientusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerClientImpl{
		s: clientService,
	}

	route := "/client"

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerCreateClient)
		c.Patch("/update/{id}", h.handlerUpdateClient)
		c.Delete("/{id}", h.handlerDeleteClient)
		c.Get("/{id}", h.handlerGetClientById)
		c.Post("/by-contact", h.handlerGetClientByContact)
		c.Get("/all", h.handlerGetAllClients)
	})

	unprotectedRoutes := []string{}
	return handler.NewHandler(route, c, unprotectedRoutes...)
}

func (h *handlerClientImpl) handlerCreateClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoClient := &clientdto.ClientCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoClient); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.s.CreateClient(ctx, dtoClient)

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, id)
}

func (h *handlerClientImpl) handlerUpdateClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoClient := &clientdto.ClientUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dtoClient); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdateClient(ctx, dtoId, dtoClient); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerClientImpl) handlerDeleteClient(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteClient(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerClientImpl) handlerGetClientById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	client, err := h.s.GetClientById(ctx, dtoId)

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, client)
}

func (h *handlerClientImpl) handlerGetClientByContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoContact := &contactdto.ContactDTO{}
	if err := jsonpkg.ParseBody(r, dtoContact); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	client, err := h.s.GetClientByContact(ctx, dtoContact)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, client)
}

func (h *handlerClientImpl) handlerGetAllClients(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// parse pagination query params
	page, perPage := headerservice.GetPageAndPerPage(r, 0, 10)

	// parse is_active query parameter (default: true)
	isActive := true
	if isActiveParam := r.URL.Query().Get("is_active"); isActiveParam != "" {
		var err error
		isActive, err = strconv.ParseBool(isActiveParam)
		if err != nil {
			jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("invalid is_active parameter"))
			return
		}
	}

	// fetch paginated clients from service
	clients, total, err := h.s.GetAllClients(ctx, page, perPage, isActive)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	// set total count header
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	// respond with paginated clients
	jsonpkg.ResponseJson(w, r, http.StatusOK, clients)
}
