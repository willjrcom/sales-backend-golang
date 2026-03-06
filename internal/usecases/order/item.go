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
	itemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/item"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
)

var (
	ErrCategoryNotFound             = errors.New("category not found")
	ErrSizeNotFound                 = errors.New("size not found")
	ErrSizeNotActive                = errors.New("size not active")
	ErrGroupItemNotBelongsToOrder   = errors.New("group item not belongs to order")
	ErrGroupNotStaging              = errors.New("group not staging")
	ErrItemNotStagingAndPending     = errors.New("item not staging or pending")
	ErrFractionalQuantityNotAllowed = errors.New("fractional quantity not allowed for this category")
)

type ItemService struct {
	db                *bun.DB
	ri                model.ItemRepository
	rgi               model.GroupItemRepository
	ro                model.OrderRepository
	rp                model.ProductRepository
	rc                model.CategoryRepository
	re                model.EmployeeRepository
	so                *OrderService
	sgi               *GroupItemService
	stockRepo         model.StockRepository
	stockMovementRepo model.StockMovementRepository
}

func NewService(db *bun.DB, ri model.ItemRepository) *ItemService {
	return &ItemService{db: db, ri: ri}
}
func (s *ItemService) AddDependencies(rgi model.GroupItemRepository, ro model.OrderRepository, rp model.ProductRepository, rc model.CategoryRepository, re model.EmployeeRepository, so *OrderService, sgi *GroupItemService, stockRepo model.StockRepository, stockMovementRepo model.StockMovementRepository) {
	s.rgi = rgi
	s.ro = ro
	s.rp = rp
	s.rc = rc
	s.re = re
	s.so = so
	s.sgi = sgi
	s.stockRepo = stockRepo
	s.stockMovementRepo = stockMovementRepo
}

