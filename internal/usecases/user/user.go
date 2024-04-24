package userusecases

import (
	"context"
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	userdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/user"
	bcryptservice "github.com/willjrcom/sales-backend-go/internal/infra/service/bcrypt"
	jwtservice "github.com/willjrcom/sales-backend-go/internal/infra/service/jwt"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidPassword   = errors.New("invalid password")
)

type Service struct {
	r companyentity.UserRepository
}

func NewService(r companyentity.UserRepository) *Service {
	return &Service{r: r}
}

func (s *Service) CreateUser(ctx context.Context, dto *userdto.CreateUserInput) (uuid.UUID, error) {
	user, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if id, _ := s.r.GetIDByEmail(ctx, user.Email); id != uuid.Nil {
		return uuid.Nil, ErrUserAlreadyExists
	}

	hash, err := bcryptservice.HashPassword(dto.Password)
	if err != nil {
		return uuid.Nil, err
	}

	user.Hash = string(hash)

	return user.ID, s.r.CreateUser(ctx, user)
}

func (s *Service) UpdateUser(ctx context.Context, dto *userdto.UpdatePasswordInput) error {
	user, err := dto.ToModel()

	if err != nil {
		return err
	}

	if id, _ := s.r.GetIDByEmail(ctx, user.Email); id == uuid.Nil {
		return ErrInvalidEmail
	}

	userLoggedIn, err := s.r.LoginUser(ctx, user)
	if err != nil {
		return ErrInvalidPassword
	}

	hash, err := bcryptservice.HashPassword(dto.NewPassword)
	if err != nil {
		return err
	}

	userLoggedIn.Hash = string(hash)

	return s.r.UpdateUser(ctx, user)
}

func (s *Service) LoginUser(ctx context.Context, dto *userdto.LoginUserInput) (data *userdto.TokenAndSchemasOutput, err error) {
	user, err := dto.ToModel()

	if err != nil {
		return nil, err
	}

	userLoggedIn, err := s.r.LoginUser(ctx, user)

	if err != nil {
		return nil, err
	}

	accessToken, err := jwtservice.CreateAccessToken(userLoggedIn)

	if err != nil {
		return nil, err
	}

	data = &userdto.TokenAndSchemasOutput{
		AccessToken: accessToken,
		Companies:   userLoggedIn.Companies,
	}
	return data, nil
}

func (s *Service) Access(ctx context.Context, dto *userdto.AccessCompanyInput, accessToken *jwt.Token) (token string, err error) {
	schema, err := dto.ToModel()

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

func (s *Service) DeleteUser(ctx context.Context, dto *userdto.DeleteUserInput) error {
	user, err := dto.ToModel()

	if err != nil {
		return err
	}

	return s.r.LoginAndDeleteUser(ctx, user)
}
