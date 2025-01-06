package processruleusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	processruledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/process_rule"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	ErrProcessRuleIsUsed = errors.New("processRule is used in products")
)

type Service struct {
	r model.ProcessRuleRepository
}

func NewService(c model.ProcessRuleRepository) *Service {
	return &Service{r: c}
}

func (s *Service) CreateProcessRule(ctx context.Context, dto *processruledto.ProcessRuleCreateDTO) (uuid.UUID, error) {
	processRule, err := dto.ToDomain()

	if err != nil {
		return uuid.Nil, err
	}

	processRuleModel := &model.ProcessRule{}
	processRuleModel.FromDomain(processRule)
	err = s.r.CreateProcessRule(ctx, processRuleModel)

	if err != nil {
		return uuid.Nil, err
	}

	return processRule.ID, nil
}

func (s *Service) UpdateProcessRule(ctx context.Context, dtoId *entitydto.IDRequest, dto *processruledto.ProcessRuleUpdateDTO) error {
	processRuleModel, err := s.r.GetProcessRuleById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	processRule := processRuleModel.ToDomain()
	if err = dto.UpdateDomain(processRule); err != nil {
		return err
	}

	processRuleModel.FromDomain(processRule)
	if err = s.r.UpdateProcessRule(ctx, processRuleModel); err != nil {
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
	if processRuleModel, err := s.r.GetProcessRuleById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return processRuleModel.ToDomain(), nil
	}
}

func (s *Service) GetProcessRulesByCategoryId(ctx context.Context, dto *entitydto.IDRequest) ([]productentity.ProcessRule, error) {
	if processRuleModels, err := s.r.GetProcessRulesByCategoryId(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		processRules := []productentity.ProcessRule{}
		for _, processRuleModel := range processRuleModels {
			processRule := processRuleModel.ToDomain()
			processRules = append(processRules, *processRule)
		}
		return processRules, nil
	}
}

func (s *Service) GetAllProcessRules(ctx context.Context) ([]processruledto.ProcessRuleDTO, error) {
	if processRuleModels, err := s.r.GetAllProcessRules(ctx); err != nil {
		return nil, err
	} else {
		return s.processRulesToDto(processRuleModels), nil
	}
}
func (s *Service) GetAllProcessRulesWithOrderProcess(ctx context.Context) ([]processruledto.ProcessRuleDTO, error) {
	if processRuleModels, err := s.r.GetAllProcessRulesWithOrderProcess(ctx); err != nil {
		return nil, err
	} else {
		return s.processRulesToDto(processRuleModels), nil
	}
}

func (s *Service) processRulesToDto(processRuleModels []model.ProcessRule) []processruledto.ProcessRuleDTO {
	var processRuleDTOs []processruledto.ProcessRuleDTO

	for _, processRuleModel := range processRuleModels {
		processRule := processRuleModel.ToDomain()
		processRuleDTO := processruledto.ProcessRuleDTO{}
		processRuleDTO.FromDomain(processRule)
		processRuleDTOs = append(processRuleDTOs, processRuleDTO)
	}

	return processRuleDTOs
}
