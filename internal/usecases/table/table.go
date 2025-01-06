package tableusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
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

	if err = s.r.CreateTable(ctx, table); err != nil {
		return uuid.Nil, err
	}

	return table.ID, nil
}

func (s *Service) UpdateTable(ctx context.Context, dtoId *entitydto.IDRequest, dto *tabledto.TableUpdateDTO) error {
	table, err := s.r.GetTableById(ctx, dtoId.ID.String())
	if err != nil {
		return err
	}

	if err := dto.UpdateDomain(table); err != nil {
		return err
	}

	if err = s.r.UpdateTable(ctx, table); err != nil {
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

func (s *Service) GetTableById(ctx context.Context, dto *entitydto.IDRequest) (*tableentity.Table, error) {
	if table, err := s.r.GetTableById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return table, nil
	}
}

func (s *Service) GetAllTables(ctx context.Context) ([]tableentity.Table, error) {
	return s.r.GetAllTables(ctx)
}

func (s *Service) GetUnusedTables(ctx context.Context) ([]tableentity.Table, error) {
	return s.r.GetUnusedTables(ctx)
}
