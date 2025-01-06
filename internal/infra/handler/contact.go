package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	keysdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/keys"
	contactusecases "github.com/willjrcom/sales-backend-go/internal/usecases/contact"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerContactImpl struct {
	s *contactusecases.Service
}

func NewHandlerContactPerson(contactService *contactusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerContactImpl{
		s: contactService,
	}

	c.With().Group(func(c chi.Router) {
		c.Get("/{id}", h.handlerGetContactById)
		c.Post("/search", h.handlerFtSearchContacts)
	})

	return handler.NewHandler("/contact", c)
}

func (h *handlerContactImpl) handlerGetContactById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
		return
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	contact, err := h.s.GetContactById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: contact})
}

func (h *handlerContactImpl) handlerFtSearchContacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	keys := &keysdto.KeysDTO{}
	if err := jsonpkg.ParseBody(r, keys); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: err.Error()})
		return
	}

	contacts, err := h.s.FtSearchContacts(ctx, keys)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: contacts})
}
