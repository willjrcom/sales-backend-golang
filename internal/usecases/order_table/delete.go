package ordertableusecases

import (
	"context"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

func (s *Service) DeleteOrderTable(ctx context.Context, dtoID *entitydto.IdRequest) error {
	return s.rto.DeleteOrderTable(ctx, dtoID.ID.String())
}
