package printmanagerusecases

import (
	"context"
	"fmt"

	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/pos"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/rabbitmq"
)

func (s *Service) RequestPrintShift(ctx context.Context, req *entitydto.IDRequest) error {
	shift, err := s.shiftService.GetOnlyShiftDomainByID(ctx, req)
	if err != nil {
		return err
	}

	company, err := s.getCompany(ctx)
	if err != nil {
		return err
	}

	printerName, _ := company.Preferences.GetString(companyentity.PrinterShiftReport)
	if err := s.rabbitmq.SendPrintMessage(rabbitmq.SHIFT_EX, company.SchemaName, shift.ID.String(), printerName); err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to send print message")
	}

	return nil
}

// PrintShift retrieves the shift by ID and returns its printable representation.
func (s *Service) PrintShift(ctx context.Context, req *entitydto.IDRequest) ([]byte, error) {
	shift, err := s.shiftService.GetOnlyShiftDomainByID(ctx, req)
	if err != nil {
		return nil, err
	}

	data, err := pos.FormatShift(shift)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// PrintShiftHTML retrieves the shift by ID and returns its HTML representation.
func (s *Service) PrintShiftHTML(ctx context.Context, req *entitydto.IDRequest) ([]byte, error) {
	shift, err := s.shiftService.GetOnlyShiftDomainByID(ctx, req)
	if err != nil {
		return nil, err
	}

	data, err := pos.RenderShiftHTML(shift)
	if err != nil {
		return nil, err
	}

	return data, nil
}
