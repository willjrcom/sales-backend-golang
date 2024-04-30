package queueusecases

import (
	"context"

	"github.com/google/uuid"
	processentity "github.com/willjrcom/sales-backend-go/internal/domain/process"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	queuedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/queue"
)

type Service struct {
	r  processentity.QueueRepository
	rp processentity.ProcessRepository
}

func NewService(c processentity.QueueRepository, rp processentity.ProcessRepository) *Service {
	return &Service{r: c, rp: rp}
}

func (s *Service) StartQueue(ctx context.Context, dto *queuedto.StartQueueInput) (uuid.UUID, error) {
	groupItemID, joinedAt, err := dto.ToModel()
	if err != nil {
		return uuid.Nil, err
	}

	queue, err := processentity.NewQueue(groupItemID, *joinedAt)
	if err != nil {
		return uuid.Nil, err
	}

	if err := s.r.RegisterQueue(ctx, queue); err != nil {
		return uuid.Nil, err
	}

	return queue.ID, nil
}

func (s *Service) FinishQueue(ctx context.Context, process *processentity.Process) error {
	queue, err := s.r.GetOpenedQueueByGroupItemId(ctx, process.GroupItemID.String())
	if err != nil {
		return err
	}

	queue.FinishQueue(process.ProcessRuleID, *process.StartedAt)

	if err := s.r.UpdateQueue(ctx, queue); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetQueueById(ctx context.Context, dto *entitydto.IdRequest) (*processentity.Queue, error) {
	if queue, err := s.r.GetQueueById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return queue, nil
	}
}

func (s *Service) GetQueuesByGroupItemId(ctx context.Context, dto *entitydto.IdRequest) ([]processentity.Queue, error) {
	if queues, err := s.r.GetQueuesByGroupItemId(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return queues, nil
	}
}

func (s *Service) GetAllQueues(ctx context.Context) ([]processentity.Queue, error) {
	if queue, err := s.r.GetAllQueues(ctx); err != nil {
		return nil, err
	} else {
		return queue, nil
	}
}
