package orderusecases

import (
	"context"
	"errors"

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

func (s *ItemService) AddAdditionalItemOrder(ctx context.Context, dto *entitydto.IDRequest, dtoAdditional *itemdto.OrderAdditionalItemCreateDTO) (id uuid.UUID, err error) {
	itemModel, err := s.ri.GetItemById(ctx, dto.ID.String())
	if err != nil {
		return uuid.Nil, errors.New("item not found: " + err.Error())
	}

	groupItemModel, err := s.rgi.GetGroupByID(ctx, itemModel.GroupItemID.String(), true)
	if err != nil {
		return uuid.Nil, errors.New("group item not found: " + err.Error())
	}

	item := itemModel.ToDomain()
	groupItem := groupItemModel.ToDomain()

	if ok, err := groupItem.CanAddItems(); !ok {
		return uuid.Nil, err
	}

	userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)
	attendantID := uuid.Nil
	if ok {
		userIDUUID := uuid.MustParse(userID)
		employee, _ := s.re.GetEmployeeByUserID(ctx, userIDUUID.String())
		if employee != nil {
			attendantID = employee.ID
		}
	}

	// Debit stock for additional item
	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, s.db)
	if err != nil {
		return uuid.Nil, err
	}
	defer cancel()
	defer tx.Rollback()

	if id, err = s.addAdditionalItemWithTx(ctx, tx, item.ID, item.Size, groupItem, dtoAdditional, attendantID); err != nil {
		return uuid.Nil, err
	}

	if err := s.so.UpdateOrderTotalWithTx(ctx, tx, groupItem.OrderID.String()); err != nil {
		return uuid.Nil, err
	}

	if err := tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *ItemService) addAdditionalItemWithTx(ctx context.Context, tx *bun.Tx, itemID uuid.UUID, size string, groupItem *orderentity.GroupItem, dtoAdditional *itemdto.OrderAdditionalItemCreateDTO, attendantID uuid.UUID) (id uuid.UUID, err error) {
	productID, variationID, quantityValue, flavor, err := dtoAdditional.ToDomain()

	if err != nil {
		return uuid.Nil, err
	}

	productAdditionalModel, err := s.rp.GetProductById(ctx, productID.String())
	if err != nil {
		return uuid.Nil, errors.New("additional product not found: " + err.Error())
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
		return uuid.Nil, errors.New("additional variation not found")
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

	additionalItem := orderentity.NewItem(productAdditional.Name, variationAdditional.Price, quantityValue, size, productAdditional.ID, variationAdditional.ID, productAdditional.CategoryID, normalizedFlavor)
	additionalItem.IsAdditional = true
	additionalItem.GroupItemID = groupItem.ID

	additionalItemModel := &model.Item{}
	additionalItemModel.FromDomain(additionalItem)

	// Debit stock for additional item
	if err := s.reserveStockFromItemWithTx(ctx, tx, additionalItem, groupItem.OrderID, attendantID); err != nil {
		return uuid.Nil, err
	}

	if err = s.ri.AddAdditionalItemWithTx(ctx, tx, itemID, productAdditional.ID, additionalItemModel); err != nil {
		return uuid.Nil, errors.New("add additional item error: " + err.Error())
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

	userID, ok := ctx.Value(companyentity.UserValue("user_id")).(string)
	attendantID := uuid.Nil
	if ok {
		userIDUUID := uuid.MustParse(userID)
		employee, _ := s.re.GetEmployeeByUserID(ctx, userIDUUID.String())
		if employee != nil {
			attendantID = employee.ID
		}
	}

	ctx, tx, cancel, err := database.GetTenantTransaction(ctx, s.db)
	if err != nil {
		return err
	}
	defer cancel()
	defer tx.Rollback()

	// Fix: Restaurar estoque do adicional removido antes do banco de dados deletar o relacionamento
	if err := s.restoreStockFromItemWithTx(ctx, tx, additionalItem.ToDomain(), groupItemModel.ToDomain(), attendantID); err != nil {
		return err
	}

	if err = s.ri.DeleteItemWithTx(ctx, tx, dtoAdditional.ID.String()); err != nil {
		return err
	}

	if err := s.so.UpdateOrderTotalWithTx(ctx, tx, groupItemModel.OrderID.String()); err != nil {
		return err
	}

	return tx.Commit()
}
