package companyusecases

import (
	"context"
	"errors"

	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
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
		// createUserInput := &companydto.UserCreateDTO{
		// 	Email:            email,
		// 	GeneratePassword: true,
		// }

		// if newUserID, err := s.us.CreateUser(ctx, createUserInput); err != nil {
		// 	return err
		// } else {
		// 	userID = newUserID
		// }
	}

	if err := s.r.AddUserToPublicCompany(ctx, *userID); err != nil {
		return err
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
