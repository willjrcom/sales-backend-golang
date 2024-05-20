package orderprocessusecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderprocessdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_process"
	orderqueuedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_queue"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/kafka"
	groupitemusecases "github.com/willjrcom/sales-backend-go/internal/usecases/group_item"
	orderqueueusecases "github.com/willjrcom/sales-backend-go/internal/usecases/order_queue"
)

type Service struct {
	r             orderprocessentity.ProcessRepository
	rpr           productentity.ProcessRuleRepository
	sq            *orderqueueusecases.Service
	rsgi          *groupitemusecases.Service
	producerKafka *kafka.KafkaProducer
}

func NewService(c orderprocessentity.ProcessRepository) *Service {
	return &Service{r: c}
}

func (s *Service) AddDependencies(sq *orderqueueusecases.Service, rpr productentity.ProcessRuleRepository, rsgi *groupitemusecases.Service, producerKafka *kafka.KafkaProducer) {
	s.rpr = rpr
	s.sq = sq
	s.rsgi = rsgi
	s.producerKafka = producerKafka
}

func (s *Service) CreateProcess(ctx context.Context, dto *orderprocessdto.CreateProcessInput) (uuid.UUID, error) {
	process, err := dto.ToModel()
	if err != nil {
		return uuid.Nil, err
	}

	idRequest := entitydto.NewIdRequest(process.GroupItemID)
	groupItem, err := s.rsgi.GetGroupByID(ctx, idRequest)
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

	if err := s.r.CreateProcess(ctx, process); err != nil {
		return uuid.Nil, err
	}

	if err := s.producerKafka.NewMessage(ctx, "order_process", process); err != nil {
		return uuid.Nil, err
	}

	return process.ID, nil
}

func (s *Service) StartProcess(ctx context.Context, dtoID *entitydto.IdRequest, dto *orderprocessdto.StartProcessInput) error {
	employeeID, err := dto.ToModel()
	if err != nil {
		return err
	}

	process, err := s.r.GetProcessById(ctx, dtoID.ID.String())
	if err != nil {
		return err
	}

	processRule, err := s.rpr.GetProcessRuleById(ctx, process.ProcessRuleID.String())
	if err != nil {
		return err
	}

	if processRule.Order == 1 {
		entityDtoID := entitydto.NewIdRequest(process.GroupItemID)
		if err := s.rsgi.StartGroupItem(ctx, entityDtoID); err != nil {
			return err
		}
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
		entityDtoID := &entitydto.IdRequest{ID: process.GroupItemID}
		if err := s.rsgi.ReadyGroupItem(ctx, entityDtoID); err != nil {
			return uuid.Nil, err
		}

		return uuid.Nil, nil
	}

	startQueueInput := &orderqueuedto.StartQueueInput{
		OrderQueueCommonAttributes: orderprocessentity.OrderQueueCommonAttributes{
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
	createProcessInput := &orderprocessdto.CreateProcessInput{
		OrderProcessCommonAttributes: orderprocessentity.OrderProcessCommonAttributes{
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

func (s *Service) GetProcessById(ctx context.Context, dto *entitydto.IdRequest) (*orderprocessdto.ProcessOutput, error) {
	if process, err := s.r.GetProcessById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		processOutput := &orderprocessdto.ProcessOutput{}
		processOutput.FromModel(process)
		return processOutput, nil
	}
}

func (s *Service) GetAllProcesses(ctx context.Context) ([]orderprocessdto.ProcessOutput, error) {
	if process, err := s.r.GetAllProcesses(ctx); err != nil {
		return nil, err
	} else {
		return s.processesToOutputs(process), nil
	}
}

func (s *Service) GetProcessesByProductID(ctx context.Context, dtoID *entitydto.IdRequest) ([]orderprocessdto.ProcessOutput, error) {
	if process, err := s.r.GetProcessesByProductID(ctx, dtoID.ID.String()); err != nil {
		return nil, err
	} else {
		return s.processesToOutputs(process), nil
	}
}

func (s *Service) GetProcessesByGroupItemID(ctx context.Context, dtoID *entitydto.IdRequest) ([]orderprocessdto.ProcessOutput, error) {
	if process, err := s.r.GetProcessesByGroupItemID(ctx, dtoID.ID.String()); err != nil {
		return nil, err
	} else {
		return s.processesToOutputs(process), nil
	}
}

func (s *Service) processesToOutputs(processes []orderprocessentity.OrderProcess) []orderprocessdto.ProcessOutput {
	outputs := make([]orderprocessdto.ProcessOutput, 0)
	for _, process := range processes {
		processOutput := &orderprocessdto.ProcessOutput{}
		processOutput.FromModel(&process)
		outputs = append(outputs, *processOutput)
	}

	return outputs
}
