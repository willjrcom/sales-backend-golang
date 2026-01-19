package handlerimpl

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	employeedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/employee"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	headerservice "github.com/willjrcom/sales-backend-go/internal/infra/service/header"
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
		c.Post("/new", h.handlerCreateEmployee)
		c.Patch("/update/{id}", h.handlerUpdateEmployee)
		c.Put("/update/{id}/payment", h.handlerAddPayment)
		c.Delete("/{id}", h.handlerDeleteEmployee)
		c.Get("/{id}", h.handlerGetEmployee)
		c.Get("/all", h.handlerGetAllEmployees)
		c.Get("/all/without-delivery-drivers", h.handlerGetAllEmployeesWithoutDeliveryDrivers)
		c.Get("/deleted", h.handlerGetAllEmployeeDeleted)
		c.Get("/{id}/salary-history", h.handlerGetSalaryHistory)
		c.Get("/{id}/payments", h.handlerGetPayments)
		c.Post("/{id}/salary-history", h.handlerCreateSalaryHistory)
		c.Post("/{id}/payments", h.handlerCreatePayment)
	})

	return handler.NewHandler("/employee", c)
}

func (h *handlerEmployeeImpl) handlerCreateEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtoEmployee := &employeedto.EmployeeCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoEmployee); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	id, err := h.s.CreateEmployee(ctx, dtoEmployee)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusCreated, id)
}

func (h *handlerEmployeeImpl) handlerUpdateEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoEmployee := &employeedto.EmployeeUpdateDTO{}
	if err := jsonpkg.ParseBody(r, dtoEmployee); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.UpdateEmployee(ctx, dtoId, dtoEmployee); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}
func (h *handlerEmployeeImpl) handlerDeleteEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	if err := h.s.DeleteEmployee(ctx, dtoId); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerEmployeeImpl) handlerGetEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	employee, err := h.s.GetEmployeeById(ctx, dtoId)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, employee)
}

func (h *handlerEmployeeImpl) handlerGetAllEmployees(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	// fetch paginated employees
	employees, total, err := h.s.GetAllEmployees(ctx, page, perPage, isActive)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	// set total count header
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	jsonpkg.ResponseJson(w, r, http.StatusOK, employees)
}

// handlerAddPayment adds a payment record for an employee.
func (h *handlerEmployeeImpl) handlerAddPayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	dtoId := &entitydto.IDRequest{ID: uuid.MustParse(id)}

	dtoPayment := &employeedto.EmployeePaymentCreateDTO{}
	if err := jsonpkg.ParseBody(r, dtoPayment); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	if err := h.s.AddPayment(ctx, dtoId, dtoPayment); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, nil)
}

func (h *handlerEmployeeImpl) handlerGetAllEmployeeDeleted(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page, perPage := headerservice.GetPageAndPerPage(r, 0, 10)

	// fetch paginated deleted employees
	employees, total, err := h.s.GetAllEmployeeDeleted(ctx, page, perPage)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	// set total count header
	w.Header().Set("X-Total-Count", strconv.Itoa(total))
	jsonpkg.ResponseJson(w, r, http.StatusOK, employees)
}

func (h *handlerEmployeeImpl) handlerGetSalaryHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	employeeID, err := uuid.Parse(id)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}
	history, err := h.s.GetSalaryHistory(ctx, employeeID)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, history)
}

func (h *handlerEmployeeImpl) handlerGetPayments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	employeeID, err := uuid.Parse(id)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}
	payments, err := h.s.GetPayments(ctx, employeeID)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, payments)
}

func (h *handlerEmployeeImpl) handlerCreateSalaryHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	employeeID, err := uuid.Parse(id)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}
	dto := &employeedto.EmployeeSalaryHistoryCreateDTO{}
	if err := jsonpkg.ParseBody(r, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}
	history, err := h.s.CreateSalaryHistory(ctx, employeeID, dto)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusCreated, history)
}

func (h *handlerEmployeeImpl) handlerCreatePayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, errors.New("id is required"))
		return
	}
	employeeID, err := uuid.Parse(id)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}
	dto := &employeedto.EmployeePaymentCreateDTO{}
	if err := jsonpkg.ParseBody(r, dto); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}
	payment, err := h.s.CreatePayment(ctx, employeeID, dto)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusCreated, payment)
}

func (h *handlerEmployeeImpl) handlerGetAllEmployeesWithoutDeliveryDrivers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	employees, err := h.s.GetAllEmployeesWithoutDeliveryDrivers(ctx)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, err)
		return
	}
	jsonpkg.ResponseJson(w, r, http.StatusOK, employees)
}
