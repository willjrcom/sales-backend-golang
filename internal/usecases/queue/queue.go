package queueusecases

import (
	"context"
	"time"

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
	groupItemID, startedAt, err := dto.ToModel()
	if err != nil {
		return uuid.Nil, err
	}

	queue, err := processentity.NewQueue(groupItemID, *startedAt)
	if err != nil {
		return uuid.Nil, err
	}

	if err := s.r.RegisterQueue(ctx, queue); err != nil {
		return uuid.Nil, err
	}

	return queue.ID, nil
}

func (s *Service) FinishQueue(ctx context.Context, groupItemID uuid.UUID, finishedAt time.Time) error {
	queue, err := s.r.GetQueueByGroupItemId(ctx, groupItemID.String())
	if err != nil {
		return err
	}

	queue.FinishQueue(finishedAt)

	if err := s.r.UpdateQueue(ctx, queue); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetQueueById(ctx context.Context, dto *entitydto.IdRequest) (*processentity.Queue, error) {
	if process, err := s.r.GetQueueById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return process, nil
	}
}

func (s *Service) GetQueueByGroupItemId(ctx context.Context, dto *entitydto.IdRequest) (*processentity.Queue, error) {
	if process, err := s.r.GetQueueByGroupItemId(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return process, nil
	}
}

func (s *Service) GetAllQueues(ctx context.Context) ([]processentity.Queue, error) {
	if queue, err := s.r.GetAllQueues(ctx); err != nil {
		return nil, err
	} else {
		return queue, nil
	}
}
