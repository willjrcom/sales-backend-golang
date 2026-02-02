package fiscalinvoiceusecases

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
	fiscalinvoice "github.com/willjrcom/sales-backend-go/internal/domain/fiscal_invoice"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/focusnfe"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
)

var (
	ErrFiscalNotEnabled                 = errors.New("fiscal invoice functionality is not enabled for this company")
	ErrMissingFiscalData                = errors.New("company is missing required fiscal data (IE, regime tributário)")
	ErrTransmitenotaNotConfigured       = errors.New("focus nfe client is not configured")
	ErrInvoiceAlreadyExists             = errors.New("invoice already exists for this order")
	ErrInvoiceNotFound                  = errors.New("fiscal invoice not found")
	ErrCannotCancelInvoice              = errors.New("invoice cannot be cancelled (not authorized or already cancelled)")
	ErrOrderNotFound                    = errors.New("order not found")
	ErrFunctionalityNotAvailableForPlan = errors.New("funcionalidade não disponível para o plano atual")
)

type Service struct {
	invoiceRepo             model.FiscalInvoiceRepository
	companyRepo             model.CompanyRepository
	companySubscriptionRepo model.CompanySubscriptionRepository
	fiscalSettingsRepo      model.FiscalSettingsRepository
	orderRepo               model.OrderRepository
	usageCostService        *companyusecases.UsageCostService
	focusClient             *focusnfe.Client
}

func NewService(
	invoiceRepo model.FiscalInvoiceRepository,
	companyRepo model.CompanyRepository,
	companySubscriptionRepo model.CompanySubscriptionRepository,
	fiscalSettingsRepo model.FiscalSettingsRepository,
	orderRepo model.OrderRepository,
	usageCostService *companyusecases.UsageCostService,
	focusClient *focusnfe.Client,
) *Service {
	return &Service{
		invoiceRepo:             invoiceRepo,
		companyRepo:             companyRepo,
		companySubscriptionRepo: companySubscriptionRepo,
		fiscalSettingsRepo:      fiscalSettingsRepo,
		orderRepo:               orderRepo,
		usageCostService:        usageCostService,
		focusClient:             focusClient,
	}
}

