package groupitemusecases

import (
	"context"
	"errors"

	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	groupitemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/group_item"
)

var (
	ErrItemsFinished              = errors.New("items already finished")
	ErrSizeMustBeTheSame          = errors.New("size must be the same")
	ErrComplementItemAlreadyAdded = errors.New("complement item already added")
	ErrComplementItemNotFound     = errors.New("complement item not found")
)

type Service struct {
	ri  itementity.ItemRepository
	rgi groupitementity.GroupItemRepository
	rp  productentity.ProductRepository
}

func NewService(ri itementity.ItemRepository, rgi groupitementity.GroupItemRepository, rp productentity.ProductRepository) *Service {
	return &Service{ri: ri, rgi: rgi, rp: rp}
}

func (s *Service) GetGroupByID(ctx context.Context, dto *entitydto.IdRequest) (groupItem *groupitementity.GroupItem, err error) {
	groupItem, err = s.rgi.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return nil, err
	}

	groupItem.CalculateTotalPrice()
	return
}

func (s *Service) GetGroupsByStatus(ctx context.Context, dto *groupitemdto.GroupItemByStatusInput) (groups []groupitementity.GroupItem, err error) {
	return s.rgi.GetGroupsByStatus(ctx, dto.Status)
}

func (s *Service) GetGroupsByOrderIDAndStatus(ctx context.Context, dto *groupitemdto.GroupItemByOrderIDAndStatusInput) (groups []groupitementity.GroupItem, err error) {
	return s.rgi.GetGroupsByOrderIDAndStatus(ctx, dto.OrderID.String(), dto.Status)
}

func (s *Service) DeleteGroupItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
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

func (s *Service) AddComplementItem(ctx context.Context, dto *entitydto.IdRequest, dtoComplement *entitydto.IdRequest) (err error) {
	groupItem, err := s.rgi.GetGroupByID(ctx, dto.ID.String(), false)

	if err != nil {
		return err
	}

	if groupItem.ComplementItemID != nil {
		return ErrComplementItemAlreadyAdded
	}

	product, err := s.rp.GetProductById(ctx, dtoComplement.ID.String())

	if err != nil {
		return err
	}

	if groupItem.Size != product.Size.Name {
		return ErrSizeMustBeTheSame
	}

	item := itementity.NewItem(product.Name, product.Price, groupItem.Quantity, groupItem.Size, itementity.StatusItem(groupItem.Status))

	if err = s.ri.AddItem(ctx, item); err != nil {
		return err
	}

	groupItem.ComplementItemID = &item.ID

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return err
	}

	groupItem, err = s.rgi.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	groupItem.CalculateTotalPrice()

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteComplementItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
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

	groupItem.CalculateTotalPrice()

	if err = s.rgi.UpdateGroupItem(ctx, groupItem); err != nil {
		return err
	}

	return nil
}
