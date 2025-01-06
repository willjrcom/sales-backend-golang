package contactusecases

import (
	"context"

	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	keysdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/keys"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type Service struct {
	r model.ContactRepository
}

func NewService(c model.ContactRepository) *Service {
	return &Service{r: c}
}

func (s *Service) GetContactById(ctx context.Context, dto *entitydto.IDRequest) (*contactdto.ContactDTO, error) {
	if contactModel, err := s.r.GetContactById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		contact := contactModel.ToDomain()
		output := &contactdto.ContactDTO{}
		output.FromDomain(contact)
		return output, nil
	}
}

func (s *Service) FtSearchContacts(ctx context.Context, keys *keysdto.KeysDTO) ([]contactdto.ContactDTO, error) {
	if keys.Query == "" {
		return nil, keysdto.ErrInvalidQuery
	}

	if contacts, err := s.r.FtSearchContacts(ctx, keys.Query, string(personentity.ContactTypeClient)); err != nil {
		return nil, err
	} else {
		dtos := modelsToDTOs(contacts)
		return dtos, nil
	}
}

func modelsToDTOs(contactModels []model.Contact) []contactdto.ContactDTO {
	dtos := make([]contactdto.ContactDTO, len(contactModels))

	for i, contactModel := range contactModels {
		contact := contactModel.ToDomain()
		dtos[i].FromDomain(contact)
	}

	return dtos
}
