package processruleusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	productcategorydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/product_category"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	entitymodel "github.com/willjrcom/sales-backend-go/internal/infra/repository/model/entity"
)

var (
	ErrProcessRuleIsUsed = errors.New("processRule is used in products")
)

type Service struct {
	r  model.ProcessRuleRepository
	rc model.CategoryRepository
}

func NewService(c model.ProcessRuleRepository) *Service {
	return &Service{r: c}
}

func (s *Service) AddDependencies(rc model.CategoryRepository) {
	s.rc = rc
}

func (s *Service) CreateProcessRule(ctx context.Context, dto *productcategorydto.ProcessRuleCreateDTO) (uuid.UUID, error) {
	processRule, err := dto.ToDomain()

	if err != nil {
		return uuid.Nil, err
	}

	// Validate Order uniqueness
	processRules, err := s.r.GetProcessRulesByCategoryId(ctx, processRule.CategoryID.String())
	if err == nil {
		for _, processRule := range processRules {
			if processRule.Order == dto.Order {
				return uuid.Nil, errors.New("order number already exists in this category")
			}
		}
	}

	// compare os processRules e crie o novo com 1 numero a mais que o maior
	var maxOrder int8
	for _, processRule := range processRules {
		if processRule.Order > maxOrder {
			maxOrder = processRule.Order
		}
	}
	processRule.Order = maxOrder + 1

	processRuleModel := &model.ProcessRule{}
	processRuleModel.FromDomain(processRule)
	err = s.r.CreateProcessRule(ctx, processRuleModel)

	if err != nil {
		return uuid.Nil, err
	}

	category, err := s.rc.GetCategoryById(ctx, processRule.CategoryID.String())
	if err != nil {
		return uuid.Nil, err
	}

	category.UseProcessRule = true
	if err = s.rc.UpdateCategory(ctx, category); err != nil {
		return uuid.Nil, err
	}

	return processRule.ID, nil
}

func (s *Service) UpdateProcessRule(ctx context.Context, dtoId *entitydto.IDRequest, dto *productcategorydto.ProcessRuleUpdateDTO) error {
	processRuleModel, err := s.r.GetProcessRuleById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	processRule := processRuleModel.ToDomain()
	if err = dto.UpdateDomain(processRule); err != nil {
		return err
	}

	// Validate Order uniqueness
	if existingRule, err := s.r.GetProcessRuleByCategoryIdAndOrder(ctx, processRule.CategoryID.String(), int8(processRule.Order)); err == nil {
		if existingRule.ID != processRule.ID {
			return errors.New("order number already exists in this category")
		}
	}

	processRuleModel.FromDomain(processRule)
	if err = s.r.UpdateProcessRule(ctx, processRuleModel); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteProcessRule(ctx context.Context, dto *entitydto.IDRequest) error {
	processRuleModel, err := s.r.GetProcessRuleById(ctx, dto.ID.String())
	if err != nil {
		return err
	}

	if err := s.r.DeleteProcessRule(ctx, dto.ID.String()); err != nil {
		return err
	}

	category, err := s.rc.GetCategoryById(ctx, processRuleModel.CategoryID.String())
	if err != nil {
		return err
	}

	if len(category.ProcessRules) == 0 {
		category.UseProcessRule = false
		if err = s.rc.UpdateCategory(ctx, category); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetProcessRuleById(ctx context.Context, dto *entitydto.IDRequest) (*productcategorydto.ProcessRuleDTO, error) {
	if processRuleModel, err := s.r.GetProcessRuleById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		processRule := processRuleModel.ToDomain()

		processRuleDto := &productcategorydto.ProcessRuleDTO{}
		processRuleDto.FromDomain(processRule)
		return processRuleDto, nil
	}
}

func (s *Service) GetProcessRulesByCategoryId(ctx context.Context, dto *entitydto.IDRequest) ([]productcategorydto.ProcessRuleDTO, error) {
	if processRuleModels, err := s.r.GetProcessRulesByCategoryId(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		processRules := []productcategorydto.ProcessRuleDTO{}
		for _, processRuleModel := range processRuleModels {
			processRule := processRuleModel.ToDomain()

			processRuleDto := &productcategorydto.ProcessRuleDTO{}
			processRuleDto.FromDomain(processRule)
			processRules = append(processRules, *processRuleDto)
		}
		return processRules, nil
	}
}

func (s *Service) GetProcessRulesWithOrderProcessByCategoryId(ctx context.Context, dto *entitydto.IDRequest) ([]productcategorydto.ProcessRuleWithOrderProcessDTO, error) {
	if processRuleModels, err := s.r.GetProcessRulesWithOrderProcessByCategoryId(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		processRules := []productcategorydto.ProcessRuleWithOrderProcessDTO{}
		for _, processRuleModel := range processRuleModels {
			processRule := processRuleModel.ToDomain()

			processRuleDto := &productcategorydto.ProcessRuleWithOrderProcessDTO{}
			processRuleDto.FromDomain(processRule)
			processRules = append(processRules, *processRuleDto)
		}
		return processRules, nil
	}
}

func (s *Service) GetAllProcessRules(ctx context.Context, page, perPage int, isActive bool) ([]productcategorydto.ProcessRuleDTO, int, error) {
	processRuleModels, total, err := s.r.GetAllProcessRules(ctx, page, perPage, isActive)
	if err != nil {
		return nil, 0, err
	}
	return s.processRulesToDto(processRuleModels), total, nil
}
func (s *Service) GetAllProcessRulesWithOrderProcess(ctx context.Context) ([]productcategorydto.ProcessRuleDTO, error) {
	if processRuleModels, err := s.r.GetAllProcessRulesWithOrderProcess(ctx); err != nil {
		return nil, err
	} else {
		return s.processRulesToDto(processRuleModels), nil
	}
}

func (s *Service) processRulesToDto(processRuleModels []model.ProcessRule) []productcategorydto.ProcessRuleDTO {
	var processRuleDTOs []productcategorydto.ProcessRuleDTO

	for _, processRuleModel := range processRuleModels {
		processRule := processRuleModel.ToDomain()
		processRuleDTO := productcategorydto.ProcessRuleDTO{}
		processRuleDTO.FromDomain(processRule)
		processRuleDTOs = append(processRuleDTOs, processRuleDTO)
	}

	return processRuleDTOs
}

func (s *Service) ReorderProcessRules(ctx context.Context, dto *productcategorydto.ProcessRuleReorderDTO) error {
	processRules := make([]model.ProcessRule, len(dto.ProcessRules))

	for i, item := range dto.ProcessRules {
		processRules[i] = model.ProcessRule{
			Entity: entitymodel.Entity{ID: item.ID},
			ProcessRuleCommonAttributes: model.ProcessRuleCommonAttributes{
				Order: item.Order,
			},
		}
	}

	return s.r.UpdateProcessRulesOrder(ctx, processRules)
}
