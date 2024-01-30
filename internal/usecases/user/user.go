package userusecases

import (
	"context"

	userentity "github.com/willjrcom/sales-backend-go/internal/domain/user"
	userdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/user"
	bcryptservice "github.com/willjrcom/sales-backend-go/internal/infra/service/bcrypt"
	jwtservice "github.com/willjrcom/sales-backend-go/internal/infra/service/jwt"
)

type Service struct {
	r userentity.Repository
}

func NewService(c userentity.Repository) *Service {
	return &Service{r: c}
}

func (s *Service) CreateUser(ctx context.Context, dto *userdto.CreateUserInput) error {
	user, err := dto.ToModel()

	if err != nil {
		return err
	}

	hash, err := bcryptservice.HashPassword(dto.Password)
	if err != nil {
		return err
	}

	user.Hash = string(hash)

	return s.r.CreateUser(ctx, user)
}

func (s *Service) UpdateUser(ctx context.Context, dto *userdto.UpdatePasswordInput) error {
	user, err := dto.ToModel()

	if err != nil {
		return err
	}

	userLoggedIn, err := s.r.LoginUser(ctx, user)
	if err != nil {
		return err
	}

	hash, err := bcryptservice.HashPassword(dto.NewPassword)
	if err != nil {
		return err
	}

	userLoggedIn.Hash = string(hash)

	return s.r.UpdateUser(ctx, user)
}

func (s *Service) LoginUser(ctx context.Context, dto *userdto.LoginUserInput) (token string, err error) {
	user, err := dto.ToModel()

	if err != nil {
		return "", err
	}

	userLoggedIn, err := s.r.LoginUser(ctx, user)

	if err != nil {
		return "", err
	}

	return jwtservice.CreateToken(&userLoggedIn.UserCommonAttributes)
}

func (s *Service) DeleteUser(ctx context.Context, dto *userdto.DeleteUserInput) error {
	user, err := dto.ToModel()

	if err != nil {
		return err
	}

	return s.r.DeleteUser(ctx, user)
}
