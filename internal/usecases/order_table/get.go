package ordertableusecases

import (
	"context"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	ordertabledto "github.com/willjrcom/sales-backend-go/internal/infra/dto/order_table"
)

func (s *Service) GetTableById(ctx context.Context, dto *entitydto.IDRequest) (*ordertabledto.OrderTableDTO, error) {
	if orderTableModel, err := s.rto.GetOrderTableById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		orderTable := orderTableModel.ToDomain()

		orderTableDTO := &ordertabledto.OrderTableDTO{}
		orderTableDTO.FromDomain(orderTable)
		return orderTableDTO, nil
	}
}

func (s *Service) GetAllTables(ctx context.Context) ([]ordertabledto.OrderTableDTO, error) {
	if orderTableModels, err := s.rto.GetAllOrderTables(ctx); err != nil {
		return nil, err
	} else {
		orderTableDTOs := make([]ordertabledto.OrderTableDTO, 0)
		for _, orderTableModel := range orderTableModels {
			orderTable := orderTableModel.ToDomain()

			orderTableModelDTO := &ordertabledto.OrderTableDTO{}
			orderTableModelDTO.FromDomain(orderTable)
			orderTableDTOs = append(orderTableDTOs, *orderTableModelDTO)
		}
		return orderTableDTOs, nil
	}
}
