package orderqueueusecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderqueuedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_queue"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type Service struct {
	r  model.QueueRepository
	rp model.OrderProcessRepository
}

func NewService(c model.QueueRepository) *Service {
	return &Service{r: c}
}

func (s *Service) AddDependencies(rp model.OrderProcessRepository) {
	s.rp = rp
}

func (s *Service) StartQueue(ctx context.Context, dto *orderqueuedto.QueueCreateDTO) (uuid.UUID, error) {
	groupItemID, joinedAt, err := dto.ToDomain()
	if err != nil {
		return uuid.Nil, err
	}

	queue, err := orderprocessentity.NewOrderQueue(groupItemID, *joinedAt)
	if err != nil {
		return uuid.Nil, err
	}

	queueModel := &model.OrderQueue{}
	queueModel.FromDomain(queue)
	if err := s.r.CreateQueue(ctx, queueModel); err != nil {
		return uuid.Nil, err
	}

	return queue.ID, nil
}

func (s *Service) FinishQueue(ctx context.Context, process *orderprocessentity.OrderProcess) error {
	queueModel, err := s.r.GetOpenedQueueByGroupItemId(ctx, process.GroupItemID.String())
	if err != nil {
		return fmt.Errorf("opened queue not found")
	}

	queue := queueModel.ToDomain()
	queue.Finish(process.ProcessRuleID, *process.StartedAt)

	queueModel.FromDomain(queue)
	if err := s.r.UpdateQueue(ctx, queueModel); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetQueueById(ctx context.Context, dto *entitydto.IDRequest) (*orderqueuedto.QueueDTO, error) {
	if queueModel, err := s.r.GetQueueById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		queue := queueModel.ToDomain()

		queueDto := &orderqueuedto.QueueDTO{}
		queueDto.FromDomain(queue)
		return queueDto, nil
	}
}

func (s *Service) GetQueuesByGroupItemId(ctx context.Context, dto *entitydto.IDRequest) ([]orderqueuedto.QueueDTO, error) {
	if queueModels, err := s.r.GetQueuesByGroupItemId(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		queues := []orderqueuedto.QueueDTO{}
		for _, queueModel := range queueModels {
			queue := queueModel.ToDomain()

			queueDto := &orderqueuedto.QueueDTO{}
			queueDto.FromDomain(queue)
			queues = append(queues, *queueDto)
		}
		return queues, nil
	}
}

func (s *Service) GetAllQueues(ctx context.Context) ([]orderqueuedto.QueueDTO, error) {
	if queueModels, err := s.r.GetAllQueues(ctx); err != nil {
		return nil, err
	} else {
		queues := []orderqueuedto.QueueDTO{}
		for _, queueModel := range queueModels {
			queue := queueModel.ToDomain()

			queueDto := &orderqueuedto.QueueDTO{}
			queueDto.FromDomain(queue)
			queues = append(queues, *queueDto)
		}
		return queues, nil
	}
}
