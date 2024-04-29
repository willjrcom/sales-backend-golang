package queueusecases

import (
	"context"

	"github.com/google/uuid"
	processentity "github.com/willjrcom/sales-backend-go/internal/domain/process"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

type Service struct {
	r  processentity.QueueRepository
	rp processentity.ProcessRepository
}

func NewService(c processentity.QueueRepository, rp processentity.ProcessRepository) *Service {
	return &Service{r: c, rp: rp}
}

func (s *Service) JoinQueue(ctx context.Context, previousProcess *processentity.Process) (uuid.UUID, error) {
	queue, err := processentity.NewQueue(previousProcess)
	if err != nil {
		return uuid.Nil, err
	}

	if err := s.r.RegisterQueue(ctx, queue); err != nil {
		return uuid.Nil, err
	}

	return queue.ID, nil
}

func (s *Service) LeftQueue(ctx context.Context, previousProcess *processentity.Process, nextProcess *processentity.Process) error {
	queue, err := s.r.GetQueueByPreviousProcessId(ctx, previousProcess.ID.String())
	if err != nil {
		return err
	}

	if err := queue.LeftQueue(nextProcess); err != nil {
		return err
	}

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

func (s *Service) GetAllQueuees(ctx context.Context) ([]processentity.Queue, error) {
	if queue, err := s.r.GetAllQueuees(ctx); err != nil {
		return nil, err
	} else {
		return queue, nil
	}
}
