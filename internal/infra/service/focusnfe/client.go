package focusnfe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	productionBaseURL   = "https://api.focusnfe.com.br"
	homologationBaseURL = "https://homologacao.focusnfe.com.br"
	defaultTimeout      = 30 * time.Second
)

// Client wraps HTTP client for Focus NFe API
type Client struct {
	baseURL             string
	httpClient          *http.Client
	mainProductionToken string
	environment         string // "production" or "homologation"
}

// CompanyRegistryRequest represents the payload to register a company
type CompanyRegistryRequest struct {
	Nome                    string `json:"nome"`
	NomeFantasia            string `json:"nome_fantasia"`
	InscricaoEstadual       string `json:"inscricao_estadual"`
	InscricaoMunicipal      string `json:"inscricao_municipal,omitempty"`
	CNPJ                    string `json:"cnpj"`
	RegimeTributario        string `json:"regime_tributario"` // 1=Simples Nacional, 3=Regime Normal
	Email                   string `json:"email"`
	Telefone                string `json:"telefone"`
	Logradouro              string `json:"logradouro"`
	Numero                  string `json:"numero"`
	Complemento             string `json:"complemento,omitempty"`
	Bairro                  string `json:"bairro"`
	CEP                     string `json:"cep"`
	Municipio               string `json:"municipio"`
	UF                      string `json:"uf"`
	DiscriminaImpostos      bool   `json:"discrimina_impostos"`
	EnviarEmailDestinatario bool   `json:"enviar_email_destinatario"`
	CscNfceProducao         string `json:"csc_nfce_producao,omitempty"`
	IdTokenNfceProducao     string `json:"id_token_nfce_producao,omitempty"`
	CscNfceHomologacao      string `json:"csc_nfce_homologacao,omitempty"`
	IdTokenNfceHomologacao  string `json:"id_token_nfce_homologacao,omitempty"`
}

// CompanyRegistryResponse represents the response from company registration
type CompanyRegistryResponse struct {
	ID                int64           `json:"id"`
	Mensagem          string          `json:"mensagem,omitempty"`
	TokenProduction   string          `json:"token_producao"`
	TokenHomologation string          `json:"token_homologacao"`
	Errors            json.RawMessage `json:"erros,omitempty"` // Can be string or array
}

// NewClient creates a new Focus NFe API client
func NewClient() *Client {

	token := strings.TrimSpace(os.Getenv("FOCUS_NFE_API_KEY"))

	environment := os.Getenv("FOCUS_NFE_ENV")
	if environment == "" {
		environment = "homologation"
	}

	baseURL := homologationBaseURL
	if environment == "production" {
		baseURL = productionBaseURL
	}

	timeout := defaultTimeout
	if timeoutStr := os.Getenv("FOCUS_NFE_TIMEOUT"); timeoutStr != "" {
		if d, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = d
		}
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		mainProductionToken: token,
		environment:         environment,
	}
}

// CadastrarEmpresa registers a new company in Focus NFe
func (c *Client) CadastrarEmpresa(ctx context.Context, req *CompanyRegistryRequest) (*CompanyRegistryResponse, error) {
	// The Companies API operates EXCLUSIVELY in the production environment.
	// https://focusnfe.com.br/doc/#empresas-ambientes
	endpoint := "/v2/empresas"

	// If we are in homologation, we should use dry_run to simulate
	if c.environment != "production" {
		endpoint += "?dry_run=1"
	}

	// For this specific call, we MUST use the production URL, regardless of the environment
	// We'll temporarily override the base URL for this request
	originalBaseURL := c.baseURL
	c.baseURL = productionBaseURL
	defer func() { c.baseURL = originalBaseURL }()

	resp := &CompanyRegistryResponse{}
	if err := c.doRequest(ctx, "POST", endpoint, req, resp, c.mainProductionToken); err != nil {
		return nil, err
	}

	return resp, nil
}

