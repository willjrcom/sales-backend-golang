package companyentity

import (
	"context"

	"github.com/google/uuid"
)

type CompanyRepository interface {
	NewCompany(ctx context.Context, company *Company) error
	GetCompany(ctx context.Context) (*Company, error)
	AddUserToPublicCompany(ctx context.Context, userID uuid.UUID) error
	RemoveUserFromPublicCompany(ctx context.Context, userID uuid.UUID) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, user *User) error
	LoginUser(ctx context.Context, user *User) (*User, error)
	GetIDByEmail(ctx context.Context, email string) (uuid.UUID, error)
}
