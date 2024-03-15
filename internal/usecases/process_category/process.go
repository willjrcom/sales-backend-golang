package processRuleusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	processdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/processruledto"
)

var (
	ErrProcessRuleIsUsed = errors.New("processRule is used in products")
)

type Service struct {
	r productentity.ProcessRuleRepository
}

func NewService(c productentity.ProcessRuleRepository) *Service {
	return &Service{r: c}
}

func (s *Service) RegisterProcessRule(ctx context.Context, dto *processdto.RegisterProcessRuleInput) (uuid.UUID, error) {
	processRule, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	err = s.r.RegisterProcessRule(ctx, processRule)

	if err != nil {
		return uuid.Nil, err
	}

	return processRule.ID, nil
}

func (s *Service) UpdateProcessRule(ctx context.Context, dtoId *entitydto.IdRequest, dto *processdto.UpdateProcessRuleInput) error {
	processRule, err := s.r.GetProcessRuleById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err = dto.UpdateModel(processRule); err != nil {
		return err
	}

	if err = s.r.UpdateProcessRule(ctx, processRule); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteProcessRule(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.r.GetProcessRuleById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.r.DeleteProcessRule(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetProcessRuleById(ctx context.Context, dto *entitydto.IdRequest) (*productentity.ProcessRule, error) {
	if processRule, err := s.r.GetProcessRuleById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return processRule, nil
	}
}
