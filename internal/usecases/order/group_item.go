package orderusecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	groupitemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/group_item"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	ErrItemsFinished              = errors.New("items already finished")
	ErrSizeMustBeTheSame          = errors.New("size must be the same")
	ErrComplementItemAlreadyAdded = errors.New("complement item already added")
	ErrComplementItemNotFound     = errors.New("complement item not found")
)

type GroupItemService struct {
	db  *bun.DB
	r   model.GroupItemRepository
	ri  model.ItemRepository
	rp  model.ProductRepository
	re  model.EmployeeRepository
	sop *OrderProcessService
	so  *OrderService
	si  *ItemService
}

func NewGroupItemService(db *bun.DB, rgi model.GroupItemRepository) *GroupItemService {
	return &GroupItemService{db: db, r: rgi}
}

func (s *GroupItemService) AddDependencies(ri model.ItemRepository, rp model.ProductRepository, so *OrderService, sop *OrderProcessService, re model.EmployeeRepository, si *ItemService) {
	s.ri = ri
	s.rp = rp
	s.so = so
	s.sop = sop
	s.re = re
	s.si = si
}

func (s *GroupItemService) GetGroupByID(ctx context.Context, dto *entitydto.IDRequest) (groupItemDTO *groupitemdto.GroupItemDTO, err error) {
	groupItemModel, err := s.r.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return nil, err
	}

	groupItem := groupItemModel.ToDomain()
	groupItemDTO = &groupitemdto.GroupItemDTO{}
	groupItemDTO.FromDomain(groupItem)
	return
}

func (s *GroupItemService) DeleteGroupItem(ctx context.Context, dto *entitydto.IDRequest) (err error) {
	groupItem, err := s.r.GetGroupByID(ctx, dto.ID.String(), true)
	if err != nil {
		return err
	}

	var complementItemID *string
	if groupItem.ComplementItemID != nil {
		complementItemID = new(string)
		*complementItemID = groupItem.ComplementItemID.String()
	}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	if err := s.r.DeleteGroupItemWithTx(ctx, tx, groupItem.ID.String(), complementItemID); err != nil {
		return err
	}

	if err := s.so.UpdateOrderTotalWithTx(ctx, tx, groupItem.OrderID.String()); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *GroupItemService) AddComplementItem(ctx context.Context, dto *entitydto.IDRequest, dtoComplement *entitydto.IDRequest, dtoVariationId *entitydto.IDRequest) (err error) {
	groupItemModel, err := s.r.GetGroupByID(ctx, dto.ID.String(), true)

	if err != nil {
		return err
	}

	groupItem := groupItemModel.ToDomain()

	if groupItem.ComplementItemID != nil {
		return ErrComplementItemAlreadyAdded
	}

	productComplementModel, err := s.rp.GetProductById(ctx, dtoComplement.ID.String())

	if err != nil {
		return err
	}

	productComplement := productComplementModel.ToDomain()

	var variation *productentity.ProductVariation
	if dtoVariationId != nil {
		for _, v := range productComplement.Variations {
			if v.ID == dtoVariationId.ID {
				variation = &v
				break
			}
		}
	} else {
		// Fallback for backwards compatibility if no variation passed
		for _, v := range productComplement.Variations {
			if v.Size != nil && v.Size.Name == groupItem.Size {
				variation = &v
				break
			}
		}
	}

	if variation == nil {
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

	size := ""
	if variation.Size != nil {
		size = variation.Size.Name
	}

	itemComplement := orderentity.NewItem(productComplement.Name, variation.Price, groupItem.Quantity, size, productComplement.ID, variation.ID, productComplement.CategoryID, nil)
	itemComplement.AddSizeToName()

	userID, ok := ctx.Value(model.UserValue("user_id")).(string)
	attendantID := uuid.Nil
	if ok {
		userIDUUID := uuid.MustParse(userID)
		employee, err := s.re.GetEmployeeByUserID(ctx, userIDUUID.String())
		if err != nil {
			return err
		}
		attendantID = employee.ID
	}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	// Reserve stock for the complement item
	if err := s.si.reserveStockFromItemWithTx(ctx, tx, itemComplement, groupItem.OrderID, attendantID); err != nil {
		return err
	}

	itemComplementModel := &model.Item{}
	itemComplementModel.FromDomain(itemComplement)
	if err = s.ri.AddItemWithTx(ctx, tx, itemComplementModel); err != nil {
		return err
	}

	groupItem.ComplementItemID = &itemComplement.ID
	groupItem.ComplementItem = itemComplement

	groupItem.CalculateTotal()

	groupItemModel.FromDomain(groupItem)
	if err := s.r.UpdateGroupItemWithTx(ctx, tx, groupItemModel); err != nil {
		return err
	}

	if err := s.so.UpdateOrderTotalWithTx(ctx, tx, groupItem.OrderID.String()); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *GroupItemService) DeleteComplementItem(ctx context.Context, dto *entitydto.IDRequest) (err error) {
	groupItemModel, err := s.r.GetGroupByID(ctx, dto.ID.String(), true)
	if err != nil {
		return err
	}

	if groupItemModel.ComplementItemID == nil {
		return ErrComplementItemNotFound
	}

	userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)
	attendantID := uuid.Nil
	if ok {
		userIDUUID := uuid.MustParse(userID)
		employee, err := s.re.GetEmployeeByUserID(ctx, userIDUUID.String())
		if err != nil {
			return err
		}
		attendantID = employee.ID
	}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	// Restore stock for the complement item before deleting it
	if err := s.si.restoreStockFromItemWithTx(ctx, tx, groupItemModel.ComplementItem.ToDomain(), groupItemModel.ToDomain(), attendantID); err != nil {
		fmt.Printf("Aviso: erro ao restaurar estoque para item complementar %s: %v\n", groupItemModel.ComplementItem.Name, err)
	}

	if err := s.ri.DeleteItemWithTx(ctx, tx, groupItemModel.ComplementItemID.String()); err != nil {
		return err
	}

	groupItemModel.ComplementItemID = nil
	groupItemModel.ComplementItem = nil

	groupItem := groupItemModel.ToDomain()
	groupItem.CalculateTotal()

	groupItemModel.FromDomain(groupItem)
	if err := s.r.UpdateGroupItemWithTx(ctx, tx, groupItemModel); err != nil {
		return err
	}

	if err := s.so.UpdateOrderTotalWithTx(ctx, tx, groupItemModel.OrderID.String()); err != nil {
		return err
	}

	return tx.Commit()
}
