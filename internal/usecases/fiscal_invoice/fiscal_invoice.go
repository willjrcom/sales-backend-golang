package fiscalinvoiceusecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	fiscalinvoice "github.com/willjrcom/sales-backend-go/internal/domain/fiscal_invoice"
	"github.com/willjrcom/sales-backend-go/internal/infra/repository/model"
	"github.com/willjrcom/sales-backend-go/internal/infra/service/transmitenota"
	companyusecases "github.com/willjrcom/sales-backend-go/internal/usecases/company"
)

var (
	ErrFiscalNotEnabled           = errors.New("fiscal invoice functionality is not enabled for this company")
	ErrMissingFiscalData          = errors.New("company is missing required fiscal data (IE, regime tributário)")
	ErrTransmitenotaNotConfigured = errors.New("transmitenota client is not configured")
	ErrInvoiceAlreadyExists       = errors.New("invoice already exists for this order")
	ErrInvoiceNotFound            = errors.New("fiscal invoice not found")
	ErrCannotCancelInvoice        = errors.New("invoice cannot be cancelled (not authorized or already cancelled)")
	ErrOrderNotFound              = errors.New("order not found")
)

type Service struct {
	invoiceRepo         model.FiscalInvoiceRepository
	companyRepo         model.CompanyRepository
	orderRepo           model.OrderRepository
	usageCostService    *companyusecases.UsageCostService
	transmitenotaClient *transmitenota.Client
}

