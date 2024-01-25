package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
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
		c.Get("/{id}", h.handlerGetCompanyById)
		c.Get("/all", h.handlerGetAllCompaniesBySchemaName)
	})

	return handler.NewHandler("/company", c)
}

func (h *handlerCompanyImpl) handlerNewCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	company := &companydto.CompanyInput{}
	jsonpkg.ParseBody(r, company)

	if id, err := h.s.NewCompany(ctx, company); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerCompanyImpl) handlerGetCompanyById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if id, err := h.s.GetCompanyById(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerCompanyImpl) handlerGetAllCompaniesBySchemaName(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	schema := &companydto.CompanyBySchemaName{}
	jsonpkg.ParseBody(r, schema)

	if companys, err := h.s.GetAllCompaniesBySchemaName(ctx, schema); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: companys})
	}
}
