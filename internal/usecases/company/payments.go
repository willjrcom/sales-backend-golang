package companyusecases

import (
	"context"

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
