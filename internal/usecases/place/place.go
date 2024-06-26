package placeusecases

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	tableentity "github.com/willjrcom/sales-backend-go/internal/domain/table"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	placedto "github.com/willjrcom/sales-backend-go/internal/infra/dto/place"
)

var (
	ErrPlacePositionIsUsed = func(name string) error { return fmt.Errorf("place position already used by table: %s", name) }
	ErrTableAlreadyInPlace = func(name string, column, row int) error {
		return fmt.Errorf("table already in place: %s on position: column %d, row %d", name, column, row)
	}
	ErrToSearchUsedTable = errors.New("error to search used table")
)

type Service struct {
	r tableentity.PlaceRepository
}

func NewService(c tableentity.PlaceRepository) *Service {
	return &Service{r: c}
}

func (s *Service) CreatePlace(ctx context.Context, dto *placedto.CreatePlaceInput) (uuid.UUID, error) {
	place, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	err = s.r.CreatePlace(ctx, place)

	if err != nil {
		return uuid.Nil, err
	}

	return place.ID, nil
}

func (s *Service) UpdatePlace(ctx context.Context, dtoId *entitydto.IdRequest, dto *placedto.UpdatePlaceInput) error {
	place, err := s.r.GetPlaceById(ctx, dtoId.ID.String())
	if err != nil {
		return err
	}

	if err := dto.UpdateModel(place); err != nil {
		return err
	}

	if err = s.r.UpdatePlace(ctx, place); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeletePlace(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.r.GetPlaceById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.r.DeletePlace(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetPlaceById(ctx context.Context, dto *entitydto.IdRequest) (*tableentity.Place, error) {
	if place, err := s.r.GetPlaceById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return place, nil
	}
}

func (s *Service) GetAllPlaces(ctx context.Context) ([]tableentity.Place, error) {
	return s.r.GetAllPlaces(ctx)
}

func (s *Service) AddTableToPlace(ctx context.Context, dto *placedto.AddTableToPlaceInput) error {
	placeToTable, err := dto.ToModel()
	if err != nil {
		return err
	}

	// If place position already used
	if usedPlacePosition, _ := s.r.GetTableToPlaceByPlaceIDAndPosition(ctx, placeToTable.PlaceID, placeToTable.Column, placeToTable.Row); usedPlacePosition != nil {
		if usedPlacePosition.TableID == placeToTable.TableID {
			return nil
		}

		return ErrPlacePositionIsUsed(usedPlacePosition.Table.Name)
	}

	// If table ID already used
	if err := s.r.AddTableToPlace(ctx, placeToTable); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			tableFromPlace, err := s.r.GetTableToPlaceByTableID(ctx, placeToTable.TableID)
			if err != nil {
				return ErrToSearchUsedTable
			}

			place, err := s.r.GetPlaceById(ctx, tableFromPlace.PlaceID.String())
			if err != nil {
				return ErrToSearchUsedTable
			}

			return ErrTableAlreadyInPlace(place.Name, tableFromPlace.Column, tableFromPlace.Row)
		}

		return err
	}

	return nil
}

func (s *Service) RemoveTableFromPlace(ctx context.Context, dto *entitydto.IdRequest) error {
	if err := s.r.RemoveTableFromPlace(ctx, dto.ID); err != nil {
		return err
	}

	return nil
}
