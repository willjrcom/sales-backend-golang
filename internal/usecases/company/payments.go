package companyusecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
)

// ListCompanyPayments returns subscription payments for the authenticated company.
func (s *Service) ListCompanyPayments(ctx context.Context, page, perPage, month, year int) ([]companydto.CompanyPaymentDTO, int, error) {
	companyModel, err := s.r.GetCompany(ctx)
	if err != nil {
		return nil, 0, err
	}

	payments, total, err := s.companyPaymentRepo.ListCompanyPayments(ctx, companyModel.ID, page, perPage, month, year)
	if err != nil {
		return nil, 0, err
	}

	dtos := make([]companydto.CompanyPaymentDTO, len(payments))
	for i := range payments {
		dto := companydto.CompanyPaymentDTO{}
		dto.FromDomain(payments[i].ToDomain())
		dtos[i] = dto
	}

	return dtos, total, nil
}

func (s *Service) CancelPayment(ctx context.Context, paymentID uuid.UUID) error {
	payment, err := s.companyPaymentRepo.GetCompanyPaymentByID(ctx, paymentID)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	if payment.Status != string(companyentity.PaymentStatusPending) {
		return fmt.Errorf("payment cannot be cancelled (status: %s)", payment.Status)
	}

	if payment.IsMandatory {
		return fmt.Errorf("payment cannot be cancelled (mandatory)")
	}

	// Unlink costs if any
	if err := s.costRepo.UnlinkCostsFromPayment(ctx, paymentID); err != nil {
		return fmt.Errorf("failed to unlink costs: %w", err)
	}

	payment.Status = string(companyentity.PaymentStatusCancelled)
	// Ensure UpdateCompanyPayment is available and works as expected
	if err := s.companyPaymentRepo.UpdateCompanyPayment(ctx, payment); err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}
