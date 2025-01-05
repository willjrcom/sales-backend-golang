package processruleusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	processruledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/process_rule"
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

func (s *Service) CreateProcessRule(ctx context.Context, dto *processruledto.CreateProcessRuleInput) (uuid.UUID, error) {
	processRule, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	err = s.r.CreateProcessRule(ctx, processRule)

	if err != nil {
		return uuid.Nil, err
	}

	return processRule.ID, nil
}

func (s *Service) UpdateProcessRule(ctx context.Context, dtoId *entitydto.IdRequest, dto *processruledto.UpdateProcessRuleInput) error {
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

func (s *Service) GetProcessRulesByCategoryId(ctx context.Context, dto *entitydto.IdRequest) ([]productentity.ProcessRule, error) {
	if processRules, err := s.r.GetProcessRulesByCategoryId(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return processRules, nil
	}
}

func (s *Service) GetAllProcessRules(ctx context.Context) ([]processruledto.ProcessRuleOutput, error) {
	if processRules, err := s.r.GetAllProcessRules(ctx); err != nil {
		return nil, err
	} else {
		return s.processRulesToDto(processRules), nil
	}
}
func (s *Service) GetAllProcessRulesWithOrderProcess(ctx context.Context) ([]processruledto.ProcessRuleOutput, error) {
	if processRules, err := s.r.GetAllProcessRulesWithOrderProcess(ctx); err != nil {
		return nil, err
	} else {
		return s.processRulesToDto(processRules), nil
	}
}

func (s *Service) processRulesToDto(processRules []productentity.ProcessRule) []processruledto.ProcessRuleOutput {
	var processRuleDtos []processruledto.ProcessRuleOutput

	for _, processRule := range processRules {
		processRuleOutput := processruledto.ProcessRuleOutput{}
		processRuleOutput.FromModel(&processRule)
		processRuleDtos = append(processRuleDtos, processRuleOutput)
	}

	return processRuleDtos
}
