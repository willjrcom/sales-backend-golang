package userrepositorylocal

import (
	"context"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type UserRepositoryLocal struct {}

func NewUserRepositoryLocal() model.UserRepository {
	return &UserRepositoryLocal{}
}

func (r *UserRepositoryLocal) CreateUser(ctx context.Context, user *model.User) error {
	return nil
}

func (r *UserRepositoryLocal) UpdateUser(ctx context.Context, user *model.User) error {
	return nil
}

func (r *UserRepositoryLocal) LoginAndDeleteUser(ctx context.Context, user *model.User) error {
	return nil
}

func (r *UserRepositoryLocal) LoginUser(ctx context.Context, user *model.User) (*model.User, error) {
	return nil, nil
}

func (r *UserRepositoryLocal) GetIDByEmail(ctx context.Context, email string) (*uuid.UUID, error) {
	return nil, nil
}

func (r *UserRepositoryLocal) GetIDByEmailOrCPF(ctx context.Context, email string, cpf string) (*uuid.UUID, error) {
	return nil, nil
}

func (r *UserRepositoryLocal) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return nil, nil
}

func (r *UserRepositoryLocal) GetByCPF(ctx context.Context, cpf string) (*model.User, error) {
	return nil, nil
}

func (r *UserRepositoryLocal) ExistsUserByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return false, nil
}
