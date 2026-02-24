package sponsorusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	sponsordto "github.com/willjrcom/sales-backend-go/internal/infra/dto/sponsor"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	ErrSponsorNotFound = errors.New("sponsor not found")
)

type SponsorService struct {
	repo model.SponsorRepository
}

func NewSponsorService(repo model.SponsorRepository) *SponsorService {
	return &SponsorService{repo: repo}
}

func (s *SponsorService) CreateSponsor(ctx context.Context, dto *sponsordto.CreateSponsorDTO) (uuid.UUID, error) {
	sponsor, err := dto.ToDomain()
	if err != nil {
		return uuid.Nil, err
	}

	sponsorModel := &model.Sponsor{}
	sponsorModel.FromDomain(sponsor)

	err = s.repo.Create(ctx, sponsorModel)
	if err != nil {
		return uuid.Nil, err
	}

	return sponsorModel.ID, nil
}

func (s *SponsorService) UpdateSponsor(ctx context.Context, idDto *entitydto.IDRequest, dto *sponsordto.UpdateSponsorDTO) error {
	sponsorModel, err := s.repo.GetByID(ctx, idDto.ID)
	if err != nil {
		return err
	}

	sponsor := sponsorModel.ToDomain()

	if err := dto.UpdateDomain(sponsor); err != nil {
		return err
	}

	sponsorModel.FromDomain(sponsor)
	return s.repo.Update(ctx, sponsorModel)
}

func (s *SponsorService) DeleteSponsor(ctx context.Context, idDto *entitydto.IDRequest) error {
	return s.repo.Delete(ctx, idDto.ID)
}

func (s *SponsorService) GetSponsorById(ctx context.Context, idDto *entitydto.IDRequest) (*sponsordto.SponsorDTO, error) {
	sponsorModel, err := s.repo.GetByID(ctx, idDto.ID)
	if err != nil {
		return nil, err
	}

	dto := &sponsordto.SponsorDTO{}
	dto.FromDomain(sponsorModel.ToDomain())
	return dto, nil
}

func (s *SponsorService) GetAllSponsors(ctx context.Context) ([]sponsordto.SponsorDTO, error) {
	models, err := s.repo.GetAllSponsors(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]sponsordto.SponsorDTO, len(models))
	for i, m := range models {
		dtos[i].FromDomain(m.ToDomain())
	}

	return dtos, nil
}
