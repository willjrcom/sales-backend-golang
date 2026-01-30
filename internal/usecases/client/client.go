package clientusecases

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	cliententity "github.com/willjrcom/sales-backend-go/internal/domain/client"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	clientdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/client"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	geocodeservice "github.com/willjrcom/sales-backend-go/internal/infra/service/geocode"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
)

type Service struct {
	rclient  model.ClientRepository
	rcontact model.ContactRepository
	cs       *companyusecases.Service
}

func NewService(rcliente model.ClientRepository) *Service {
	return &Service{rclient: rcliente}
}

func (s *Service) AddDependencies(rcontact model.ContactRepository, cs *companyusecases.Service) {
	s.rcontact = rcontact
	s.cs = cs
}

func (s *Service) CreateClient(ctx context.Context, dto *clientdto.ClientCreateDTO) (uuid.UUID, error) {
	client, err := dto.ToDomain()

	if err != nil {
		return uuid.Nil, err
	}

	company, err := s.cs.GetCompany(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	if minDeliveryTax, err := company.Preferences.GetDecimal(companyentity.MinDeliveryTax); err == nil {
		if client.Address.DeliveryTax.LessThan(minDeliveryTax) {
			return uuid.Nil, errors.New("delivery tax is less than minimum delivery tax")
		}
	}

	s.UpdateClientWithCoordinates(ctx, client)

	clientModel := &model.Client{}
	clientModel.FromDomain(client)
	if err := s.rclient.CreateClient(ctx, clientModel); err != nil {
		if strings.Contains(err.Error(), `duplicate key value violates unique constraint "idx_unique_contact"`) {
			return uuid.Nil, errors.New("contact already exists")
		}
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
	clientModel, err := s.rclient.GetClientById(ctx, dtoId.ID.String())
	if err != nil {
		return err
	}

	client := clientModel.ToDomain()
	if err := dto.UpdateDomain(client); err != nil {
		return err
	}

	company, err := s.cs.GetCompany(ctx)
	if err != nil {
		return err
	}

	if minDeliveryTax, err := company.Preferences.GetDecimal(companyentity.MinDeliveryTax); err == nil {
		if client.Address.DeliveryTax.LessThan(minDeliveryTax) {
			return errors.New("delivery tax is less than minimum delivery tax")
		}
	}

	s.UpdateClientWithCoordinates(ctx, client)

	clientModel.FromDomain(client)
	if err := s.rclient.UpdateClient(ctx, clientModel); err != nil {
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
	if clientModel, err := s.rclient.GetClientById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		client := clientModel.ToDomain()
		dto := &clientdto.ClientDTO{}
		dto.FromDomain(client)
		return dto, nil
	}
}

func (s *Service) GetClientByContact(ctx context.Context, dto *contactdto.ContactDTO) (*clientdto.ClientDTO, error) {
	contactModel, err := s.rcontact.GetContactByNumber(ctx, dto.Number, string(personentity.ContactTypeClient))

	if err != nil {
		return nil, err
	}

	if contactModel == nil {
		return nil, errors.New(("contact not found"))
	}

	if clientModel, err := s.rclient.GetClientById(ctx, contactModel.ObjectID.String()); err != nil {
		return nil, err
	} else {
		client := clientModel.ToDomain()
		dto := &clientdto.ClientDTO{}
		dto.FromDomain(client)
		return dto, nil
	}
}

// GetAllClients retrieves a paginated list of clients and the total count.
func (s *Service) GetAllClients(ctx context.Context, page, perPage int, isActive bool) ([]clientdto.ClientDTO, int, error) {
	clientModels, total, err := s.rclient.GetAllClients(ctx, page, perPage, isActive)
	if err != nil {
		return nil, 0, err
	}
	dtos := modelsToDTOs(clientModels)
	return dtos, total, nil
}

func modelsToDTOs(clientModels []model.Client) []clientdto.ClientDTO {
	dtos := make([]clientdto.ClientDTO, len(clientModels))
	for i, clientModel := range clientModels {
		client := clientModel.ToDomain()
		dtos[i].FromDomain(client)
	}

	return dtos
}
