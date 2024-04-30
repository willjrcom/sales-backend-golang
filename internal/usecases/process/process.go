package processusecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	processentity "github.com/willjrcom/sales-backend-go/internal/domain/process"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	processdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/process"
	queuedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/queue"
	queueusecases "github.com/willjrcom/sales-backend-go/internal/usecases/queue"
)

type Service struct {
	r   processentity.ProcessRepository
	ri  groupitementity.GroupItemRepository
	rpr productentity.ProcessRuleRepository
	sq  *queueusecases.Service
}

func NewService(c processentity.ProcessRepository, ri groupitementity.GroupItemRepository, sq *queueusecases.Service, rpr productentity.ProcessRuleRepository) *Service {
	return &Service{r: c, ri: ri, sq: sq, rpr: rpr}
}

func (s *Service) CreateProcess(ctx context.Context, dto *processdto.CreateProcessInput) (uuid.UUID, error) {
	process, err := dto.ToModel()
	if err != nil {
		return uuid.Nil, err
	}

	groupItem, err := s.ri.GetGroupByID(ctx, process.GroupItemID.String(), true)
	if err != nil {
		return uuid.Nil, err
	}

	productIDs, err := groupItem.GetDistinctProductIDs()
	if err != nil {
		return uuid.Nil, err
	}

	for _, productID := range productIDs {
		process.Products = append(process.Products, productentity.Product{
			Entity: entity.Entity{
				ID: productID,
			},
		})
	}

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

	if err := s.sq.FinishQueue(ctx, process); err != nil {
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

func (s *Service) FinishProcess(ctx context.Context, dtoID *entitydto.IdRequest) (nextProcessID uuid.UUID, err error) {
	process, err := s.r.GetProcessById(ctx, dtoID.ID.String())
	if err != nil {
		return uuid.Nil, err
	}

	if err := process.FinishProcess(); err != nil {
		return uuid.Nil, err
	}

	if err := s.r.UpdateProcess(ctx, process); err != nil {
		return uuid.Nil, err
	}

	last, err := s.rpr.IsLastProcessRuleByID(ctx, process.ProcessRuleID)
	if err != nil {
		return uuid.Nil, err
	}

	// Processes finished
	if last {
		return uuid.Nil, nil
	}

	startQueueInput := &queuedto.StartQueueInput{
		QueueCommonAttributes: processentity.QueueCommonAttributes{
			GroupItemID: process.GroupItemID,
		},
		JoinedAt: *process.FinishedAt,
	}

	if _, err := s.sq.StartQueue(ctx, startQueueInput); err != nil {
		return uuid.Nil, err
	}

	processRule, err := s.rpr.GetProcessRuleById(ctx, process.ProcessRuleID.String())
	if err != nil {
		return uuid.Nil, err
	}

	nextProcessRule, err := s.rpr.GetProcessRuleByCategoryIdAndOrder(ctx, processRule.CategoryID.String(), processRule.Order+1)
	if err != nil {
		return uuid.Nil, err
	}

	// Create next process
	createProcessInput := &processdto.CreateProcessInput{
		ProcessCommonAttributes: processentity.ProcessCommonAttributes{
			GroupItemID:   process.GroupItemID,
			ProcessRuleID: nextProcessRule.ID,
		},
	}

	nextProcessID, err = s.CreateProcess(ctx, createProcessInput)
	if err != nil {
		return uuid.Nil, err
	}

	return nextProcessID, nil
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

func (s *Service) GetProcessesByProductID(ctx context.Context, dtoID *entitydto.IdRequest) ([]processdto.ProcessOutput, error) {
	if process, err := s.r.GetProcessesByProductID(ctx, dtoID.ID.String()); err != nil {
		return nil, err
	} else {
		return s.processesToOutputs(process), nil
	}
}

func (s *Service) GetProcessesByGroupItemID(ctx context.Context, dtoID *entitydto.IdRequest) ([]processdto.ProcessOutput, error) {
	if process, err := s.r.GetProcessesByGroupItemID(ctx, dtoID.ID.String()); err != nil {
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
