package handlerimpl

import (
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

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerNewCompany)
		c.Get("/", h.handlerGetCompany)
	})

	return handler.NewHandler("/company", c)
}

func (h *handlerCompanyImpl) handlerNewCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	company := &companydto.CompanyInput{}
	jsonpkg.ParseBody(r, company)

	if id, schemaName, err := h.s.NewCompany(ctx, company); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: struct {
			ID     uuid.UUID `json:"company_id"`
			Schema *string   `json:"schema"`
		}{
			ID:     id,
			Schema: schemaName,
		}})
	}
}

func (h *handlerCompanyImpl) handlerGetCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if id, err := h.s.GetCompany(ctx); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: id})
	}
}
