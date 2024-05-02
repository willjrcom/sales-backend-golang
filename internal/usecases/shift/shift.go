package shiftusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	shiftdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/shift"
)

var (
	ErrShiftAlreadyClosed = errors.New("shift already closed")
	ErrShiftAlreadyOpened = errors.New("shift already opened")
)

type Service struct {
	r shiftentity.ShiftRepository
}

func NewService(c shiftentity.ShiftRepository) *Service {
	return &Service{r: c}
}

func (s *Service) OpenShift(ctx context.Context, dto *shiftdto.OpenShift) (id uuid.UUID, err error) {
	shift, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if openedShift, _ := s.r.GetOpenedShift(ctx); openedShift != nil {
		return uuid.Nil, ErrShiftAlreadyOpened
	}

	shift.OpenShift()

	if err = s.r.CreateShift(ctx, shift); err != nil {
		return uuid.Nil, err
	}

	return shift.ID, nil
}

func (s *Service) CloseShift(ctx context.Context, dto *shiftdto.CloseShift) (err error) {
	endChange, err := dto.ToModel()

	if err != nil {
		return err
	}

	shift, err := s.r.GetOpenedShift(ctx)

	if err != nil {
		return err
	}

	if shift.IsClosed() {
		return ErrShiftAlreadyClosed
	}

	shift.CloseShift(endChange)

	if err := s.r.UpdateShift(ctx, shift); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetShiftByID(ctx context.Context, dtoID *entitydto.IdRequest) (shift *shiftentity.Shift, err error) {
	return s.r.GetShiftByID(ctx, dtoID.ID.String())
}

func (s *Service) GetOpenedShift(ctx context.Context) (shift *shiftentity.Shift, err error) {
	return s.r.GetOpenedShift(ctx)
}

func (s *Service) GetAllShifts(ctx context.Context) (shift []shiftentity.Shift, err error) {
	return s.r.GetAllShifts(ctx)
}
