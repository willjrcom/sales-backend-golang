package userentity

import (
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, user *User) error
	LoginUser(ctx context.Context, user *User) (*User, error)
}
