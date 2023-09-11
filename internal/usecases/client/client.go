package clientusecases

import (
	addressentity "github.com/willjrcom/sales-backend-go/internal/domain/address"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	addressdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/address"
	clientdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/client"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	filterdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/filter"
)

type Service struct {
	rClient  cliententity.Repository
	rAddress addressentity.Repository
}

func NewService(rc cliententity.Repository, ra addressentity.Repository) *Service {
	return &Service{rClient: rc, rAddress: ra}
}

func (s *Service) RegisterClient(dto *clientdto.RegisterClientInput) error {
	person, err := dto.ToModel()

	if err != nil {
		return err
	}

	if err := s.rClient.RegisterClient(person); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateClient(dtoId *entitydto.IdRequest, dto *clientdto.UpdateClientInput) error {
	client, err := s.rClient.GetClientById(dtoId.ID.String())

	if err != nil {
		return err
	}

	if err := dto.UpdateModel(client); err != nil {
		return err
	}

	if err := s.rClient.UpdateClient(client); err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveClient(dto *entitydto.IdRequest) error {
	if _, err := s.rClient.GetClientById(dto.ID.String()); err != nil {
		return err
	}

	if err := s.rClient.DeleteClient(dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetClientById(dto *entitydto.IdRequest) (*cliententity.Client, error) {
	if client, err := s.rClient.GetClientById(dto.ID.String()); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func (s *Service) GetAllClient(dto *filterdto.Filter) ([]cliententity.Client, error) {
	if clients, err := s.rClient.GetAllClient(dto.Key, dto.Value); err != nil {
		return nil, err
	} else {
		return clients, nil
	}
}

func (s *Service) RegisterAddress(dtoId *entitydto.IdRequest, dto *addressdto.RegisterAddressInput) error {
	address, err := dto.ToModel()

	if err != nil {
		return err
	}

	address.PersonID = dtoId.ID
	return nil
}

func (s *Service) RemoveAddress(dto *entitydto.IdRequest) error {
	return nil
}