// EmitirNFCeParaPedido emits NFC-e for an order and registers the cost
func (s *Service) EmitirNFCeParaPedido(ctx context.Context, orderID uuid.UUID) (*fiscalinvoice.FiscalInvoice, error) {
	// Get company from context
	companyModel, err := s.companyRepo.GetCompany(ctx)
	if err != nil {
		return nil, err
	}

	company := companyModel.ToDomain()

	sub, _, err := s.companySubscriptionRepo.GetActiveAndUpcomingSubscriptions(ctx, company.ID)
	if err != nil || sub == nil {
		return nil, ErrFunctionalityNotAvailableForPlan
	}

	if sub.PlanType == companyentity.PlanFree {
		return nil, ErrFunctionalityNotAvailableForPlan
	}

	// Fetch Fiscal Settings
	settings, err := s.fiscalSettingsRepo.GetByCompanyID(ctx, sub.CompanyID)
	if err != nil || settings == nil {
		return nil, ErrFiscalNotEnabled
	}

	// Validate fiscal is enabled
	if !settings.IsActive {
		return nil, ErrFiscalNotEnabled
	}

	// Validate required fiscal data
	if settings.TaxRegime == 0 {
		return nil, ErrMissingFiscalData
	}

	// Check if invoice already exists for this order
	if existing, err := s.invoiceRepo.GetByOrderID(ctx, orderID); err == nil && existing != nil {
		return nil, ErrInvoiceAlreadyExists
	}

	// Check Focus client
	if s.focusClient == nil || !s.focusClient.Enabled() {
		return nil, ErrTransmitenotaNotConfigured
	}

	// Get next invoice number
	series := 1 // Default series
	number, err := s.invoiceRepo.GetNextNumber(ctx, company.ID, series)
	if err != nil {
		return nil, fmt.Errorf("failed to get next invoice number: %w", err)
	}

	// Create invoice entity
	invoice := fiscalinvoice.NewFiscalInvoice(company.ID, orderID, number, series)

	// Fetch order to get items and payment info
	orderModel, err := s.orderRepo.GetOrderById(ctx, orderID.String())
	if err != nil || orderModel == nil {
		return nil, ErrOrderNotFound
	}

	// Build NFC-e items from order using default food fiscal values
	nfceItems := make([]focusnfe.NFCeItem, 0)
	itemNumber := 1
	for _, group := range orderModel.GroupItems {
		for _, item := range group.Items {
			valorUnitario, _ := item.Price.Float64()
			valorTotal, _ := item.TotalPrice.Float64()

			nfceItem := focusnfe.NFCeItem{
				NumeroItem:             itemNumber,
				CodigoProduto:          item.ProductID.String()[:8], // First 8 chars of product ID
				Descricao:              item.Name,
				NCM:                    fiscalinvoice.DefaultFoodNCM,
				CFOP:                   fiscalinvoice.DefaultCFOP,
				UnidadeComercial:       fiscalinvoice.DefaultUnidade,
				QuantidadeComercial:    float64(item.Quantity),
				ValorUnitarioComercial: valorUnitario,
				ValorBruto:             valorTotal,
				ICMSOrigem:             fmt.Sprintf("%d", fiscalinvoice.DefaultOrigem),
				ICMSSituacaoTributaria: fiscalinvoice.GetCSOSNForRegime(settings.TaxRegime),
			}
			nfceItems = append(nfceItems, nfceItem)
			itemNumber++
		}
	}

	// Build payment info from order
	formasPagamento := make([]focusnfe.FormaPagamento, 0)
	for _, payment := range orderModel.Payments {
		valor, _ := payment.TotalPaid.Float64()
		forma := mapPaymentMethod(payment.Method)
		formasPagamento = append(formasPagamento, focusnfe.FormaPagamento{
			FormaPagamento: forma,
			ValorPagamento: valor,
		})
	}

	// If no payments yet, add "dinheiro" with total
	if len(formasPagamento) == 0 {
		valorTotal, _ := orderModel.TotalPayable.Float64()
		formasPagamento = append(formasPagamento, focusnfe.FormaPagamento{
			FormaPagamento: "01", // 01 = Dinheiro
			ValorPagamento: valorTotal,
		})
	}

	// Sanitize CNPJ to ensure only digits
	reg, _ := regexp.Compile("[^0-9]+")
	sanitizedCNPJ := reg.ReplaceAllString(settings.Cnpj, "")

	nfceRequest := &focusnfe.NFCeRequest{
		NaturezaOperacao:  "Venda ao Consumidor",
		DataEmissao:       time.Now().UTC().Format("2006-01-02T15:04:05-07:00"),
		Itens:             nfceItems,
		FormasPagamento:   formasPagamento,
		Numero:            fmt.Sprintf("%d", number),
		Serie:             fmt.Sprintf("%d", series),
		CNPJ:              sanitizedCNPJ,
		PresencaComprador: "1", // Operação presencial
	}

	// Use invoice ID as reference
	reference := invoice.ID.String()

	// Select correct token based on environment
	token := settings.TokenHomologation
	if s.focusClient.GetEnvironment() == "production" {
		token = settings.TokenProduction
	}

	// Emit NFC-e via Focus NFe API
	response, err := s.focusClient.EmitirNFCe(
		ctx,
		reference,
		nfceRequest,
		token,
	)

	if err != nil {
		// Mark as rejected
		invoice.Reject(err.Error())
		invoiceModel := &model.FiscalInvoice{}
		invoiceModel.FromDomain(invoice)
		_ = s.invoiceRepo.Create(ctx, invoiceModel)
		return nil, fmt.Errorf("failed to emit NFC-e: %w", err)
	}

	// Check response status
	// Focus NFe returns status like "autorizado", "processando", "erro_autorizacao"
	if response.Status == "erro_autorizacao" {
		// ... handle error
		errorMsg := response.Mensagem
		if len(response.Erros) > 0 {
			errorMsg += fmt.Sprintf(" %s", string(response.Erros))
		}
		invoice.Reject(errorMsg)
		invoiceModel := &model.FiscalInvoice{}
		invoiceModel.FromDomain(invoice)
		_ = s.invoiceRepo.Create(ctx, invoiceModel)
		return nil, fmt.Errorf("NFC-e rejected: %s", errorMsg)
	}

	// Note: If status is "processando", we might need to poll later.
	// For now, we save what we have. If it's authorized, we get paths.

	if response.Status == "autorizado" {
		// Mark as authorized
		invoice.Authorize(response.ChaveNFe, response.Protocolo, response.CaminhoXML, response.CaminhoPDF)
	} else {
		// Use "Processing" status if available in domain, or just save generic status.
		// Assuming current domain only has Pending, Authorized, Rejected, Cancelled.
		// If "processando", we might leave it as Pending (Created).
		// But we should update ChaveNFe if available.
	}

	// Save invoice
	invoiceModel := &model.FiscalInvoice{}
	invoiceModel.FromDomain(invoice)
	if err := s.invoiceRepo.Create(ctx, invoiceModel); err != nil {
		return nil, fmt.Errorf("failed to save invoice: %w", err)
	}

	// Register cost (R$ 0.10 per NFC-e)
	if s.usageCostService != nil && response.Status == "autorizado" {
		description := fmt.Sprintf("Emissão NFC-e #%d - Série %d", number, series)
		pricePerInvoice, _ := decimal.NewFromString(os.Getenv("PRICE_PER_NFCE"))

		cost := companyentity.NewUsageCost(company.ID, companyentity.CostTypeNFCe, pricePerInvoice, description, &invoice.ID)
		if err := s.usageCostService.RegisterUsageCost(ctx, cost); err != nil {
			// Log error but don't fail the emission
			fmt.Printf("Warning: failed to register NFC-e cost: %v\n", err)
		}
	}

	return invoice, nil
}

