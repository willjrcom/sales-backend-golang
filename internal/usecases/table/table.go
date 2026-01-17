package tableusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	tabledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/table"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	ErrTableIsUsed = errors.New("table is used in products")
)

type Service struct {
	r model.TableRepository
}

func NewService(c model.TableRepository) *Service {
	return &Service{r: c}
}

func (s *Service) CreateTable(ctx context.Context, dto *tabledto.TableCreateDTO) (uuid.UUID, error) {
	table, err := dto.ToDomain()

	if err != nil {
		return uuid.Nil, err
	}

	tableModel := &model.Table{}
	tableModel.FromDomain(table)
	if err = s.r.CreateTable(ctx, tableModel); err != nil {
		return uuid.Nil, err
	}

	return table.ID, nil
}

func (s *Service) UpdateTable(ctx context.Context, dtoId *entitydto.IDRequest, dto *tabledto.TableUpdateDTO) error {
	tableModel, err := s.r.GetTableById(ctx, dtoId.ID.String())
	if err != nil {
		return err
	}

	table := tableModel.ToDomain()
	if err := dto.UpdateDomain(table); err != nil {
		return err
	}

	tableModel.FromDomain(table)
	if err = s.r.UpdateTable(ctx, tableModel); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteTable(ctx context.Context, dto *entitydto.IDRequest) error {
	if _, err := s.r.GetTableById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.r.DeleteTable(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetTableById(ctx context.Context, dto *entitydto.IDRequest) (*tabledto.TableDTO, error) {
	if tableModel, err := s.r.GetTableById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		table := tableModel.ToDomain()

		tableDTO := &tabledto.TableDTO{}
		tableDTO.FromDomain(table)
		return tableDTO, nil
	}
}

func (s *Service) GetAllTables(ctx context.Context, isActive bool) ([]tabledto.TableDTO, error) {
	tableModels, err := s.r.GetAllTables(ctx, isActive)
	if err != nil {
		return nil, err
	}

	tables := []tabledto.TableDTO{}
	for _, tableModel := range tableModels {
		table := tableModel.ToDomain()

		tableDTO := &tabledto.TableDTO{}
		tableDTO.FromDomain(table)
		tables = append(tables, *tableDTO)
	}

	return tables, nil
}

func (s *Service) GetUnusedTables(ctx context.Context, isActive bool) ([]tabledto.TableDTO, error) {
	tableModels, err := s.r.GetUnusedTables(ctx, isActive)
	if err != nil {
		return nil, err
	}

	tables := []tabledto.TableDTO{}
	for _, tableModel := range tableModels {
		table := tableModel.ToDomain()

		tableDTO := &tabledto.TableDTO{}
		tableDTO.FromDomain(table)
		tables = append(tables, *tableDTO)
	}

	return tables, nil
}