// doRequest performs HTTP request to Focus NFe API
func (c *Client) doRequest(ctx context.Context, method, endpoint string, reqBody, respBody interface{}, token string) error {
	url := c.baseURL + endpoint

	var body io.Reader
	if reqBody != nil {
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Basic Auth with Token
	req.SetBasicAuth(token, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		// Try to parse error message if possible
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	if respBody != nil {
		if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// NFCeRequest represents the request to emit NFC-e
type NFCeRequest struct {
	NaturezaOperacao string          `json:"natureza_operacao"`
	DataEmissao      string          `json:"data_emissao,omitempty"`
	Itens            []NFCeItem      `json:"items"` // Note: Focus might use "items" or "itens" depending on version. Usually "items".
	Cliente          *NFCeClient     `json:"cliente,omitempty"`
	Pagamento        *NFCePayment    `json:"pagamento,omitempty"` // Or FormasPagamento
	FormasPagamento  []PaymentMethod `json:"formas_pagamento,omitempty"`
	// Additional fields
	Serie             string `json:"serie,omitempty"`
	Numero            string `json:"numero,omitempty"`
	PresencaComprador string `json:"presenca_comprador,omitempty"` // 1=Presencial
	CNPJ              string `json:"cnpj_emitente,omitempty"`      // sometimes used if token covers multiple
}

type NFCeItem struct {
	NumeroItem             int     `json:"numero_item"`
	CodigoProduto          string  `json:"codigo_produto"`
	Descricao              string  `json:"descricao"`
	CFOP                   string  `json:"cfop"` // 5102
	UnidadeComercial       string  `json:"unidade_comercial"`
	QuantidadeComercial    float64 `json:"quantidade_comercial"`
	ValorUnitarioComercial float64 `json:"valor_unitario_comercial"`
	ValorBruto             float64 `json:"valor_bruto"` // Qty * UnitPrice
	NCM                    string  `json:"ncm"`
	ICMSOrigem             string  `json:"icms_origem"`              // 0
	ICMSSituacaoTributaria string  `json:"icms_situacao_tributaria"` // 102, etc
	// PIS/COFINS usually needed too
}

type NFCeClient struct {
	CPF   string `json:"cpf,omitempty"`
	Nome  string `json:"nome_completo,omitempty"`
	Email string `json:"email,omitempty"`
}

type PaymentMethod struct {
	FormaPagamento string  `json:"forma_pagamento"` // 01, 03...
	ValorPagamento float64 `json:"valor_pagamento"`
}

// Keep older struct for compat if needed, but we are building new.
type NFCePayment struct {
	FormasPagamento []PaymentMethod `json:"formas_pagamento"`
	Troco           float64         `json:"troco,omitempty"`
}

type NFCeResponse struct {
	Status     string          `json:"status"` // authorized, processando, erro_autorizacao
	CaminhoXML string          `json:"caminho_xml_nota_fiscal"`
	CaminhoPDF string          `json:"caminho_danfe"`
	ChaveNFe   string          `json:"chave_nfe"`
	Numero     interface{}     `json:"numero"` // int or string
	Serie      interface{}     `json:"serie"`
	Protocolo  string          `json:"protocolo"`
	Mensagem   string          `json:"mensagem_sefaz,omitempty"`
	Erros      json.RawMessage `json:"erros,omitempty"`
}

// CancelRequest
type CancelRequest struct {
	Justificativa string `json:"justificativa"`
}

// EmitNFCe emits a new NFC-e
// Reference maps "reference" to our internal ID to query later if stuck in processing
func (c *Client) EmitNFCe(ctx context.Context, reference string, req *NFCeRequest, token string) (*NFCeResponse, error) {
	// POST /v2/nfce?ref=reference
	endpoint := fmt.Sprintf("/v2/nfce?ref=%s", reference)

	// If dry_run is requested (e.g. implicitly by environment), append to query
	// User requested "dry_run=1" on all permitted endpoints.
	if c.environment != "production" {
		endpoint += "&dry_run=1"
	}

	resp := &NFCeResponse{}
	if err := c.doRequest(ctx, "POST", endpoint, req, resp, token); err != nil {
		return nil, err
	}

	return resp, nil
}

// SearchNFCe queries NFC-e. Can query by reference.
func (c *Client) SearchNFCe(ctx context.Context, reference string, token string) (*NFCeResponse, error) {
	endpoint := fmt.Sprintf("/v2/nfce/%s", reference)

	resp := &NFCeResponse{}
	if err := c.doRequest(ctx, "GET", endpoint, nil, resp, token); err != nil {
		return nil, err
	}

	return resp, nil
}

// CancelNFCe cancels an NFC-e using its reference (or key)
func (c *Client) CancelNFCe(ctx context.Context, reference string, req *CancelRequest, token string) error {
	endpoint := fmt.Sprintf("/v2/nfce/%s", reference)

	// Cancellation in Focus NFe: DELETE /v2/nfce/{ref} with body?
	// Or POST /v2/nfce/{ref}/cancelar?
	// Documentation says DELETE /v2/nfce/{ref} cancels it.

	// Usually DELETE accepts body with justificativa.
	// Check if doRequest supports body in DELETE.

	resp := &struct {
		Status string `json:"status"`
	}{}
	if err := c.doRequest(ctx, "DELETE", endpoint, req, resp, token); err != nil {
		return err
	}

	return nil
}

// Enabled checks if client is properly configured with credentials
func (c *Client) Enabled() bool {
	return c != nil && c.baseURL != "" && c.mainProductionToken != ""
}

// GetEnvironment returns the current environment
func (c *Client) GetEnvironment() string {
	return c.environment
}
