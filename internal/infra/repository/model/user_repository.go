package model

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	UpdateUserPassword(ctx context.Context, user *User) error
	LoginAndDeleteUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, email string) (*User, error)
	LoginUser(ctx context.Context, user *User) (*User, error)
	GetIDByEmail(ctx context.Context, email string) (*uuid.UUID, error)
	GetIDByEmailOrCPF(ctx context.Context, email string, cpf string) (*uuid.UUID, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByCPF(ctx context.Context, cpf string) (*User, error)
	ExistsUserByID(ctx context.Context, id uuid.UUID) (bool, error)
}
