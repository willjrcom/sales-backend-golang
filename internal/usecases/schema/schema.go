package schemausecases

import (
	"context"

	schemaentity "github.com/willjrcom/sales-backend-go/internal/domain/schema"
	schemadto "github.com/willjrcom/sales-backend-go/internal/infra/dto/schema"
)

type Service struct {
	r schemaentity.Repository
}

func NewService(r schemaentity.Repository) *Service {
	return &Service{r: r}
}

func (s *Service) NewSchema(ctx context.Context, dto *schemadto.SchemaInput) {
	s.r.NewSchema(ctx, dto.Name)
}
