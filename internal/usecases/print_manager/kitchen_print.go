package printmanagerusecases

import (
	"context"
	"fmt"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/pos"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/rabbitmq"
)

func (s *Service) RequestPrintGroupItemKitchen(ctx context.Context, req *entitydto.IDRequest) error {
	groupItem, err := s.groupItemRepository.GetGroupByID(ctx, req.ID.String(), false)
	if err != nil {
		return err
	}

	company, err := s.getCompany(ctx)
	if err != nil {
		return err
	}

	if err := s.rabbitmq.SendPrintMessage(rabbitmq.GROUP_ITEM_EX, company.SchemaName, groupItem.ID.String(), groupItem.PrinterName); err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to send print message")
	}

	return nil
}

// PrintGroupItemKitchen retrieves the order by ID and returns its kitchen-printable bytes
// showing only items and complements, without prices or totals.
func (s *Service) PrintGroupItemKitchen(ctx context.Context, req *entitydto.IDRequest) ([]byte, error) {
	// fetch full order model
	modelGroupItem, err := s.groupItemRepository.GetGroupByID(ctx, req.ID.String(), true)
	if err != nil {
		return nil, err
	}

	company, err := s.getCompany(ctx)
	if err != nil {
		return nil, err
	}

	// convert to domain
	groupItem := modelGroupItem.ToDomain()
	data, err := pos.FormatGroupItemKitchen(groupItem, company)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// PrintGroupItemKitchenHTML retrieves the order by ID and returns its kitchen-printable HTML bytes
func (s *Service) PrintGroupItemKitchenHTML(ctx context.Context, req *entitydto.IDRequest) ([]byte, error) {
	// fetch full order model
	modelGroupItem, err := s.groupItemRepository.GetGroupByID(ctx, req.ID.String(), true)
	if err != nil {
		return nil, err
	}

	company, err := s.getCompany(ctx)
	if err != nil {
		return nil, err
	}

	// convert to domain
	groupItem := modelGroupItem.ToDomain()
	data, err := pos.RenderGroupItemKitchenHTML(groupItem, company)
	if err != nil {
		return nil, err
	}
	return data, nil
}
