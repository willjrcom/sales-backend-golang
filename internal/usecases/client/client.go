package clientusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	clientdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/client"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	geocodeservice "github.com/willjrcom/sales-backend-go/internal/infra/service/geocode"
)

type Service struct {
	rclient  model.ClientRepository
	rcontact model.ContactRepository
}

func NewService(rcliente model.ClientRepository) *Service {
	return &Service{rclient: rcliente}
}

func (s *Service) AddDependencies(rcontact model.ContactRepository) {
	s.rcontact = rcontact
}

func (s *Service) CreateClient(ctx context.Context, dto *clientdto.ClientCreateDTO) (uuid.UUID, error) {
	client, err := dto.ToDomain()

	if err != nil {
		return uuid.Nil, err
	}

	s.UpdateClientWithCoordinates(ctx, client)

	if err := s.rclient.CreateClient(ctx, client); err != nil {
		return uuid.Nil, err
	}

	return client.ID, nil
}

func (s *Service) UpdateClientWithCoordinates(ctx context.Context, client *cliententity.Client) {
	coordinates, _ := geocodeservice.GetCoordinates(&client.Address.AddressCommonAttributes)
	if coordinates == nil {
		return
	}

	client.Address.AddressCommonAttributes.Coordinates = *coordinates
}

func (s *Service) UpdateClient(ctx context.Context, dtoId *entitydto.IDRequest, dto *clientdto.ClientUpdateDTO) error {
	client, err := s.rclient.GetClientById(ctx, dtoId.ID.String())
	if err != nil {
		return err
	}

	if err := dto.UpdateDomain(client); err != nil {
		return err
	}

	s.UpdateClientWithCoordinates(ctx, client)

	if err := s.rclient.UpdateClient(ctx, client); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteClient(ctx context.Context, dto *entitydto.IDRequest) error {
	if _, err := s.rclient.GetClientById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.rclient.DeleteClient(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetClientById(ctx context.Context, dto *entitydto.IDRequest) (*clientdto.ClientDTO, error) {
	if client, err := s.rclient.GetClientById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		dto := &clientdto.ClientDTO{}
		dto.FromDomain(client)
		return dto, nil
	}
}

func (s *Service) GetClientByContact(ctx context.Context, dto *contactdto.ContactDTO) (*clientdto.ClientDTO, error) {
	contact, err := s.rcontact.GetContactByDddAndNumber(ctx, dto.Ddd, dto.Number, personentity.ContactTypeClient)

	if err != nil {
		return nil, err
	}

	if contact == nil {
		return nil, errors.New(("contact not found"))
	}

	if client, err := s.rclient.GetClientById(ctx, contact.ObjectID.String()); err != nil {
		return nil, err
	} else {
		dto := &clientdto.ClientDTO{}
		dto.FromDomain(client)
		return dto, nil
	}
}

func (s *Service) GetAllClients(ctx context.Context) ([]clientdto.ClientDTO, error) {
	if clients, err := s.rclient.GetAllClients(ctx); err != nil {
		return nil, err
	} else {
		dtos := clientsToDtos(clients)
		return dtos, nil
	}
}

func clientsToDtos(clients []cliententity.Client) []clientdto.ClientDTO {
	dtos := make([]clientdto.ClientDTO, len(clients))
	for i, client := range clients {
		dtos[i].FromDomain(&client)
	}

	return dtos
}