func (s *ItemService) AddItemOrder(ctx context.Context, dto *itemdto.OrderItemCreateDTO) (ids *itemdto.ItemIDDTO, err error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	if exists, _ := s.ro.ExistsOrderById(ctx, dto.OrderID.String()); !exists {
		return nil, errors.New("order not found")
	}

	productModel, err := s.rp.GetProductById(ctx, dto.ProductID.String())
	if err != nil {
		return nil, errors.New("product not found: " + err.Error())
	}

	product := productModel.ToDomain()

	var variation *productentity.ProductVariation
	for _, v := range product.Variations {
		if v.ID == dto.VariationID {
			variation = &v
			break
		}
	}

	if variation == nil {
		return nil, errors.New("variation not found")
	}

	if !variation.IsAvailable {
		return nil, errors.New("product variation not available")
	}

	if !product.IsActive {
		return nil, errors.New("product not active")
	}

	if product.Category == nil {
		return nil, ErrCategoryNotFound
	}

	if product.Category.AllowFractional == false && dto.Quantity != float64(int64(dto.Quantity)) {
		return nil, ErrFractionalQuantityNotAllowed
	}

	if variation.Size == nil {
		return nil, ErrSizeNotFound
	}

	if !variation.Size.IsActive {
		return nil, ErrSizeNotActive
	}

	// If group item id is not provided, create a new group item
	if dto.GroupItemID == nil {
		groupItem, err := s.newGroupItem(ctx, dto.OrderID, product, variation)

		if err != nil {
			return nil, errors.New("new group item error: " + err.Error())
		}

		dto.GroupItemID = &groupItem.ID
	}

	groupItemModel, err := s.rgi.GetGroupByID(ctx, dto.GroupItemID.String(), true)
	if err != nil {
		return nil, errors.New("group item not found: " + err.Error())
	}

	groupItem := groupItemModel.ToDomain()

	// Check if group item belongs to order
	if groupItem.OrderID != dto.OrderID {
		return nil, ErrGroupItemNotBelongsToOrder
	}

	if ok, err := groupItem.CanAddItems(); !ok {
		return nil, err
	}

	item, err := dto.ToDomain(product, variation, groupItem, dto.Quantity)
	if err != nil {
		return nil, err
	}

	// Process Removed Items
	for _, removedDTO := range dto.RemovedItems {
		name, err := removedDTO.ToDomain()
		if err != nil {
			return nil, err
		}
		item.AddRemovedItem(*name)
	}

	itemModel := &model.Item{}
	itemModel.FromDomain(item)

	userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)
	attendantID := uuid.Nil
	if ok {
		userIDUUID := uuid.MustParse(userID)
		employee, err := s.re.GetEmployeeByUserID(ctx, userIDUUID.String())
		if err != nil {
			return nil, err
		}
		attendantID = employee.ID
	}

	// Reservar estoque antes de salvar o item — se estoque insuficiente, loga aviso e continua.
	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, s.db)
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer tx.Rollback()

	// Reserva o estoque do item
	if err := s.reserveStockFromItemWithTx(ctx, tx, item, groupItem.OrderID, attendantID); err != nil {
		return nil, err
	}

	// Salvar item após confirmar estoque.
	if err = s.ri.AddItemWithTx(ctx, tx, itemModel); err != nil {
		return nil, errors.New("add item error: " + err.Error())
	}

	// Update item total if it had additions
	if len(dto.Additions) > 0 {
		for _, additionDTO := range dto.Additions {
			if _, err := s.addAdditionalItemWithTx(ctx, tx, item.ID, item.Size, groupItem, &additionDTO, attendantID); err != nil {
				return nil, err
			}
		}
	}

	// Update complement item
	if groupItem.ComplementItemID != nil {
		complementDelta := *groupItem.ComplementItem
		complementDelta.Quantity = item.Quantity

		// Reserve stock only for the proportional delta of the new item quantity
		if err := s.reserveStockFromItemWithTx(ctx, tx, &complementDelta, groupItem.OrderID, attendantID); err != nil {
			return nil, err
		}

		groupItem.ComplementItem.Quantity += item.Quantity

		complementItemModel := &model.Item{}
		complementItemModel.FromDomain(groupItem.ComplementItem)
		if err = s.ri.UpdateItemWithTx(ctx, tx, complementItemModel); err != nil {
			return nil, errors.New("update complement item error: " + err.Error())
		}
	}

	if err := s.so.UpdateOrderTotalWithTx(ctx, tx, dto.OrderID.String()); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return itemdto.FromDomain(item.ID, groupItem.ID), nil
}

