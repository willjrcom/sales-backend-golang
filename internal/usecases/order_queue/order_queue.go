package orderqueueusecases

import (
	"context"

	"github.com/google/uuid"
	orderprocessentity "github.com/willjrcom/sales-backend-go/internal/domain/order_process"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	orderqueuedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_queue"
)

type Service struct {
	r  orderprocessentity.QueueRepository
	rp orderprocessentity.ProcessRepository
}

func NewService(c orderprocessentity.QueueRepository) *Service {
	return &Service{r: c}
}

func (s *Service) AddDependencies(rp orderprocessentity.ProcessRepository) {
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

	if err := s.r.CreateQueue(ctx, queue); err != nil {
		return uuid.Nil, err
	}

	return queue.ID, nil
}

func (s *Service) FinishQueue(ctx context.Context, process *orderprocessentity.OrderProcess) error {
	queue, err := s.r.GetOpenedQueueByGroupItemId(ctx, process.GroupItemID.String())
	if err != nil {
		return err
	}

	queue.Finish(process.ProcessRuleID, *process.StartedAt)

	if err := s.r.UpdateQueue(ctx, queue); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetQueueById(ctx context.Context, dto *entitydto.IDRequest) (*orderprocessentity.OrderQueue, error) {
	if queue, err := s.r.GetQueueById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return queue, nil
	}
}

func (s *Service) GetQueuesByGroupItemId(ctx context.Context, dto *entitydto.IDRequest) ([]orderprocessentity.OrderQueue, error) {
	if queues, err := s.r.GetQueuesByGroupItemId(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return queues, nil
	}
}

func (s *Service) GetAllQueues(ctx context.Context) ([]orderprocessentity.OrderQueue, error) {
	if queue, err := s.r.GetAllQueues(ctx); err != nil {
		return nil, err
	} else {
		return queue, nil
	}
}
