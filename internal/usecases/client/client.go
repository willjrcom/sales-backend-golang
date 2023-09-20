package clientusecases

import (
	"context"

	"github.com/google/uuid"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	clientdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/client"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

type Service struct {
	rclient  cliententity.Repository
	rcontact personentity.ContactRepository
}

func NewService(rcliente cliententity.Repository, rcontact personentity.ContactRepository) *Service {
	return &Service{rclient: rcliente, rcontact: rcontact}
}

func (s *Service) RegisterClient(ctx context.Context, dto *clientdto.RegisterClientInput) (uuid.UUID, error) {
	client, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	if err := s.rclient.RegisterClient(ctx, client); err != nil {
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

func (s *Service) GetClientsBy(ctx context.Context, dto *clientdto.FilterClientInput) ([]clientdto.ClientOutput, error) {
	clients, err := s.rclient.GetAllClients(ctx)

	if err != nil {
		return nil, err
	}

	dtos := clientsToDtos(clients)
	return dtos, nil
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

func (s *Service) RegisterContactToClient(ctx context.Context, dto *contactdto.RegisterContactInput) (uuid.UUID, error) {
	contact, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	// Validate if exists
	if _, err := s.rclient.GetClientById(ctx, contact.PersonID.String()); err != nil {
		return uuid.Nil, err
	}

	if err := s.rcontact.RegisterContact(ctx, contact); err != nil {
		return uuid.Nil, err
	}

	return contact.ID, nil
}