func NewService(
	invoiceRepo model.FiscalInvoiceRepository,
	companyRepo model.CompanyRepository,
	orderRepo model.OrderRepository,
	usageCostService *companyusecases.UsageCostService,
	transmitenotaClient *transmitenota.Client,
) *Service {
	return &Service{
		invoiceRepo:         invoiceRepo,
		companyRepo:         companyRepo,
		orderRepo:           orderRepo,
		usageCostService:    usageCostService,
		transmitenotaClient: transmitenotaClient,
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

	// Validate fiscal is enabled
	if !company.FiscalEnabled {
		return nil, ErrFiscalNotEnabled
	}

	// Validate required fiscal data
	if company.InscricaoEstadual == "" || company.RegimeTributario == 0 {
		return nil, ErrMissingFiscalData
	}

	// Check if invoice already exists for this order
	if existing, err := s.invoiceRepo.GetByOrderID(ctx, orderID); err == nil && existing != nil {
		return nil, ErrInvoiceAlreadyExists
	}

	// Check Transmitenota client
	if s.transmitenotaClient == nil || !s.transmitenotaClient.Enabled() {
		return nil, ErrTransmitenotaNotConfigured
	}

	// Get next invoice number
	serie := 1 // Default series
	numero, err := s.invoiceRepo.GetNextNumber(ctx, company.ID, serie)
	if err != nil {
		return nil, fmt.Errorf("failed to get next invoice number: %w", err)
	}

	// Create invoice entity
	invoice := fiscalinvoice.NewFiscalInvoice(company.ID, orderID, numero, serie)

	// Fetch order to get items and payment info
	orderModel, err := s.orderRepo.GetOrderById(ctx, orderID.String())
	if err != nil || orderModel == nil {
		return nil, ErrOrderNotFound
	}

	// Build NFC-e items from order using default food fiscal values
	nfceItems := make([]transmitenota.NFCeItem, 0)
	itemNumber := 1
	for _, group := range orderModel.GroupItems {
		for _, item := range group.Items {
			valorUnitario, _ := item.Price.Float64()
			valorTotal, _ := item.TotalPrice.Float64()

			nfceItem := transmitenota.NFCeItem{
				Numero:        itemNumber,
				Codigo:        item.ProductID.String()[:8], // First 8 chars of product ID
				Descricao:     item.Name,
				NCM:           fiscalinvoice.DefaultFoodNCM,
				CFOP:          fiscalinvoice.DefaultCFOP,
				Unidade:       fiscalinvoice.DefaultUnidade,
				Quantidade:    float64(item.Quantity),
				ValorUnitario: valorUnitario,
				ValorTotal:    valorTotal,
				ICMS: &transmitenota.ICMS{
					Situacao: fiscalinvoice.GetCSOSNForRegime(company.RegimeTributario),
					Origem:   fmt.Sprintf("%d", fiscalinvoice.DefaultOrigem),
					Aliquota: fiscalinvoice.DefaultAliquotaICMS,
				},
			}
			nfceItems = append(nfceItems, nfceItem)
			itemNumber++
		}
	}

	// Build payment info from order
	nfcePagamentos := make([]transmitenota.NFCePagamento, 0)
	for _, payment := range orderModel.Payments {
		valor, _ := payment.TotalPaid.Float64()
		forma := mapPaymentMethod(payment.Method)
		nfcePagamentos = append(nfcePagamentos, transmitenota.NFCePagamento{
			Forma: forma,
			Valor: valor,
		})
	}

	// If no payments yet, add "dinheiro" with total
	if len(nfcePagamentos) == 0 {
		valorTotal, _ := orderModel.TotalPayable.Float64()
		nfcePagamentos = append(nfcePagamentos, transmitenota.NFCePagamento{
			Forma: "01", // 01 = Dinheiro
			Valor: valorTotal,
		})
	}

	nfceRequest := &transmitenota.NFCeRequest{
		CNPJ:              company.Cnpj,
		InscricaoEstadual: company.InscricaoEstadual,
		RegimeTributario:  company.RegimeTributario,
		Numero:            numero,
		Serie:             serie,
		DataEmissao:       time.Now().Format("2006-01-02T15:04:05-07:00"),
		Items:             nfceItems,
		Pagamento:         nfcePagamentos,
	}

	// Emit NFC-e via Transmitenota API
	response, err := s.transmitenotaClient.EmitirNFCe(
		ctx,
		nfceRequest,
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
	if response.Status != "autorizado" && response.Status != "authorized" {
		errorMsg := response.ErroMensagem
		if errorMsg == "" {
			errorMsg = response.Mensagem
		}
		invoice.Reject(errorMsg)
		invoiceModel := &model.FiscalInvoice{}
		invoiceModel.FromDomain(invoice)
		_ = s.invoiceRepo.Create(ctx, invoiceModel)
		return nil, fmt.Errorf("NFC-e rejected: %s", errorMsg)
	}

	// Mark as authorized
	invoice.Authorize(response.ChaveAcesso, response.Protocolo, response.XMLPath, response.PDFPath)

	// Save invoice
	invoiceModel := &model.FiscalInvoice{}
	invoiceModel.FromDomain(invoice)
	if err := s.invoiceRepo.Create(ctx, invoiceModel); err != nil {
		return nil, fmt.Errorf("failed to save invoice: %w", err)
	}

	// Register cost (R$ 0.10 per NFC-e)
	if s.usageCostService != nil {
		description := fmt.Sprintf("Emissão NFC-e #%d - Série %d", numero, serie)
		if err := s.usageCostService.RegisterNFCeCost(ctx, company.ID, invoice.ID, description); err != nil {
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

	// If already has chave_acesso, query Transmitenota
	if invoice.ChaveAcesso != "" && s.transmitenotaClient != nil {
		response, err := s.transmitenotaClient.ConsultarNFCe(
			ctx,
			invoice.ChaveAcesso,
		)

		if err == nil && response != nil {
			// Update paths if changed
			if response.XMLPath != "" {
				invoice.XMLPath = response.XMLPath
			}
			if response.PDFPath != "" {
				invoice.PDFPath = response.PDFPath
			}

			// Update model and save
			invoiceModel.FromDomain(invoice)
			_ = s.invoiceRepo.Update(ctx, invoiceModel)
		}
	}

	return invoice, nil
}

// CancelarNFCe cancels an NFC-e
func (s *Service) CancelarNFCe(ctx context.Context, invoiceID uuid.UUID, justificativa string) error {
	invoiceModel, err := s.invoiceRepo.GetByID(ctx, invoiceID)
	if err != nil {
		return ErrInvoiceNotFound
	}

	invoice := invoiceModel.ToDomain()

	if !invoice.CanBeCancelled() {
		return ErrCannotCancelInvoice
	}

	// Validate justification (minimum 15 characters as per SEFAZ rules)
	if len(justificativa) < 15 {
		return errors.New("justificativa deve ter no mínimo 15 caracteres")
	}

	// Cancel via Transmitenota API
	if s.transmitenotaClient != nil && s.transmitenotaClient.Enabled() {
		cancelRequest := &transmitenota.CancelamentoRequest{
			ChaveAcesso:   invoice.ChaveAcesso,
			Justificativa: justificativa,
		}

		_, err := s.transmitenotaClient.CancelarNFCe(
			ctx,
			cancelRequest,
		)

		if err != nil {
			return fmt.Errorf("failed to cancel NFC-e: %w", err)
		}
	}

	// Mark as cancelled
	invoice.Cancel(justificativa)
	invoiceModel.FromDomain(invoice)

	if err := s.invoiceRepo.Update(ctx, invoiceModel); err != nil {
		return fmt.Errorf("failed to update invoice: %w", err)
	}

	// TODO: Register refund cost (negative amount)
	// if s.usageCostService != nil {
	// 	s.usageCostService.RegisterUsageCost(...)
	// }

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
