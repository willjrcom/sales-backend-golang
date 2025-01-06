package userusecases

import (
	"context"
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	userdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/user"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	bcryptservice "github.com/willjrcom/sales-backend-go/internal/infra/service/bcrypt"
	emailservice "github.com/willjrcom/sales-backend-go/internal/infra/service/email"
	jwtservice "github.com/willjrcom/sales-backend-go/internal/infra/service/jwt"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidPassword   = errors.New("invalid password")
)

type Service struct {
	r model.UserRepository
}

func NewService(r model.UserRepository) *Service {
	return &Service{r: r}
}

func (s *Service) CreateUser(ctx context.Context, dto *userdto.UserCreateDTO) (*uuid.UUID, error) {
	user, err := dto.ToDomain()

	if err != nil {
		return nil, err
	}

	if id, _ := s.r.GetIDByEmail(ctx, user.Email); id != nil {
		return nil, ErrUserAlreadyExists
	}

	hash, err := bcryptservice.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	user.Hash = string(hash)

	userModel := &model.User{}
	userModel.FromDomain(user)

	err = s.r.CreateUser(ctx, userModel)
	return &user.ID, err
}

func (s *Service) UpdateUserPassword(ctx context.Context, dto *userdto.UserUpdatePasswordDTO) error {
	user, err := dto.ToDomain()

	if err != nil {
		return err
	}

	if id, _ := s.r.GetIDByEmail(ctx, user.Email); id == nil {
		return ErrInvalidEmail
	}

	userModel := &model.User{}
	userModel.FromDomain(user)

	userLoggedInModel, err := s.r.LoginUser(ctx, userModel)
	if err != nil {
		return ErrInvalidPassword
	}

	userLoggedIn := userLoggedInModel.ToDomain()

	hash, err := bcryptservice.HashPassword(dto.NewPassword)
	if err != nil {
		return err
	}

	userLoggedIn.Hash = string(hash)

	userModel.FromDomain(userLoggedIn)
	return s.r.UpdateUser(ctx, userModel)
}

func (s *Service) ForgetUserPassword(ctx context.Context, dto *userdto.UserForgetPasswordDTO) error {
	email, err := dto.ToDomain()
	if err != nil {
		return err
	}

	// Send email service
	emailservice.SendEmail(email)
	return nil
}

func (s *Service) UpdateUser(ctx context.Context, dtoID *entitydto.IDRequest, dto *userdto.UserUpdateDTO) error {
	userModel, err := s.r.GetUserByID(ctx, dtoID.ID)
	if err != nil {
		return err
	}

	user := userModel.ToDomain()

	if err = dto.UpdateDomain(user); err != nil {
		return err
	}

	userModel.FromDomain(user)
	return s.r.UpdateUser(ctx, userModel)
}

func (s *Service) LoginUser(ctx context.Context, dto *userdto.UserLoginDTO) (data *userdto.UserTokenAndSchemasDTO, err error) {
	user, err := dto.ToDomain()

	if err != nil {
		return nil, err
	}

	userModel := &model.User{}
	userModel.FromDomain(user)
	userLoggedInModel, err := s.r.LoginUser(ctx, userModel)

	if err != nil {
		return nil, err
	}

	userLoggedIn := userLoggedInModel.ToDomain()

	accessToken, err := jwtservice.CreateAccessToken(userLoggedIn)

	if err != nil {
		return nil, err
	}

	data = &userdto.UserTokenAndSchemasDTO{
		Person:      userLoggedIn.Person,
		AccessToken: accessToken,
		Companies:   userLoggedIn.Companies,
	}
	return data, nil
}

func (s *Service) Access(ctx context.Context, dto *userdto.UserSchemaDTO, accessToken *jwt.Token) (token string, err error) {
	schema, err := dto.ToDomain()

	if err != nil {
		return "", err
	}

	schemasInterface := jwtservice.GetSchemasFromToken(accessToken)

	if len(schemasInterface) == 0 {
		return "", errors.New("schemas not found in token")
	}

	if !findSchemaInSchemas(schemasInterface, *schema) {
		return "", errors.New("schema not found in schemas in token")
	}

	return jwtservice.CreateIDToken(accessToken, *schema)
}

func findSchemaInSchemas(schemas []interface{}, schema string) bool {
	for _, s := range schemas {
		schemaCompany, ok := s.(string)

		if !ok {
			continue
		}
		if strings.EqualFold(schemaCompany, schema) {
			return true
		}
	}

	return false
}

func (s *Service) DeleteUser(ctx context.Context, dto *userdto.UserDeleteDTO) error {
	user, err := dto.ToDomain()

	if err != nil {
		return err
	}

	userModel := &model.User{}
	userModel.FromDomain(user)
	return s.r.LoginAndDeleteUser(ctx, userModel)
}
