package employeedto

import (
	"errors"

	"github.com/google/uuid"
	employeeentity "github.com/willjrcom/sales-backend-go/internal/domain/employee"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrUserIDRequired = errors.New("user_id is required")
)

type CreateEmployeeInput struct {
	UserID *uuid.UUID `json:"user_id"`
}

func (r *CreateEmployeeInput) validate() error {
	if r.UserID == nil {
		return ErrUserIDRequired
	}

	return nil
}

func (r *CreateEmployeeInput) ToModel() (*employeeentity.Employee, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	// Get exists user
	return &employeeentity.Employee{
		Entity: entity.NewEntity(),
		UserID: r.UserID,
	}, nil
}
