package orderusecases

import (
	"context"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

func (s *Service) DeleteOrderByID(ctx context.Context, dtoId *entitydto.IDRequest) error {
	return s.ro.DeleteOrder(ctx, dtoId.ID.String())
}
