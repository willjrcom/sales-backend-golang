package handlerimpl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
	employeedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/employee"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	employeeusecases "github.com/willjrcom/sales-backend-go/internal/usecases/employee"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

type handlerEmployeeImpl struct {
	s *employeeusecases.Service
}

func NewHandlerEmployee(employeeService *employeeusecases.Service) *handler.Handler {
	c := chi.NewRouter()

	h := &handlerEmployeeImpl{
		s: employeeService,
	}

	c.With().Group(func(c chi.Router) {
		c.Post("/new", h.handlerRegisterEmployee)
		c.Patch("/update/{id}", h.handlerUpdateEmployee)
		c.Delete("/delete/{id}", h.handlerDeleteEmployee)
		c.Get("/{id}", h.handlerGetEmployee)
		c.Get("/all", h.handlerGetAllEmployees)
	})

	c.With().Group(func(c chi.Router) {
		c.Post("/contact/new", h.handlerRegisterContactEmployee)
		c.Patch("/contact/update/{id}", h.handlerUpdateContactEmployee)
		c.Delete("/contact/delete/{id}", h.handlerDeleteContactEmployee)
	})

	return handler.NewHandler("/employee", c)
}

func (h *handlerEmployeeImpl) handlerRegisterEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	employee := &employeedto.RegisterEmployeeInput{}
	jsonpkg.ParseBody(r, employee)

	if id, err := h.s.RegisterEmployee(ctx, employee); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerEmployeeImpl) handlerUpdateEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	employee := &employeedto.UpdateEmployeeInput{}
	jsonpkg.ParseBody(r, employee)

	if err := h.s.UpdateEmployee(ctx, dtoId, employee); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}

func (h *handlerEmployeeImpl) handlerDeleteEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteEmployee(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}

func (h *handlerEmployeeImpl) handlerGetEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if employee, err := h.s.GetEmployeeById(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: employee})
	}
}

func (h *handlerEmployeeImpl) handlerGetAllEmployees(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if categories, err := h.s.GetAllEmployees(ctx); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, jsonpkg.HTTPResponse{Data: categories})
	}
}

func (h *handlerEmployeeImpl) handlerRegisterContactEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contact := &contactdto.RegisterContactInput{}
	jsonpkg.ParseBody(r, contact)

	if id, err := h.s.RegisterContactToEmployee(ctx, contact); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusCreated, jsonpkg.HTTPResponse{Data: id})
	}
}

func (h *handlerEmployeeImpl) handlerUpdateContactEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	contact := &contactdto.UpdateContactInput{}
	jsonpkg.ParseBody(r, contact)

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.UpdateContact(ctx, dtoId, contact); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}

func (h *handlerEmployeeImpl) handlerDeleteContactEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseJson(w, r, http.StatusBadRequest, jsonpkg.Error{Message: "id is required"})
	}

	dtoId := &entitydto.IdRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteContact(ctx, dtoId); err != nil {
		jsonpkg.ResponseJson(w, r, http.StatusInternalServerError, jsonpkg.Error{Message: err.Error()})
	} else {
		jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
	}
}
