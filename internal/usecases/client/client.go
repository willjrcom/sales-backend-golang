package clientusecases

import (
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	clientdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/client"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	filterdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/filter"
)

type Service struct {
	Repository cliententity.Repository
}

func NewService(repository cliententity.Repository) *Service {
	return &Service{Repository: repository}
}

func (s *Service) RegisterClient(dto *clientdto.RegisterClientInput) error {
	person, err := dto.ToModel()

	if err != nil {
		return err
	}

	if err := s.Repository.RegisterClient(person); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateClient(dtoId *entitydto.IdRequest, dto *clientdto.UpdateClientInput) error {
	client, err := s.Repository.GetClientById(dtoId.Id.String())

	if err != nil {
		return err
	}

	if err := dto.UpdateModel(client); err != nil {
		return err
	}

	if err := s.Repository.UpdateClient(client); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteClient(dto *entitydto.IdRequest) error {
	if _, err := s.Repository.GetClientById(dto.Id.String()); err != nil {
		return err
	}

	if err := s.Repository.DeleteClient(dto.Id.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetClientById(dto *entitydto.IdRequest) (*cliententity.Client, error) {
	if client, err := s.Repository.GetClientById(dto.Id.String()); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func (s *Service) GetAllClient(dto *filterdto.Filter) ([]cliententity.Client, error) {
	if clients, err := s.Repository.GetAllClient(dto.Key, dto.Value); err != nil {
		return nil, err
	} else {
		return clients, nil
	}
}
