package processusecases

import (
	"context"

	"github.com/google/uuid"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	processentity "github.com/willjrcom/sales-backend-go/internal/domain/process"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	processdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/process"
)

type Service struct {
	r  processentity.ProcessRepository
	ri itementity.ItemRepository
}

func NewService(c processentity.ProcessRepository, ri itementity.ItemRepository) *Service {
	return &Service{r: c, ri: ri}
}

func (s *Service) CreateProcess(ctx context.Context, dto *processdto.CreateProcessInput) (uuid.UUID, error) {
	process, err := dto.ToModel()
	if err != nil {
		return uuid.Nil, err
	}

	item, err := s.ri.GetItemById(ctx, process.ItemID.String())
	if err != nil {
		return uuid.Nil, err
	}

	process.ProductID = item.ProductID

	if err := s.r.RegisterProcess(ctx, process); err != nil {
		return uuid.Nil, err
	}

	return process.ID, nil
}

func (s *Service) StartProcess(ctx context.Context, dtoID *entitydto.IdRequest, dto *processdto.StartProcessInput) error {
	employeeID, err := dto.ToModel()
	if err != nil {
		return err
	}

	process, err := s.r.GetProcessById(ctx, dtoID.ID.String())
	if err != nil {
		return err
	}

	if err := process.StartProcess(employeeID); err != nil {
		return err
	}

	if err := s.r.UpdateProcess(ctx, process); err != nil {
		return err
	}

	return nil
}

func (s *Service) PauseProcess(ctx context.Context, dtoID *entitydto.IdRequest) error {
	process, err := s.r.GetProcessById(ctx, dtoID.ID.String())
	if err != nil {
		return err
	}

	if err := process.PauseProcess(); err != nil {
		return err
	}

	if err := s.r.UpdateProcess(ctx, process); err != nil {
		return err
	}

	return nil
}

func (s *Service) ContinueProcess(ctx context.Context, dtoID *entitydto.IdRequest) error {
	process, err := s.r.GetProcessById(ctx, dtoID.ID.String())
	if err != nil {
		return err
	}

	if err := process.ContinueProcess(); err != nil {
		return err
	}

	if err := s.r.UpdateProcess(ctx, process); err != nil {
		return err
	}

	return nil
}

func (s *Service) FinishProcess(ctx context.Context, dtoID *entitydto.IdRequest) error {
	process, err := s.r.GetProcessById(ctx, dtoID.ID.String())
	if err != nil {
		return err
	}

	if err := process.FinishProcess(); err != nil {
		return err
	}

	if err := s.r.UpdateProcess(ctx, process); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetProcessById(ctx context.Context, dto *entitydto.IdRequest) (*processdto.ProcessOutput, error) {
	if process, err := s.r.GetProcessById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		processOutput := &processdto.ProcessOutput{}
		processOutput.FromModel(process)
		return processOutput, nil
	}
}

func (s *Service) GetAllProcesses(ctx context.Context) ([]processdto.ProcessOutput, error) {
	if process, err := s.r.GetAllProcesses(ctx); err != nil {
		return nil, err
	} else {
		return s.processesToOutputs(process), nil
	}
}

func (s *Service) processesToOutputs(processes []processentity.Process) []processdto.ProcessOutput {
	outputs := make([]processdto.ProcessOutput, 0)
	for _, process := range processes {
		processOutput := &processdto.ProcessOutput{}
		processOutput.FromModel(&process)
		outputs = append(outputs, *processOutput)
	}

	return outputs
}