func (s *ItemService) DeleteItemOrder(ctx context.Context, dto *entitydto.IDRequest) (groupItemDeleted bool, err error) {
	itemModel, err := s.ri.GetItemById(ctx, dto.ID.String())
	if err != nil {
		return false, err
	}

	item := itemModel.ToDomain()

	groupItemModel, err := s.rgi.GetGroupByID(ctx, item.GroupItemID.String(), true)
	if err != nil {
		return false, err
	}

	groupItem := groupItemModel.ToDomain()

	userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)
	attendantID := uuid.Nil
	if ok {
		userIDUUID := uuid.MustParse(userID)
		employee, err := s.re.GetEmployeeByUserID(ctx, userIDUUID.String())
		if err != nil {
			return false, err
		}
		attendantID = employee.ID
	}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, s.db)
	if err != nil {
		return false, err
	}
	defer cancel()
	defer tx.Rollback()

	// Restore stock for the item itself
	if err := s.restoreStockFromItemWithTx(ctx, tx, item, groupItem, attendantID); err != nil {
		return false, err
	}

	if err = s.ri.DeleteItemWithTx(ctx, tx, dto.ID.String()); err != nil {
		return false, errors.New("delete item error: " + err.Error())
	}

	// Restaurar estoque dos adicionais do item
	for _, addItem := range item.AdditionalItems {
		if err := s.restoreStockFromItemWithTx(ctx, tx, &addItem, groupItem, attendantID); err != nil {
			return false, err
		}
	}

	// Re-load group item to get updated items list (it will be empty if it was the last item)
	groupItemModel, err = s.rgi.GetGroupByIDWithTx(ctx, tx, itemModel.GroupItemID.String(), true)
	if err != nil {
		return false, errors.New("group item not found: " + err.Error())
	}

	groupItem = groupItemModel.ToDomain()

	// Update complement item quantity
	if groupItem.ComplementItemID != nil && len(groupItem.Items) != 0 {
		// Restore stock for the complement item proportional to the removed item's quantity
		complementDelta := *groupItem.ComplementItem
		complementDelta.Quantity = itemModel.Quantity
		if err := s.restoreStockFromItemWithTx(ctx, tx, &complementDelta, groupItem, attendantID); err != nil {
			return false, err
		}

		groupItem.ComplementItem.Quantity -= itemModel.Quantity
		complementItemModel := &model.Item{}
		complementItemModel.FromDomain(groupItem.ComplementItem)
		if err = s.ri.UpdateItemWithTx(ctx, tx, complementItemModel); err != nil {
			return false, errors.New("update complement item error: " + err.Error())
		}
	}

	// Delete group item if there are no items
	if len(groupItem.Items) == 0 {
		var complementItemID *string
		if groupItem.ComplementItemID != nil {
			// If deleting the last item and a complement exists, restore ALL remaining stock of the complement
			if err := s.restoreStockFromItemWithTx(ctx, tx, groupItem.ComplementItem, groupItem, attendantID); err != nil {
				fmt.Printf("Aviso: erro ao restaurar estoque total do item complementar %s na deleção do grupo: %v\n", groupItem.ComplementItem.Name, err)
			}

			complementItemID = new(string)
			*complementItemID = groupItem.ComplementItemID.String()
		}

		if err = s.rgi.DeleteGroupItemWithTx(ctx, tx, groupItem.ID.String(), complementItemID); err != nil {
			return false, err
		}

		if err := s.so.UpdateOrderTotalWithTx(ctx, tx, groupItem.OrderID.String()); err != nil {
			return false, err
		}

		if err := tx.Commit(); err != nil {
			return false, err
		}

		return true, nil
	}

	if err := s.so.UpdateOrderTotalWithTx(ctx, tx, groupItem.OrderID.String()); err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return false, nil
}

func (s *ItemService) newGroupItem(ctx context.Context, orderID uuid.UUID, product *productentity.Product, variation *productentity.ProductVariation) (groupItem *orderentity.GroupItem, err error) {
	groupCommonAttributes := orderentity.GroupCommonAttributes{
		OrderID: orderID,
		GroupDetails: orderentity.GroupDetails{
			CategoryID:     product.CategoryID,
			Size:           variation.Size.Name,
			NeedPrint:      product.Category.NeedPrint,
			PrinterName:    product.Category.PrinterName,
			UseProcessRule: product.Category.UseProcessRule,
		},
	}

	groupItem = orderentity.NewGroupItem(groupCommonAttributes)

	groupItemModel := &model.GroupItem{}
	groupItemModel.FromDomain(groupItem)
	err = s.rgi.CreateGroupItem(ctx, groupItemModel)
	return
}

func (s *ItemService) getStockToMove(ctx context.Context, item *orderentity.Item) (*model.Stock, error) {
	// 1. Tentar buscar estoque específico para a variação
	stockModel, err := s.stockRepo.GetStockByVariationID(ctx, item.ProductVariationID.String())
	if err == nil {
		return stockModel, nil
	}

	// 2. Se não encontrar, buscar todos os estoques do produto
	stocks, err := s.stockRepo.GetStockByProductID(ctx, item.ProductID.String())
	if err != nil || len(stocks) == 0 {
		fmt.Printf("produto %s não tem controle de estoque configurado", item.Name)
		return nil, nil
	}

	// 3. Procurar por um registro de estoque "global" (sem variation_id)
	for _, st := range stocks {
		if st.ProductVariationID == nil {
			return &st, nil
		}
	}

	// Se chegou aqui, existem registros de estoque para o produto, mas nenhum para a variação solicitada nem global
	return nil, fmt.Errorf("produto %s não tem estoque para a variação solicitada e não possui estoque global", item.Name)
}
