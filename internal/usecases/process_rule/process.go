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

func (s *Service) CreateProcessRule(ctx context.Context, dto *processruledto.ProcessRuleCreateDTO) (uuid.UUID, error) {
	processRule, err := dto.ToDomain()

	if err != nil {
		return uuid.Nil, err
	}

	err = s.r.CreateProcessRule(ctx, processRule)

	if err != nil {
		return uuid.Nil, err
	}

	return processRule.ID, nil
}

func (s *Service) UpdateProcessRule(ctx context.Context, dtoId *entitydto.IDRequest, dto *processruledto.ProcessRuleUpdateDTO) error {
	processRule, err := s.r.GetProcessRuleById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err = dto.UpdateDomain(processRule); err != nil {
		return err
	}

	if err = s.r.UpdateProcessRule(ctx, processRule); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteProcessRule(ctx context.Context, dto *entitydto.IDRequest) error {
	if _, err := s.r.GetProcessRuleById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.r.DeleteProcessRule(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetProcessRuleById(ctx context.Context, dto *entitydto.IDRequest) (*productentity.ProcessRule, error) {
	if processRule, err := s.r.GetProcessRuleById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return processRule, nil
	}
}

func (s *Service) GetProcessRulesByCategoryId(ctx context.Context, dto *entitydto.IDRequest) ([]productentity.ProcessRule, error) {
	if processRules, err := s.r.GetProcessRulesByCategoryId(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return processRules, nil
	}
}

func (s *Service) GetAllProcessRules(ctx context.Context) ([]processruledto.ProcessRuleDTO, error) {
	if processRules, err := s.r.GetAllProcessRules(ctx); err != nil {
		return nil, err
	} else {
		return s.processRulesToDto(processRules), nil
	}
}
func (s *Service) GetAllProcessRulesWithOrderProcess(ctx context.Context) ([]processruledto.ProcessRuleDTO, error) {
	if processRules, err := s.r.GetAllProcessRulesWithOrderProcess(ctx); err != nil {
		return nil, err
	} else {
		return s.processRulesToDto(processRules), nil
	}
}

func (s *Service) processRulesToDto(processRules []productentity.ProcessRule) []processruledto.ProcessRuleDTO {
	var processRuleDTOs []processruledto.ProcessRuleDTO

	for _, processRule := range processRules {
		processRuleDTO := processruledto.ProcessRuleDTO{}
		processRuleDTO.FromDomain(&processRule)
		processRuleDTOs = append(processRuleDTOs, processRuleDTO)
	}

	return processRuleDTOs
}
