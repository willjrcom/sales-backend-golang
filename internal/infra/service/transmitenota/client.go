package transmitenota

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	defaultBaseURL = "https://api.transmitenota.com.br"
	defaultTimeout = 30 * time.Second
)

// Client wraps HTTP client for Transmitenota API
type Client struct {
	baseURL    string
	httpClient *http.Client
	sandbox    bool
	usuario    string
	senha      string
}

// NFCeRequest represents the request to emit NFC-e
type NFCeRequest struct {
	CNPJ              string          `json:"cnpj"`
	InscricaoEstadual string          `json:"inscricao_estadual"`
	RegimeTributario  int             `json:"regime_tributario"`
	Numero            int             `json:"numero"`
	Serie             int             `json:"serie"`
	DataEmissao       string          `json:"data_emissao"`
	Items             []NFCeItem      `json:"items"`
	Cliente           *NFCeCliente    `json:"cliente,omitempty"`
	Pagamento         []NFCePagamento `json:"pagamento"`
	InformacoesAdic   string          `json:"informacoes_adicionais,omitempty"`
}

type NFCeItem struct {
	Numero        int     `json:"numero"`
	Codigo        string  `json:"codigo"`
	Descricao     string  `json:"descricao"`
	NCM           string  `json:"ncm"`
	CFOP          string  `json:"cfop"`
	Unidade       string  `json:"unidade"`
	Quantidade    float64 `json:"quantidade"`
	ValorUnitario float64 `json:"valor_unitario"`
	ValorTotal    float64 `json:"valor_total"`
	ICMS          *ICMS   `json:"icms,omitempty"`
}

type ICMS struct {
	Situacao  string  `json:"situacao_tributaria"`
	Origem    string  `json:"origem"`
	Aliquota  float64 `json:"aliquota,omitempty"`
	ValorBase float64 `json:"valor_base_calculo,omitempty"`
	ValorICMS float64 `json:"valor_icms,omitempty"`
}

type NFCeCliente struct {
	CPF  string `json:"cpf,omitempty"`
	CNPJ string `json:"cnpj,omitempty"`
	Nome string `json:"nome,omitempty"`
}

type NFCePagamento struct {
	Forma string  `json:"forma"`
	Valor float64 `json:"valor"`
}

// NFCeResponse represents the response from NFC-e emission
type NFCeResponse struct {
	Status       string `json:"status"`
	ChaveAcesso  string `json:"chave_acesso"`
	Numero       int    `json:"numero"`
	Serie        int    `json:"serie"`
	DataEmissao  string `json:"data_emissao"`
	Protocolo    string `json:"protocolo"`
	XMLPath      string `json:"xml_path"`
	PDFPath      string `json:"pdf_path"`
	Mensagem     string `json:"mensagem,omitempty"`
	ErroMensagem string `json:"erro_mensagem,omitempty"`
}

// CancelamentoRequest represents cancellation request
type CancelamentoRequest struct {
	ChaveAcesso   string `json:"chave_acesso"`
	Justificativa string `json:"justificativa"`
}

// CancelamentoResponse represents cancellation response
type CancelamentoResponse struct {
	Status       string `json:"status"`
	ChaveAcesso  string `json:"chave_acesso"`
	Protocolo    string `json:"protocolo"`
	Mensagem     string `json:"mensagem,omitempty"`
	ErroMensagem string `json:"erro_mensagem,omitempty"`
}

// ConsultaResponse represents query response
type ConsultaResponse struct {
	Status      string `json:"status"`
	ChaveAcesso string `json:"chave_acesso"`
	Numero      int    `json:"numero"`
	Serie       int    `json:"serie"`
	DataEmissao string `json:"data_emissao"`
	Protocolo   string `json:"protocolo"`
	XMLPath     string `json:"xml_path"`
	PDFPath     string `json:"pdf_path"`
	Situacao    string `json:"situacao"`
	Mensagem    string `json:"mensagem,omitempty"`
}

// NewClient creates a new Transmitenota API client with credentials from ENV
func NewClient() *Client {
	baseURL := os.Getenv("TRANSMITENOTA_API_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	sandbox := os.Getenv("TRANSMITENOTA_SANDBOX") == "true"

	timeout := defaultTimeout
	if timeoutStr := os.Getenv("TRANSMITENOTA_TIMEOUT"); timeoutStr != "" {
		if d, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = d
		}
	}

	// Load credentials from ENV
	usuario := os.Getenv("TRANSMITENOTA_USUARIO")
	senha := os.Getenv("TRANSMITENOTA_SENHA")

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		sandbox: sandbox,
		usuario: usuario,
		senha:   senha,
	}
}

// EmitirNFCe emits a new NFC-e
func (c *Client) EmitirNFCe(ctx context.Context, req *NFCeRequest) (*NFCeResponse, error) {
	endpoint := "/nfce/emitir"
	if c.sandbox {
		endpoint = "/sandbox" + endpoint
	}

	resp := &NFCeResponse{}
	if err := c.doRequest(ctx, "POST", endpoint, req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ConsultarNFCe queries NFC-e status
func (c *Client) ConsultarNFCe(ctx context.Context, chaveAcesso string) (*ConsultaResponse, error) {
	endpoint := fmt.Sprintf("/nfce/consultar/%s", chaveAcesso)
	if c.sandbox {
		endpoint = "/sandbox" + endpoint
	}

	resp := &ConsultaResponse{}
	if err := c.doRequest(ctx, "GET", endpoint, nil, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// CancelarNFCe cancels an NFC-e
func (c *Client) CancelarNFCe(ctx context.Context, req *CancelamentoRequest) (*CancelamentoResponse, error) {
	endpoint := "/nfce/cancelar"
	if c.sandbox {
		endpoint = "/sandbox" + endpoint
	}

	resp := &CancelamentoResponse{}
	if err := c.doRequest(ctx, "POST", endpoint, req, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// doRequest performs HTTP request to Transmitenota API
func (c *Client) doRequest(ctx context.Context, method, endpoint string, reqBody, respBody interface{}) error {
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
	req.Header.Set("Accept", "application/json")

	// Basic auth with credentials from ENV
	if c.usuario != "" && c.senha != "" {
		req.SetBasicAuth(c.usuario, c.senha)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	if respBody != nil {
		if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Enabled checks if client is properly configured with credentials
func (c *Client) Enabled() bool {
	return c != nil && c.baseURL != "" && c.usuario != "" && c.senha != ""
}
