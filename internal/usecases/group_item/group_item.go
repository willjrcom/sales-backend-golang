package groupitemusecases

import (
	"context"
	"errors"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	ErrItemsFinished              = errors.New("items already finished")
	ErrSizeMustBeTheSame          = errors.New("size must be the same")
	ErrComplementItemAlreadyAdded = errors.New("complement item already added")
	ErrComplementItemNotFound     = errors.New("complement item not found")
)

type Service struct {
	rgi model.GroupItemRepository
	ri  model.ItemRepository
	rp  model.ProductRepository
}

func NewService(rgi model.GroupItemRepository) *Service {
	return &Service{rgi: rgi}
}

func (s *Service) AddDependencies(ri model.ItemRepository, rp model.ProductRepository) {
	s.ri = ri
	s.rp = rp
}

func (s *Service) GetGroupByID(ctx context.Context, dto *entitydto.IDRequest) (groupItem *orderentity.GroupItem, err error) {
	groupItemModel, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return nil, err
	}

	groupItem = groupItemModel.ToDomain()
	return
}

func (s *Service) DeleteGroupItem(ctx context.Context, dto *entitydto.IDRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	var complementItemID *string
	if groupItem.ComplementItemID != nil {
		complementItemID = new(string)
		*complementItemID = groupItem.ComplementItemID.String()
	}

	return s.rgi.DeleteGroupItem(ctx, groupItem.ID.String(), complementItemID)
}

func (s *Service) AddComplementItem(ctx context.Context, dto *entitydto.IDRequest, dtoComplement *entitydto.IDRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	if groupItem.ComplementItemID != nil {
		return ErrComplementItemAlreadyAdded
	}

	productComplement, err := s.rp.GetProductById(ctx, dtoComplement.ID.String())

	if err != nil {
		return err
	}

	if groupItem.Size != productComplement.Size.Name {
		return ErrSizeMustBeTheSame
	}

	// Is valid complement to this category
	found := false
	for _, complementCategory := range groupItem.Category.ComplementCategories {
		if productComplement.Category.ID == complementCategory.ID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("complement category does not belong to this category")
	}

	itemComplement := orderentity.NewItem(productComplement.Name, productComplement.Price, groupItem.Quantity, groupItem.Size, productComplement.ID, productComplement.CategoryID)

	itemComplementModel := &model.Item{}
	itemComplementModel.FromDomain(itemComplement)
	if err = s.ri.AddItem(ctx, itemComplementModel); err != nil {
		return err
	}

	groupItem.ComplementItemID = &itemComplement.ID

	return s.rgi.UpdateGroupItem(ctx, groupItem)
}

func (s *Service) DeleteComplementItem(ctx context.Context, dto *entitydto.IDRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	if groupItem.ComplementItemID == nil {
		return ErrComplementItemNotFound
	}

	if err := s.ri.DeleteItem(ctx, groupItem.ComplementItemID.String()); err != nil {
		return err
	}

	groupItem.ComplementItemID = nil
	groupItem.ComplementItem = nil

	return s.rgi.UpdateGroupItem(ctx, groupItem)
}
