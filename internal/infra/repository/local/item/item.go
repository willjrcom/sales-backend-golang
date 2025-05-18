package itemrepositorylocal

import (
	"context"
	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

type ItemRepositoryLocal struct {}

func NewItemRepositoryLocal() model.ItemRepository {
	return &ItemRepositoryLocal{}
}

func (r *ItemRepositoryLocal) AddItem(ctx context.Context, item *model.Item) error {
	return nil
}

func (r *ItemRepositoryLocal) AddAdditionalItem(ctx context.Context, id uuid.UUID, productID uuid.UUID, additionalItem *model.Item) error {
	return nil
}

func (r *ItemRepositoryLocal) DeleteItem(ctx context.Context, id string) error {
	return nil
}

func (r *ItemRepositoryLocal) DeleteAdditionalItem(ctx context.Context, idAdditional uuid.UUID) error {
	return nil
}

func (r *ItemRepositoryLocal) UpdateItem(ctx context.Context, item *model.Item) error {
	return nil
}

func (r *ItemRepositoryLocal) GetItemById(ctx context.Context, id string) (*model.Item, error) {
	return nil, nil
}

func (r *ItemRepositoryLocal) GetItemByAdditionalItemID(ctx context.Context, idAdditional uuid.UUID) (*model.Item, error) {
	return nil, nil
}
