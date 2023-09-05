package itemusecases

import (
	"github.com/google/uuid"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
)

type Service struct {
	Repository itementity.Repository
}

func NewService(repository itementity.Repository) *Service {
	return &Service{Repository: repository}
}

func (s *Service) AddItemOrder(idOrder, idProduct string, dto *itementity.Item) (uuid.UUID, error) {
	// find order and product by id
	return uuid.New(), nil
}
