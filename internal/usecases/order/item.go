package orderusecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

func NewService(ri model.ItemRepository) *ItemService {
	return &ItemService{ri: ri}
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

	if _, err := s.ro.GetOrderById(ctx, dto.OrderID.String()); err != nil {
		return nil, err
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

	if err := s.DebitStockFromItem(ctx, item, groupItem, attendantID); err != nil {
		return nil, err
	}

	if err = s.ri.AddItem(ctx, itemModel); err != nil {
		return nil, errors.New("add item error: " + err.Error())
	}

	GroupItemModel := &model.GroupItem{}
	GroupItemModel.FromDomain(groupItem)

	if err = s.rgi.UpdateGroupItem(ctx, groupItemModel); err != nil {
		return nil, errors.New("update group item error: " + err.Error())
	}

	// Update complement item
	if groupItem.ComplementItemID != nil {
		groupItem.ComplementItem.Quantity += item.Quantity

		complementItemModel := &model.Item{}
		complementItemModel.FromDomain(groupItem.ComplementItem)
		if err = s.ri.UpdateItem(ctx, complementItemModel); err != nil {
			return nil, errors.New("update complement item error: " + err.Error())
		}
	}

	if err := s.sgi.UpdateGroupItemTotal(ctx, dto.GroupItemID.String()); err != nil {
		return nil, err
	}

	if err := s.so.UpdateOrderTotal(ctx, dto.OrderID.String()); err != nil {
		return nil, err
	}

	return itemdto.FromDomain(item.ID, groupItem.ID), nil
}

func (s *ItemService) DebitStockFromItem(ctx context.Context, item *orderentity.Item, groupItem *orderentity.GroupItem, attendantID uuid.UUID) error {
	// Buscar estoque do produto/variação
	stockModel, err := s.stockRepo.GetStockByVariationID(ctx, item.ProductVariationID.String())
	if err != nil {
		// Fallback para buscar apenas por ProductID se não houver variação específica (ex: adicionais sem tamanho)
		stocks, err := s.stockRepo.GetStockByProductID(ctx, item.ProductID.String())
		if err != nil || len(stocks) == 0 {
			// Se não há controle de estoque para o produto, continuar
			fmt.Printf("Produto/Variação %s não tem controle de estoque configurado\n", item.Name)
			return err
		}
		stockModel = &stocks[0]
	}

	stock := stockModel.ToDomain()

	// Reservar estoque (permite estoque negativo)
	movement, err := stock.ReserveStock(
		decimal.NewFromFloat(item.Quantity),
		groupItem.OrderID,
		attendantID,
		item.SubTotal,
		item.Total,
	)
	if err != nil {
		fmt.Printf("Erro ao reservar estoque para produto %s: %v\n", item.Name, err)
		return err
	}

	// Salvar movimento
	movementModel := &model.StockMovement{}
	movementModel.FromDomain(movement)
	if err := s.stockMovementRepo.CreateMovement(ctx, movementModel); err != nil {
		fmt.Printf("Erro ao salvar movimento de estoque: %v\n", err)
		return err
	}

	// Atualizar estoque
	stockModel.FromDomain(stock)
	if err := s.stockRepo.UpdateStock(ctx, stockModel); err != nil {
		fmt.Printf("Erro ao atualizar estoque: %v\n", err)
		return err
	}

	fmt.Printf("Estoque debitado para produto %s: %f\n", item.Name, item.Quantity)

	return nil
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

	s.RestoreStockFromItem(ctx, item, groupItemModel.ToDomain(), attendantID)

	if err = s.ri.DeleteItem(ctx, dto.ID.String()); err != nil {
		return false, errors.New("delete item error: " + err.Error())
	}

	groupItem, err := s.rgi.GetGroupByID(ctx, itemModel.GroupItemID.String(), true)

	if err != nil {
		return false, errors.New("group item not found: " + err.Error())
	}

	// Update complement item quantity
	if groupItem.ComplementItemID != nil && len(groupItem.Items) != 0 {
		groupItem.ComplementItem.Quantity -= itemModel.Quantity
		if err = s.ri.UpdateItem(ctx, groupItem.ComplementItem); err != nil {
			return false, errors.New("update complement item error: " + err.Error())
		}
	}

	// Delete group item if there are no items
	if len(groupItem.Items) == 0 {
		var complementItemID *string
		if groupItem.ComplementItemID != nil {
			complementItemID = new(string)
			*complementItemID = groupItem.ComplementItemID.String()
		}

		if err = s.rgi.DeleteGroupItem(ctx, groupItem.ID.String(), complementItemID); err != nil {
			return false, err
		}

		if err := s.so.UpdateOrderTotal(ctx, groupItem.OrderID.String()); err != nil {
			return false, err
		}

		return true, nil
	}

	if err := s.sgi.UpdateGroupItemTotal(ctx, groupItem.ID.String()); err != nil {
		return false, err
	}

	if err := s.so.UpdateOrderTotal(ctx, groupItem.OrderID.String()); err != nil {
		return false, err
	}

	return false, nil
}

func (s ItemService) RestoreStockFromItem(ctx context.Context, item *orderentity.Item, groupItem *orderentity.GroupItem, attendantID uuid.UUID) error {
	// Buscar estoque do produto/variação
	stockModel, err := s.stockRepo.GetStockByVariationID(ctx, item.ProductVariationID.String())
	if err != nil {
		// Fallback para buscar apenas por ProductID
		stocks, err := s.stockRepo.GetStockByProductID(ctx, item.ProductID.String())
		if err != nil || len(stocks) == 0 {
			return errors.New("stock not found")
		}
		stockModel = &stocks[0]
	}

	stock := stockModel.ToDomain()

	// Restaurar estoque
	movement, err := stock.RestoreStock(
		decimal.NewFromFloat(item.Quantity),
		groupItem.OrderID,
		attendantID,
		item.SubTotal,
		item.Total,
	)
	if err != nil {
		fmt.Printf("Erro ao restaurar estoque para produto %s: %v\n", item.Name, err)
		return err
	}

	// Salvar movimento
	movementModel := &model.StockMovement{}
	movementModel.FromDomain(movement)
	if err := s.stockMovementRepo.CreateMovement(ctx, movementModel); err != nil {
		fmt.Printf("Erro ao salvar movimento de estoque: %v\n", err)
		return err
	}

	// Atualizar estoque
	stockModel.FromDomain(stock)
	if err := s.stockRepo.UpdateStock(ctx, stockModel); err != nil {
		fmt.Printf("Erro ao atualizar estoque: %v\n", err)
		return err
	}

	fmt.Printf("Estoque creditado para produto %s: %f\n", item.Name, item.Quantity)

	return nil
}

func (s *ItemService) AddAdditionalItemOrder(ctx context.Context, dto *entitydto.IDRequest, dtoAdditional *itemdto.OrderAdditionalItemCreateDTO) (id uuid.UUID, err error) {
	productID, variationID, quantityValue, flavor, err := dtoAdditional.ToDomain()

	if err != nil {
		return uuid.Nil, err
	}

	item, err := s.ri.GetItemById(ctx, dto.ID.String())

	if err != nil {
		return uuid.Nil, errors.New("item not found: " + err.Error())
	}

	groupItemModel, err := s.rgi.GetGroupByID(ctx, item.GroupItemID.String(), true)
	if err != nil {
		return uuid.Nil, errors.New("group item not found: " + err.Error())
	}

	groupItem := groupItemModel.ToDomain()

	if ok, err := groupItem.CanAddItems(); !ok {
		return uuid.Nil, err
	}

	productAdditionalModel, err := s.rp.GetProductById(ctx, productID.String())
	if err != nil {
		return uuid.Nil, errors.New("product not found: " + err.Error())
	}

	productAdditional := productAdditionalModel.ToDomain()

	var variationAdditional *productentity.ProductVariation
	for _, v := range productAdditional.Variations {
		if v.ID == variationID {
			variationAdditional = &v
			break
		}
	}

	if variationAdditional == nil {
		return uuid.Nil, errors.New("variation not found")
	}

	found := false
	for _, additionalCategory := range groupItem.Category.AdditionalCategories {
		if additionalCategory.ID == productAdditional.CategoryID {
			found = true
			break
		}
	}

	if !found {
		return uuid.Nil, errors.New("additional category does not belong to this category")
	}

	normalizedFlavor, err := itemdto.NormalizeFlavor(flavor, productAdditional.Flavors)
	if err != nil {
		return uuid.Nil, err
	}

	additionalItem := orderentity.NewItem(productAdditional.Name, variationAdditional.Price, quantityValue, item.Size, productAdditional.ID, variationAdditional.ID, productAdditional.CategoryID, normalizedFlavor)
	additionalItem.IsAdditional = true
	additionalItem.GroupItemID = groupItem.ID

	additionalItemModel := &model.Item{}
	additionalItemModel.FromDomain(additionalItem)

	if err = s.ri.AddAdditionalItem(ctx, item.ID, productAdditional.ID, additionalItemModel); err != nil {
		return uuid.Nil, errors.New("add additional item error: " + err.Error())
	}

	if err := s.UpdateItemTotal(ctx, item.ID.String()); err != nil {
		return uuid.Nil, err
	}

	if err := s.sgi.UpdateGroupItemTotal(ctx, groupItem.ID.String()); err != nil {
		return uuid.Nil, err
	}

	if err := s.so.UpdateOrderTotal(ctx, groupItem.OrderID.String()); err != nil {
		return uuid.Nil, err
	}

	return additionalItem.ID, nil
}

func (s *ItemService) DeleteAdditionalItemOrder(ctx context.Context, dtoAdditional *entitydto.IDRequest) (err error) {
	additionalItem, err := s.ri.GetItemById(ctx, dtoAdditional.ID.String())
	if err != nil {
		return errors.New("item not found: " + err.Error())
	}

	groupItemModel, err := s.rgi.GetGroupByID(ctx, additionalItem.GroupItemID.String(), true)
	if err != nil {
		return errors.New("group item not found: " + err.Error())
	}

	item, err := s.ri.GetItemByAdditionalItemID(ctx, additionalItem.ID)
	if err != nil {
		return err
	}

	if err = s.ri.DeleteAdditionalItem(ctx, dtoAdditional.ID); err != nil {
		return err
	}

	if err := s.UpdateItemTotal(ctx, item.ID.String()); err != nil {
		return err
	}

	if err := s.sgi.UpdateGroupItemTotal(ctx, groupItemModel.ID.String()); err != nil {
		return err
	}

	if err := s.so.UpdateOrderTotal(ctx, groupItemModel.OrderID.String()); err != nil {
		return err
	}

	return nil
}

func (s *ItemService) AddRemovedItem(ctx context.Context, dtoID *entitydto.IDRequest, dto *itemdto.RemovedItemDTO) (err error) {
	name, err := dto.ToDomain()
	if err != nil {
		return err
	}

	itemModel, err := s.ri.GetItemById(ctx, dtoID.ID.String())
	if err != nil {
		return err
	}

	item := itemModel.ToDomain()

	category, err := s.rc.GetCategoryById(ctx, item.CategoryID.String())
	if err != nil {
		return err
	}

	found := false
	for _, removedItem := range category.RemovableIngredients {
		if removedItem == *name {
			found = true
		}
	}

	if !found {
		return errors.New("removed item not found on category")
	}

	// Already added
	for _, removedItem := range item.RemovedItems {
		if removedItem == *name {
			return nil
		}
	}

	item.AddRemovedItem(*name)

	itemModel.FromDomain(item)
	return s.ri.UpdateItem(ctx, itemModel)
}

func (s *ItemService) RemoveRemovedItem(ctx context.Context, dtoID *entitydto.IDRequest, dto *itemdto.RemovedItemDTO) (err error) {
	name, err := dto.ToDomain()
	if err != nil {
		return err
	}

	itemModel, err := s.ri.GetItemById(ctx, dtoID.ID.String())
	if err != nil {
		return err
	}

	item := itemModel.ToDomain()
	item.RemoveRemovedItem(*name)

	itemModel.FromDomain(item)
	return s.ri.UpdateItem(ctx, itemModel)
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

func (s *ItemService) UpdateItemTotal(ctx context.Context, id string) error {
	itemModel, err := s.ri.GetItemById(ctx, id)
	if err != nil {
		return err
	}

	item := itemModel.ToDomain()

	item.CalculateTotal()

	itemModel.FromDomain(item)
	if err := s.ri.UpdateItem(ctx, itemModel); err != nil {
		return err
	}

	return nil
}
