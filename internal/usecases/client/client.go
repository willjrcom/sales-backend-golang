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
	keysdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/keys"
)

type Service struct {
	rclient  cliententity.Repository
	rcontact personentity.ContactRepository
}

func NewService(rcliente cliententity.Repository, rcontact personentity.ContactRepository) *Service {
	return &Service{rclient: rcliente, rcontact: rcontact}
}

func (s *Service) CreateClient(ctx context.Context, dto *clientdto.CreateClientInput) (uuid.UUID, error) {
	client, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if err := s.rclient.CreateClient(ctx, client); err != nil {
		return uuid.Nil, err
	}

	return client.ID, nil
}

func (s *Service) UpdateClient(ctx context.Context, dtoId *entitydto.IdRequest, dto *clientdto.UpdateClientInput) error {
	client, err := s.rclient.GetClientById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err := dto.UpdateModel(client); err != nil {
		return err
	}

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

func (s *Service) CreateContactToClient(ctx context.Context, dto *contactdto.CreateContactInput) (uuid.UUID, error) {
	contact, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	// Validate if exists
	if _, err := s.rclient.GetClientById(ctx, contact.ObjectID.String()); err != nil {
		return uuid.Nil, err
	}

	if err := s.rcontact.CreateContact(ctx, contact); err != nil {
		return uuid.Nil, err
	}

	return contact.ID, nil
}

func (s *Service) UpdateContact(ctx context.Context, dtoId *entitydto.IdRequest, dto *contactdto.UpdateContactInput) error {
	contact, err := s.rcontact.GetContactById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err := dto.UpdateModel(contact); err != nil {
		return err
	}

	if err := s.rcontact.UpdateContact(ctx, contact); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteContact(ctx context.Context, dtoId *entitydto.IdRequest) error {
	if _, err := s.rcontact.GetContactById(ctx, dtoId.ID.String()); err != nil {
		return err
	}

	if err := s.rcontact.DeleteContact(ctx, dtoId.ID.String()); err != nil {
		return err
	}

	return nil
}
