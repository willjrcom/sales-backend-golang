package handlerimpl

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerCompanyImpl struct {
	s *companyusecases.Service
}

func NewHandlerCompany(companyService *companyusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerCompanyImpl{
		s: companyService,
	}

	route := "/company"
	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerNewCompany)
		c.Get("/", h.handlerGetCompany)
		c.Post("/add/user", h.handlerAddUserToCompany)
		c.Post("/remove/user", h.handlerRemoveUserFromCompany)
	})

	unprotectedRoutes := []string{
		fmt.Sprintf("%s/new", route),
	}
	return handler.NewHandler("/company", c, unprotectedRoutes...)
}

func (h *handlerCompanyImpl) handlerNewCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	company := &companydto.CompanyInput{}
	jsonpkg.ParseBody(r, company)

	id, schemaName, err := h.s.NewCompany(ctx, company)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: struct {
		ID     uuid.UUID `json:"company_id"`
		Schema *string   `json:"schema"`
	}{
		ID:     id,
		Schema: schemaName,
	}})
}

func (h *handlerCompanyImpl) handlerGetCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := h.s.GetCompany(ctx)
	if err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: id})
}

func (h *handlerCompanyImpl) handlerAddUserToCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := &companydto.UserInput{}
	jsonpkg.ParseBody(r, user)

	if err := h.s.AddUserToCompany(ctx, user); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerCompanyImpl) handlerRemoveUserFromCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := &companydto.UserInput{}
	jsonpkg.ParseBody(r, user)

	if err := h.s.RemoveUserFromCompany(ctx, user); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
