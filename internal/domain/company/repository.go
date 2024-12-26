package companyentity

import (
	"context"

	"github.com/google/uuid"
)

type CompanyRepository interface {
	NewCompany(ctx context.Context, company *Company) error
	GetCompany(ctx context.Context) (*Company, error)
	ValidateUserToPublicCompany(ctx context.Context, userID uuid.UUID) (bool, error)
	AddUserToPublicCompany(ctx context.Context, userID uuid.UUID) error
	RemoveUserFromPublicCompany(ctx context.Context, userID uuid.UUID) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	LoginAndDeleteUser(ctx context.Context, user *User) error
	LoginUser(ctx context.Context, user *User) (*User, error)
	GetIDByEmail(ctx context.Context, email string) (*uuid.UUID, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	ExistsUserByID(ctx context.Context, id uuid.UUID) (bool, error)
}