// ConsultarNFCe queries NFC-e status
func (s *Service) ConsultarNFCe(ctx context.Context, invoiceID uuid.UUID) (*fiscalinvoice.FiscalInvoice, error) {
	invoiceModel, err := s.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return nil, ErrInvoiceNotFound
	}

	invoice := invoiceModel.ToDomain()

	if s.focusClient == nil || !s.focusClient.Enabled() {
		return invoice, nil
	}

	settingsModel, err := s.fiscalSettingsRepo.GetByCompanyID(ctx, invoice.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get fiscal settings: %w", err)
	}

	settings := settingsModel.ToDomain()
	token := settings.TokenHomologation
	if s.focusClient.GetEnvironment() == "production" {
		token = settings.TokenProduction
	}

	response, err := s.focusClient.ConsultarNFCe(
		ctx,
		invoice.ID.String(),
		token,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to consult NFC-e: %w", err)
	}

	// Update status if authorized
	switch response.Status {
	case "autorizado":
		if invoice.Status != fiscalinvoice.StatusAuthorized {
			invoice.Authorize(response.ChaveNFe, response.Protocolo, response.CaminhoXML, response.CaminhoPDF)
		}
		// Ensure paths are updated
		if response.CaminhoXML != "" {
			invoice.XMLPath = response.CaminhoXML
		}
		if response.CaminhoPDF != "" {
			invoice.PDFPath = response.CaminhoPDF
		}
	case "erro_autorizacao", "cancelado":
		// handle other statuses
		if response.Status == "cancelado" {
			invoice.Cancel("Cancelado na SEFAZ") // Or keep original logic
		}
	}

	// Update model and save
	invoiceModel.FromDomain(invoice)
	if err := s.invoiceRepo.Update(ctx, invoiceModel); err != nil {
		return nil, fmt.Errorf("failed to update invoice: %w", err)
	}

	return invoice, nil
}

// CancelarNFCe cancels an NFC-e
func (s *Service) CancelarNFCe(ctx context.Context, invoiceID uuid.UUID, justify string) error {
	invoiceModel, err := s.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return ErrInvoiceNotFound
	}

	invoice := invoiceModel.ToDomain()

	if !invoice.CanBeCancelled() {
		return ErrCannotCancelInvoice
	}

	// Validate justification (minimum 15 characters as per SEFAZ rules)
	if len(justify) < 15 {
		return errors.New("justify must be at least 15 characters long")
	}

	// Cancel via Focus NFe API
	if s.focusClient == nil || !s.focusClient.Enabled() {
		return errors.New("focus client not enabled")
	}

	cancelRequest := &focusnfe.CancelamentoRequest{
		Justify: justify,
	}

	// Fetch settings to get token
	settingsModel, err := s.fiscalSettingsRepo.GetByCompanyID(ctx, invoice.CompanyID)
	if err != nil {
		return fmt.Errorf("failed to get fiscal settings: %w", err)
	}

	settings := settingsModel.ToDomain()
	token := settings.TokenHomologation
	if s.focusClient.GetEnvironment() == "production" {
		token = settings.TokenProduction
	}

	if err = s.focusClient.CancelarNFCe(
		ctx,
		invoice.ID.String(), // Use reference
		cancelRequest,
		token,
	); err != nil {
		return fmt.Errorf("failed to cancel NFC-e: %w", err)
	}

	// Mark as cancelled
	invoice.Cancel(justify)
	invoiceModel.FromDomain(invoice)

	if err := s.invoiceRepo.Update(ctx, invoiceModel); err != nil {
		return fmt.Errorf("failed to update invoice: %w", err)
	}

	return nil
}

// ListInvoices lists fiscal invoices for the company
func (s *Service) ListInvoices(ctx context.Context, page, perPage int) ([]*fiscalinvoice.FiscalInvoice, int, error) {
	companyModel, err := s.companyRepo.GetCompany(ctx)
	if err != nil {
		return nil, 0, err
	}

	invoiceModels, total, err := s.invoiceRepo.List(ctx, companyModel.ID, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	invoices := make([]*fiscalinvoice.FiscalInvoice, len(invoiceModels))
	for i, model := range invoiceModels {
		invoices[i] = model.ToDomain()
	}

	return invoices, total, nil
}

// mapPaymentMethod converts order payment method to NFC-e payment code
// https://nfe.io/docs/campos/formas-de-pagamento-nf-e/
func mapPaymentMethod(method string) string {
	switch method {
	case "Dinheiro":
		return "01" // Dinheiro
	case "Cheque":
		return "02" // Cheque
	case "Visa", "MasterCard", "American Express", "Elo", "Diners Club", "Hipercard":
		return "03" // Cartão de Crédito
	case "Visa Electron", "Maestro":
		return "04" // Cartão de Débito
	case "VR", "Ticket", "Alelo":
		return "04" // Vale Refeição/Alimentação (tratado como débito)
	case "PIX":
		return "17" // PIX
	case "PayPal":
		return "99" // Outros
	default:
		return "99" // Outros
	}
}
