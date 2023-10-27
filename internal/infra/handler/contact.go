package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	contactusecases "github.com/willjrcom/sales-backend-go/internal/usecases/contact_person"
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
		c.Post("/by", h.handlerGetContactsBy)
		c.Get("/{id}", h.handlerGetContactById)
		c.Get("/all", h.handlerGetAllContacts)
	})

	return handler.NewHandler("/contact", c)
}

func (h *handlerContactImpl) handlerGetContactsBy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	contact := &contactdto.FilterContact{}
	jsonpkg.ParseBody(r, contact)

	if id, err := h.s.GetContactsBy(ctx, contact); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerContactImpl) handlerGetContactById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if id, err := h.s.GetContactById(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerContactImpl) handlerGetAllContacts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if contacts, err := h.s.GetAllContacts(ctx); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: contacts})
	}
}
