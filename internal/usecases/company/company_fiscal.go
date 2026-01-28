package companyusecases

import (
	"context"
	"fmt"

	"github.com/willjrcom/sales-backend-go/internal/infra/service/focusnfe"
)

// RegisterFiscalCompany registers the company in Focus NFe API
func (s *Service) RegisterFiscalCompany(ctx context.Context) error {
	companyModel, err := s.r.GetCompany(ctx)
	if err != nil {
		return err
	}

	company := companyModel.ToDomain()

	if !company.FiscalEnabled {
		return fmt.Errorf("fiscal feature is not enabled for this company")
	}

	if s.focusClient == nil || !s.focusClient.Enabled() {
		return fmt.Errorf("focus nfe client is not enabled")
	}

	// Prepare request
	// We need to map RegimeTributario to Focus NFe format
	// Focus NFe: 1=Simples Nacional, 3=Regime Normal
	regime := "1" // Default Simples
	if company.RegimeTributario == 3 {
		regime = "3"
	}

	req := &focusnfe.CompanyRegistryRequest{
		Nome:              company.BusinessName,
		NomeFantasia:      company.TradeName,
		InscricaoEstadual: company.InscricaoEstadual,
		CNPJ:              company.Cnpj,
		RegimeTributario:  regime,
		Email:             company.Email,
		// Telefone: company.Contacts[0] ??
	}

	if len(company.Contacts) > 0 {
		req.Telefone = company.Contacts[0]
	}

	if company.Address != nil {
		req.Logradouro = company.Address.Street
		req.Numero = company.Address.Number
		req.Bairro = company.Address.Neighborhood
		req.CEP = company.Address.Cep
		req.Municipio = company.Address.City
		req.UF = company.Address.UF
	}

	resp, err := s.focusClient.CadastrarEmpresa(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to register company in Focus NFe: %w", err)
	}

	// You could save the ID returned by Focus NFe if needed,
	// but usually CNPJ is the key.
	_ = resp

	return nil
}
