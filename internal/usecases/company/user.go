package companyusecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
	emailservice "github.com/willjrcom/sales-backend-go/internal/infra/service/email"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/rabbitmq"
)

func (s *Service) AddUserToCompany(ctx context.Context, dto *companydto.UserToCompanyDTO) error {
	email, err := dto.ToDomain()

	if err != nil {
		return err
	}

	userID, _ := s.u.GetIDByEmail(ctx, email)

	if userID != nil {
		if exists, _ := s.r.ValidateUserToPublicCompany(ctx, *userID); exists {
			return errors.New("user already added to company")
		}
	}

	if userID == nil {
		return errors.New("user not found")
	}

	if err := s.r.AddUserToPublicCompany(ctx, *userID); err != nil {
		return err
	}

	// send email
	if s.rabbitmq == nil {
		return nil
	}

	company, err := s.GetCompany(ctx)
	if err != nil {
		return err
	}

	bodyEmail := &emailservice.BodyEmail{
		Email:   email,
		Subject: "Bem-vindo à GFood",
		Body: `<div style="font-family: Arial, sans-serif; max-width: 480px; margin: 0 auto; background: #fff; border-radius: 8px; box-shadow: 0 2px 8px #0001; padding: 32px;">
			<h2 style="color: #eab308; margin-bottom: 16px;">Você foi adicionado à GFood</h2>
			
			<p style="color: #333; font-size: 16px; margin-bottom: 24px;">
				Sua conta (` + email + `) foi adicionada com sucesso ao painel da empresa: ` + company.TradeName + `.<br>
				Acesse agora para começar a gerenciar seus pedidos:
			</p>

			<div style="text-align: center; margin-bottom: 24px;">
				<a href="https://gfood.com.br/login" style="display: inline-block; background-color: #eab308; color: #ffffff; font-size: 16px; font-weight: bold; padding: 12px 24px; text-decoration: none; border-radius: 6px;">
					Acessar Painel
				</a>
			</div>
			
			<p style="color: #999; font-size: 13px; margin-top: 24px;">
				Atenciosamente,<br>
				Equipe GFood
			</p>
			</div>
		`,
	}

	bodyJSON, err := json.Marshal(bodyEmail)
	if err != nil {
		fmt.Println("Error marshaling email body:", err)
		return nil
	}

	if err := s.rabbitmq.SendMessage(rabbitmq.EMAIL_EX, "", string(bodyJSON)); err != nil {
		fmt.Println("Error sending email:", err)
	}

	return nil
}

func (s *Service) RemoveUserFromCompany(ctx context.Context, dto *companydto.UserToCompanyDTO) error {
	email, err := dto.ToDomain()

	if err != nil {
		return err
	}

	userID, err := s.u.GetIDByEmail(ctx, email)

	if err != nil {
		return err
	}

	if userID != nil {
		if exists, _ := s.r.ValidateUserToPublicCompany(ctx, *userID); !exists {
			return errors.New("user already removed from company")
		}
	}

	if err := s.r.RemoveUserFromPublicCompany(ctx, *userID); err != nil {
		return err
	}

	// send email
	if s.rabbitmq == nil {
		return nil
	}

	company, err := s.GetCompany(ctx)
	if err != nil {
		return err
	}

	bodyEmail := &emailservice.BodyEmail{
		Email:   email,
		Subject: "Acesso removido - GFood",
		Body: `<div style="font-family: Arial, sans-serif; max-width: 480px; margin: 0 auto; background: #fff; border-radius: 8px; box-shadow: 0 2px 8px #0001; padding: 32px;">
			<h2 style="color: #f87171; margin-bottom: 16px;">Acesso Removido</h2>
			
			<p style="color: #333; font-size: 16px; margin-bottom: 24px;">
				Olá,<br>
				Seu acesso ao painel da empresa: ` + company.TradeName + ` associado ao e-mail <b>` + email + `</b> foi removido.
			</p>
			
			<p style="color: #999; font-size: 13px; margin-top: 24px;">
				Atenciosamente,<br>
				Equipe GFood
			</p>
			</div>
		`,
	}

	bodyJSON, err := json.Marshal(bodyEmail)
	if err != nil {
		fmt.Println("Error marshaling email body:", err)
		return nil
	}

	if err := s.rabbitmq.SendMessage(rabbitmq.EMAIL_EX, "", string(bodyJSON)); err != nil {
		fmt.Println("Error sending email:", err)
	}

	return nil
}

// GetCompanyUsers retrieves a paginated list of users and the total count.
func (s *Service) GetCompanyUsers(ctx context.Context, page, perPage int) ([]companydto.UserDTO, int, error) {
	userModels, total, err := s.r.GetCompanyUsers(ctx, page, perPage)
	if err != nil {
		return nil, 0, err
	}
	dtos := make([]companydto.UserDTO, len(userModels))
	for i, userModel := range userModels {
		user := userModel.ToDomain()
		dto := &companydto.UserDTO{}
		dto.FromDomain(user)
		dtos[i] = *dto
	}
	return dtos, total, nil
}
