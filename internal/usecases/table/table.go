package tableusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	tabledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/table"
)

var (
	ErrTableIsUsed = errors.New("table is used in products")
)

type Service struct {
	r tableentity.TableRepository
}

func NewService(c tableentity.TableRepository) *Service {
	return &Service{r: c}
}

func (s *Service) CreateTable(ctx context.Context, dto *tabledto.TableCreateDTO) (uuid.UUID, error) {
	table, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if err = s.r.CreateTable(ctx, table); err != nil {
		return uuid.Nil, err
	}

	return table.ID, nil
}

func (s *Service) UpdateTable(ctx context.Context, dtoId *entitydto.IdRequest, dto *tabledto.TableUpdateDTO) error {
	table, err := s.r.GetTableById(ctx, dtoId.ID.String())
	if err != nil {
		return err
	}

	if err := dto.UpdateModel(table); err != nil {
		return err
	}

	if err = s.r.UpdateTable(ctx, table); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteTable(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.r.GetTableById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.r.DeleteTable(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetTableById(ctx context.Context, dto *entitydto.IdRequest) (*tableentity.Table, error) {
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
