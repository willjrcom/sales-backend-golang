package advertisingusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	advertisingdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/advertising"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	ErrAdvertisingNotFound = errors.New("advertising not found")
)

type AdvertisingService struct {
	repo        model.AdvertisingRepository
	sponsorRepo model.SponsorRepository
	userRepo    model.UserRepository
}

func NewAdvertisingService(repo model.AdvertisingRepository, sponsorRepo model.SponsorRepository, userRepo model.UserRepository) *AdvertisingService {
	return &AdvertisingService{
		repo:        repo,
		sponsorRepo: sponsorRepo,
		userRepo:    userRepo,
	}
}

func (s *AdvertisingService) CreateAdvertising(ctx context.Context, dto *advertisingdto.CreateAdvertisingDTO) (uuid.UUID, error) {
	advertising, err := dto.ToDomain()
	if err != nil {
		return uuid.Nil, err
	}

	sponsorModel, err := s.sponsorRepo.GetByID(ctx, advertising.SponsorID)
	if err != nil {
		return uuid.Nil, errors.New("sponsor not found")
	}
	advertising.Sponsor = sponsorModel.ToDomain()

	advertisingModel := &model.Advertising{}
	advertisingModel.FromDomain(advertising)

	err = s.repo.Create(ctx, advertisingModel)
	if err != nil {
		return uuid.Nil, err
	}

	return advertisingModel.ID, nil
}

func (s *AdvertisingService) UpdateAdvertising(ctx context.Context, idDto *entitydto.IDRequest, dto *advertisingdto.UpdateAdvertisingDTO) error {
	advertisingModel, err := s.repo.GetByID(ctx, idDto.ID)
	if err != nil {
		return err
	}

	advertising := advertisingModel.ToDomain()

	if err := dto.UpdateDomain(advertising); err != nil {
		return err
	}

	if dto.SponsorID != nil {
		sponsorModel, err := s.sponsorRepo.GetByID(ctx, advertising.SponsorID)
		if err != nil {
			return errors.New("new sponsor not found")
		}
		advertising.Sponsor = sponsorModel.ToDomain()
	}

	advertisingModel.FromDomain(advertising)
	return s.repo.Update(ctx, advertisingModel)
}

func (s *AdvertisingService) DeleteAdvertising(ctx context.Context, idDto *entitydto.IDRequest) error {
	return s.repo.Delete(ctx, idDto.ID)
}

func (s *AdvertisingService) GetAdvertisingById(ctx context.Context, idDto *entitydto.IDRequest) (*advertisingdto.AdvertisingDTO, error) {
	advertisingModel, err := s.repo.GetByID(ctx, idDto.ID)
	if err != nil {
		return nil, err
	}

	dto := &advertisingdto.AdvertisingDTO{}
	dto.FromDomain(advertisingModel.ToDomain())
	return dto, nil
}

func (s *AdvertisingService) GetAllAdvertisements(ctx context.Context) ([]advertisingdto.AdvertisingDTO, error) {
	models, err := s.repo.GetAllAdvertisements(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]advertisingdto.AdvertisingDTO, len(models))
	for i, m := range models {
		dtos[i].FromDomain(m.ToDomain())
	}

	return dtos, nil
}

func (s *AdvertisingService) GetActiveAdvertisements(ctx context.Context) ([]advertisingdto.AdvertisingDTO, error) {
	userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)
	if !ok {
		return nil, errors.New("user not found in context")
	}

	userIDUUID := uuid.MustParse(userID)
	userModel, err := s.userRepo.GetUserByID(ctx, userIDUUID, true)
	if err != nil {
		return nil, err
	}

	categoryIDsMap := make(map[uuid.UUID]struct{})
	for _, company := range userModel.Companies {
		for _, category := range company.Categories {
			categoryIDsMap[category.ID] = struct{}{}
		}
	}

	categoryIDs := make([]uuid.UUID, 0, len(categoryIDsMap))
	for id := range categoryIDsMap {
		categoryIDs = append(categoryIDs, id)
	}

	models, err := s.repo.GetActiveAdvertisements(ctx, categoryIDs)
	if err != nil {
		return nil, err
	}

	dtos := make([]advertisingdto.AdvertisingDTO, len(models))
	for i, m := range models {
		dtos[i].FromDomain(m.ToDomain())
	}

	return dtos, nil
}
