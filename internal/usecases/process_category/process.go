package processusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	processdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/process_category"
)

var (
	ErrProcessIsUsed = errors.New("process is used in products")
)

type Service struct {
	r productentity.ProcessRepository
}

func NewService(c productentity.ProcessRepository) *Service {
	return &Service{r: c}
}

func (s *Service) RegisterProcess(ctx context.Context, dto *processdto.RegisterProcessInput) (uuid.UUID, error) {
	process, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	err = s.r.RegisterProcess(ctx, process)

	if err != nil {
		return uuid.Nil, err
	}

	return process.ID, nil
}

func (s *Service) UpdateProcess(ctx context.Context, dtoId *entitydto.IdRequest, dto *processdto.UpdateProcessInput) error {
	process, err := s.r.GetProcessById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err = dto.UpdateModel(process); err != nil {
		return err
	}

	if err = s.r.UpdateProcess(ctx, process); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteProcess(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.r.GetProcessById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.r.DeleteProcess(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetProcessById(ctx context.Context, dto *entitydto.IdRequest) (*productentity.Process, error) {
	if process, err := s.r.GetProcessById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return process, nil
	}
}
