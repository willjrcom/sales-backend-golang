package handlerimpl

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
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
		c.Put("/update", h.handlerUpdateCompany)
		c.Get("/", h.handlerGetCompany)
		c.Get("/users", h.handlerGetCompanyUsers)
		c.Post("/add/user", h.handlerAddUserToCompany)
		c.Post("/remove/user", h.handlerRemoveUserFromCompany)
		c.Post("/test", h.handlerTest)
	})

	return handler.NewHandler("/company", c)
}

func (h *handlerCompanyImpl) handlerNewCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoCompany := &companydto.CompanyCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoCompany); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	companySchema, err := h.s.NewCompany(ctx, dtoCompany)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, companySchema)
}

func (h *handlerCompanyImpl) handlerUpdateCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoCompany := &companydto.CompanyUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dtoCompany); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdateCompany(ctx, dtoCompany); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerCompanyImpl) handlerGetCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := h.s.GetCompany(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, id)
}

func (h *handlerCompanyImpl) handlerGetCompanyUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// parse pagination query params
	page, perPage := headerservice.GetPageAndPerPage(r, 0, 10)

	users, total, err := h.s.GetCompanyUsers(ctx, page, perPage)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	jsonpkg.ResponseJson(w, r, http.StatusOK, users)
}

func (h *handlerCompanyImpl) handlerAddUserToCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &companydto.UserToCompanyDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.AddUserToCompany(ctx, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerCompanyImpl) handlerRemoveUserFromCompany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoUser := &companydto.UserToCompanyDTO{}
	if err := jsonpkg.ParseBody(r, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.RemoveUserFromCompany(ctx, dtoUser); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerCompanyImpl) handlerTest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := h.s.Test(ctx); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
