package clientusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	clientdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/client"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	keysdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/keys"
	geocodeservice "github.com/willjrcom/sales-backend-go/internal/infra/service/geocode"
)

type Service struct {
	rclient  cliententity.Repository
	rcontact personentity.ContactRepository
}

func NewService(rcliente cliententity.Repository) *Service {
	return &Service{rclient: rcliente}
}

func (s *Service) AddDependencies(rcontact personentity.ContactRepository) {
	s.rcontact = rcontact
}

func (s *Service) CreateClient(ctx context.Context, dto *clientdto.ClientCreateDTO) (uuid.UUID, error) {
	client, err := dto.ToModel()

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

func (s *Service) UpdateClient(ctx context.Context, dtoId *entitydto.IdRequest, dto *clientdto.ClientUpdateDTO) error {
	client, err := s.rclient.GetClientById(ctx, dtoId.ID.String())
	if err != nil {
		return err
	}

	if err := dto.UpdateModel(client); err != nil {
		return err
	}

	s.UpdateClientWithCoordinates(ctx, client)

	if err := s.rclient.UpdateClient(ctx, client); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteClient(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.rclient.GetClientById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.rclient.DeleteClient(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetClientById(ctx context.Context, dto *entitydto.IdRequest) (*clientdto.ClientOutput, error) {
	if client, err := s.rclient.GetClientById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		dto := &clientdto.ClientOutput{}
		dto.FromModel(client)
		return dto, nil
	}
}

func (s *Service) GetClientByContact(ctx context.Context, dto *keysdto.KeysInput) (*clientdto.ClientOutput, error) {
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
		dto := &clientdto.ClientOutput{}
		dto.FromModel(client)
		return dto, nil
	}
}

func (s *Service) GetAllClients(ctx context.Context) ([]clientdto.ClientOutput, error) {
	if clients, err := s.rclient.GetAllClients(ctx); err != nil {
		return nil, err
	} else {
		dtos := clientsToDtos(clients)
		return dtos, nil
	}
}

func clientsToDtos(clients []cliententity.Client) []clientdto.ClientOutput {
	dtos := make([]clientdto.ClientOutput, len(clients))
	for i, client := range clients {
		dtos[i].FromModel(&client)
	}

	return dtos
}
