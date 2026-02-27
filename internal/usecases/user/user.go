package userusecases

import (
	"context"
	"errors"
	"fmt"
	"os"

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
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
)

type Service struct {
	r            model.UserRepository
	emailService *emailservice.Service
}

func NewService(r model.UserRepository) *Service {
	return &Service{r: r}
}

func (s *Service) AddDependencies(emailService *emailservice.Service) {
	s.emailService = emailService
}

func (s *Service) CreateUser(ctx context.Context, dto *companydto.UserCreateDTO) (*uuid.UUID, error) {
	user, err := dto.ToDomain()

	if err != nil {
		return nil, err
	}

	if id, key, _ := s.r.GetIDByEmailOrCPF(ctx, user.Email, user.Cpf); id != nil {
		switch key {
		case "email":
			return nil, fmt.Errorf("user email already exists")
		case "cpf":
			return nil, fmt.Errorf("user cpf already exists")
		}
	}

	hash, err := bcryptservice.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	user.Hash = string(hash)

	userModel := &model.User{}
	userModel.FromDomain(user)

	if err = s.r.CreateUser(ctx, userModel); err != nil {
		return nil, err
	}

	// send email
	go func() {
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:3000"
		}

		bodyEmail := &emailservice.BodyEmail{
			Email:   user.Email,
			Subject: "Bem-vindo à GFood",
			Body: `<div style="font-family: Arial, sans-serif; max-width: 480px; margin: 0 auto; background: #fff; border-radius: 8px; box-shadow: 0 2px 8px #0001; padding: 32px;">
			<h2 style="color: #eab308; margin-bottom: 16px;">Bem-vindo à GFood</h2>
			
			<p style="color: #333; font-size: 16px; margin-bottom: 24px;">
				Sua conta foi criada com sucesso: ` + dto.Email + `.<br>
				Clique no botão abaixo para acessar o sistema:
			</p>

			<div style="text-align: center; margin-bottom: 24px;">
				<a href="` + frontendURL + `/login" style="display: inline-block; background-color: #eab308; color: #ffffff; font-size: 16px; font-weight: bold; padding: 12px 24px; text-decoration: none; border-radius: 6px;">
					Acessar Sistema
				</a>
			</div>
			
			<p style="color: #999; font-size: 13px; margin-top: 24px;">
				Atenciosamente,<br>
				Equipe GFood
			</p>
			</div>
		`,
		}

		// Send email service
		if s.emailService != nil {
			if err := s.emailService.SendEmail(bodyEmail); err != nil {
				fmt.Println("Error sending email:", err)
			}

			emailCC := os.Getenv("EMAIL_CC")
			if emailCC != "" {
				bodyEmailCC := *bodyEmail
				bodyEmailCC.Email = emailCC
				if err := s.emailService.SendEmail(&bodyEmailCC); err != nil {
					fmt.Println("Error sending email:", err)
				}
			}
		}
	}()
	return &user.ID, err
}

func (s *Service) UpdateUserForgetPassword(ctx context.Context, dto *companydto.UserUpdateForgetPasswordDTO) error {
	user, err := dto.ToDomain()

	if err != nil {
		return err
	}

	userModel, _ := s.r.GetUser(ctx, user.Email)
	if userModel == nil {
		return ErrInvalidEmail
	}

	hash, err := bcryptservice.HashPassword(dto.Password)
	if err != nil {
		return err
	}

	userModel.Hash = string(hash)
	return s.r.UpdateUserPassword(ctx, userModel)
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
	return s.r.UpdateUserPassword(ctx, userModel)
}

func (s *Service) ForgetUserPassword(ctx context.Context, dto *companydto.UserForgetPasswordDTO) error {
	email, err := dto.ToDomain()
	if err != nil {
		return err
	}

	if id, _ := s.r.GetIDByEmail(ctx, *email); id == nil {
		return ErrInvalidEmail
	}

	token, err := jwtservice.CreatePasswordResetToken(*email)
	if err != nil {
		return err
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	magicLink := fmt.Sprintf("%s/login/forget-password?token=%s", frontendURL, token)

	bodyEmail := &emailservice.BodyEmail{
		Email:   *email,
		Subject: "Redefinição de senha",
		Body: `<div style="font-family: Arial, sans-serif; max-width: 480px; margin: 0 auto; background: #fff; border-radius: 8px; box-shadow: 0 2px 8px #0001; padding: 32px;">
		<h2 style="color: #eab308; margin-bottom: 16px;">Redefinição de Senha - GFood</h2>
		
		<p style="color: #333; font-size: 16px; margin-bottom: 24px;">
			Recebemos uma solicitação para redefinir a senha da sua conta.<br>
			Clique no botão abaixo para redefinir sua senha:
		</p>

		<div style="text-align: center; margin-bottom: 24px;">
			<a href="` + magicLink + `" style="display: inline-block; background-color: #eab308; color: #ffffff; font-size: 16px; font-weight: bold; padding: 12px 24px; text-decoration: none; border-radius: 6px;">
				Redefinir Senha
			</a>
		</div>
		
		<p style="color: #666; font-size: 14px;">
			Este link é exclusivo para sua conta e garante a segurança da redefinição.<br>
			Se você não solicitou essa alteração, ignore este e-mail.<br>
			<b>O link expira em 30 minutos.</b>
		</p>
		
		<p style="color: #999; font-size: 13px; margin-top: 24px;">
			Atenciosamente,<br>
			Equipe GFood
		</p>
		</div>
	`,
	}
	// Send email service
	if s.emailService != nil {
		if err := s.emailService.SendEmail(bodyEmail); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) UpdateUser(ctx context.Context, dtoID *entitydto.IDRequest, dto *companydto.UserUpdateDTO) error {
	userModel, err := s.r.GetUserByID(ctx, dtoID.ID, false)
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

	idToken, err := jwtservice.CreateBasicAccessToken(userLoggedIn)

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
	userModel, err := s.r.GetUserByID(ctx, userIDUUID, true)
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

	return jwtservice.CreateFullAccessToken(accessToken, *schema)
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
	userModel, err := s.r.GetUserByID(ctx, userIDUUID, true)
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

// ListPublicUsers returns basic data for all registered users.
func (s *Service) ListPublicUsers(ctx context.Context) ([]companydto.UserBasicDTO, error) {
	userModels, err := s.r.ListPublicUsers(ctx)
	if err != nil {
		return nil, err
	}

	if len(userModels) == 0 {
		return []companydto.UserBasicDTO{}, nil
	}

	basic := make([]companydto.UserBasicDTO, len(userModels))
	for i := range userModels {
		user := userModels[i].ToDomain()
		dto := companydto.UserBasicDTO{}
		dto.FromDomain(user)
		basic[i] = dto
	}

	return basic, nil
}

// GetUserByID returns the user data by ID.
func (s *Service) GetUserByID(ctx context.Context, dtoID *entitydto.IDRequest) (*companydto.UserDTO, error) {
	userModel, err := s.r.GetUserByID(ctx, dtoID.ID, false)
	if err != nil {
		return nil, err
	}

	user := userModel.ToDomain()
	dto := &companydto.UserDTO{}
	dto.FromDomain(user)

	return dto, nil
}
