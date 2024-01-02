package tableorderusecases

import (
	"context"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

func (s *Service) DeleteTableOrder(ctx context.Context, dtoID *entitydto.IdRequest) error {
	return s.rto.DeleteTableOrder(ctx, dtoID.ID.String())
}
