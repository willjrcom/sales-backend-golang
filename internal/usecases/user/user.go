package userusecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
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

func (s *Service) CreateUser(ctx context.Context, dto *companydto.UserCreateDTO) (*uuid.UUID, error) {
	user, err := dto.ToDomain()

	if err != nil {
		return nil, err
	}

	if id, _ := s.r.GetIDByEmailOrCPF(ctx, user.Email, user.Cpf); id != nil {
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

func (s *Service) UpdateUserPassword(ctx context.Context, dto *companydto.UserUpdatePasswordDTO) error {
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

func (s *Service) ForgetUserPassword(ctx context.Context, dto *companydto.UserForgetPasswordDTO) error {
	email, err := dto.ToDomain()
	if err != nil {
		return err
	}

	// Send email service
	if err := emailservice.SendEmail(email); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateUser(ctx context.Context, dtoID *entitydto.IDRequest, dto *companydto.UserUpdateDTO) error {
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

func (s *Service) LoginUser(ctx context.Context, dto *companydto.UserLoginDTO) (data *companydto.UserTokenDTO, err error) {
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
	userLoggedIn.Companies = []companyentity.Company{}

	idToken, err := jwtservice.CreateIDToken(userLoggedIn)

	if err != nil {
		return nil, err
	}

	data = &companydto.UserTokenDTO{
		IDToken: idToken,
	}

	data.User.FromDomain(userLoggedIn)

	return data, nil
}

func (s *Service) Access(ctx context.Context, dto *companydto.UserSchemaDTO, accessToken *jwt.Token) (token string, err error) {
	schema, err := dto.ToDomain()

	if err != nil {
		return "", err
	}

	userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)
	if !ok {
		return "", errors.New("context user not found")
	}

	userIDUUID := uuid.MustParse(userID)
	userModel, err := s.r.GetUserByID(ctx, userIDUUID)
	if err != nil {
		return "", err
	}

	var company *model.Company
	for _, c := range userModel.Companies {
		if c.SchemaName == *schema {
			company = &c
		}
	}

	if company == nil {
		return "", errors.New("schema not found in schemas in token")
	}

	if company.IsBlocked {
		return "", fmt.Errorf("company %s is blocked", company.BusinessName)
	}

	return jwtservice.CreateAccessToken(accessToken, *schema)
}

func (s *Service) SearchUser(ctx context.Context, dto *companydto.UserSearchDTO) (*companydto.UserDTO, error) {
	cpf, err := dto.ToDomain()

	if err != nil {
		return nil, err
	}

	userModel, err := s.r.GetByCPF(ctx, cpf)
	if err != nil {
		return nil, err
	}

	user := userModel.ToDomain()

	userDTO := &companydto.UserDTO{}
	userDTO.FromDomain(user)

	return userDTO, nil
}

func (s *Service) DeleteUser(ctx context.Context, dto *companydto.UserDeleteDTO) error {
	user, err := dto.ToDomain()

	if err != nil {
		return err
	}

	userModel := &model.User{}
	userModel.FromDomain(user)
	return s.r.LoginAndDeleteUser(ctx, userModel)
}

func (s *Service) GetCompanies(ctx context.Context) ([]companydto.CompanyDTO, error) {
	userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)
	if !ok {
		return nil, errors.New("context user not found")
	}

	userIDUUID := uuid.MustParse(userID)
	userModel, err := s.r.GetUserByID(ctx, userIDUUID)
	if err != nil {
		return nil, err
	}

	companyDTOs := []companydto.CompanyDTO{}
	for _, companyModel := range userModel.Companies {
		company := companyModel.ToDomain()
		companyDTO := companydto.CompanyDTO{}
		companyDTO.FromDomain(company)
		companyDTOs = append(companyDTOs, companyDTO)
	}

	return companyDTOs, nil
}
