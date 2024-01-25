package schemaservice

import (
	"context"

	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
)

type Service struct {
	r schemaentity.Repository
}

func NewService(r schemaentity.Repository) *Service {
	return &Service{r: r}
}

func (s *Service) NewSchema(ctx context.Context, name string) error {
	return s.r.NewSchema(ctx, name)
}
