package processusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	processentity "github.com/willjrcom/sales-backend-go/internal/domain/process"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	processdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/process"
)

var (
	ErrProcessIsUsed = errors.New("process is used in products")
)

type Service struct {
	r processentity.ProcessRepository
}

func NewService(c processentity.ProcessRepository) *Service {
	return &Service{r: c}
}

func (s *Service) RegisterProcess(ctx context.Context, dto *processdto.CreateProcessInput) (uuid.UUID, error) {
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

func (s *Service) UpdateProcess(ctx context.Context, dtoId *entitydto.IdRequest) error {
	process, err := s.r.GetProcessById(ctx, dtoId.ID.String())

	if err != nil {
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

func (s *Service) GetProcessById(ctx context.Context, dto *entitydto.IdRequest) (*processentity.Process, error) {
	if process, err := s.r.GetProcessById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return process, nil
	}
}

func (s *Service) GetAllProcesses(ctx context.Context, dto *entitydto.IdRequest) ([]processentity.Process, error) {
	if process, err := s.r.GetAllProcesses(ctx); err != nil {
		return nil, err
	} else {
		return process, nil
	}
}
