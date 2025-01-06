package headerservice

import (
	"context"

	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type Service struct {
	r model.SchemaRepository
}

func NewService(r model.SchemaRepository) *Service {
	return &Service{r: r}
}

func (s *Service) NewSchema(ctx context.Context) error {
	return s.r.NewSchema(ctx)
}
