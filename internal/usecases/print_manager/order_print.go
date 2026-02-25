package printmanagerusecases

import (
	"context"
	"fmt"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/pos"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/rabbitmq"
)

func (s *Service) RequestPrintOrder(ctx context.Context, req *entitydto.IDRequest) error {
	order, err := s.orderRepository.GetOrderById(ctx, req.ID.String())
	if err != nil {
		return err
	}

	company, err := s.getCompany(ctx)
	if err != nil {
		return err
	}

	printerName, _ := company.Preferences.GetString(companyentity.PrinterOrder)
	path := rabbitmq.ORDER_PATH + order.ID.String()
	if err := s.rabbitmq.SendPrintMessage(rabbitmq.ORDER_EX, company.SchemaName, path, printerName); err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to send print message")
	}

	return nil
}

// PrintOrder retrieves the order by ID and returns its printable representation.
func (s *Service) PrintOrder(ctx context.Context, req *entitydto.IDRequest) ([]byte, error) {
	model, err := s.orderRepository.GetOrderById(ctx, req.ID.String())
	if err != nil {
		return nil, err
	}

	company, err := s.getCompany(ctx)
	if err != nil {
		return nil, err
	}

	order := model.ToDomain()
	data, err := pos.FormatOrder(order, company)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// PrintOrderHTML retrieves the order by ID and returns its HTML representation.
func (s *Service) PrintOrderHTML(ctx context.Context, req *entitydto.IDRequest) ([]byte, error) {
	model, err := s.orderRepository.GetOrderById(ctx, req.ID.String())
	if err != nil {
		return nil, err
	}

	company, err := s.getCompany(ctx)
	if err != nil {
		return nil, err
	}

	order := model.ToDomain()
	data, err := pos.RenderOrderHTML(order, company)
	if err != nil {
		return nil, err
	}

	return data, nil
}
