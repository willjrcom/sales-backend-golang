package userrepositorylocal

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type UserRepositoryLocal struct {
	users map[uuid.UUID]*model.User
	mu    sync.RWMutex
}

func NewUserRepositoryLocal() model.UserRepository {
	return &UserRepositoryLocal{
		users: make(map[uuid.UUID]*model.User),
	}
}

func (r *UserRepositoryLocal) CreateUser(ctx context.Context, user *model.User) error {
	if user == nil || user.Entity.ID == uuid.Nil {
		return errors.New("invalid user")
	}

	if _, exists := r.users[user.Entity.ID]; exists {
		return errors.New("user already exists")
	}
	r.users[user.Entity.ID] = user
	return nil
}
func (r *UserRepositoryLocal) GetUser(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *UserRepositoryLocal) UpdateUserPassword(ctx context.Context, user *model.User) error {
	if user == nil || user.Entity.ID == uuid.Nil {
		return errors.New("invalid user")
	}

	if _, exists := r.users[user.Entity.ID]; !exists {
		return errors.New("user not found")
	}
	r.users[user.Entity.ID] = user
	return nil
}

func (r *UserRepositoryLocal) UpdateUser(ctx context.Context, user *model.User) error {
	if user == nil || user.Entity.ID == uuid.Nil {
		return errors.New("invalid user")
	}

	if _, exists := r.users[user.Entity.ID]; !exists {
		return errors.New("user not found")
	}
	r.users[user.Entity.ID] = user
	return nil
}

func (r *UserRepositoryLocal) LoginAndDeleteUser(ctx context.Context, user *model.User) error {
	if user == nil || user.Entity.ID == uuid.Nil {
		return errors.New("invalid user")
	}

	if _, exists := r.users[user.Entity.ID]; !exists {
		return errors.New("user not found")
	}
	delete(r.users, user.Entity.ID)
	return nil
}

func (r *UserRepositoryLocal) ListPublicUsers(ctx context.Context) ([]model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]model.User, 0, len(r.users))
	for _, user := range r.users {
		result = append(result, *user)
	}
	return result, nil
}

func (r *UserRepositoryLocal) LoginUser(ctx context.Context, user *model.User) (*model.User, error) {
	if user == nil {
		return nil, errors.New("invalid user")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Simple linear search by email and password
	for _, u := range r.users {
		if u.Email == user.Email && u.Password == user.Password {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *UserRepositoryLocal) GetIDByEmail(ctx context.Context, email string) (*uuid.UUID, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Email == email {
			return &u.Entity.ID, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *UserRepositoryLocal) GetIDByEmailOrCPF(ctx context.Context, email string, cpf string) (*uuid.UUID, error) {
	if email == "" && cpf == "" {
		return nil, errors.New("email and cpf cannot both be empty")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if (email != "" && u.Email == email) || (cpf != "" && u.Cpf == cpf) {
			return &u.Entity.ID, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *UserRepositoryLocal) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid id")
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (r *UserRepositoryLocal) GetByCPF(ctx context.Context, cpf string) (*model.User, error) {
	return nil, nil
}

func (r *UserRepositoryLocal) ExistsUserByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return false, nil
}
