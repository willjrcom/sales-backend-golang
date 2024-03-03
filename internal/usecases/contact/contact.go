package contactusecases

import (
	"context"

	personentity "github.com/willjrcom/sales-backend-go/internal/domain/person"
	contactdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/contact"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	keysdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/keys"
)

type Service struct {
	r personentity.ContactRepository
}

func NewService(c personentity.ContactRepository) *Service {
	return &Service{r: c}
}

func (s *Service) GetContactById(ctx context.Context, dto *entitydto.IdRequest) (*contactdto.ContactOutput, error) {
	if contact, err := s.r.GetContactById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		output := &contactdto.ContactOutput{}
		output.FromModel(contact)
		return output, nil
	}
}

func (s *Service) FtSearchContacts(ctx context.Context, keys *keysdto.KeysInput) ([]contactdto.ContactOutput, error) {
	if keys.Query == "" {
		return nil, keysdto.ErrInvalidQuery
	}

	if contacts, err := s.r.FtSearchContacts(ctx, keys.Query, personentity.ContactTypeClient); err != nil {
		return nil, err
	} else {
		dtos := contactsToDtos(contacts)
		return dtos, nil
	}
}

func contactsToDtos(contacts []personentity.Contact) []contactdto.ContactOutput {
	dtos := make([]contactdto.ContactOutput, len(contacts))
	for i, contact := range contacts {
		dtos[i].FromModel(&contact)
	}

	return dtos
}
