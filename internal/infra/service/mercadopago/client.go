package mercadopagoservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	defaultAPIBaseURL      = "https://api.mercadopago.com"
	defaultMonthlyPriceBRL = 99.90
)

type Client struct {
	httpClient      *http.Client
	baseURL         string
	accessToken     string
	integratorID    string
	notificationURL string
	webhookSecret   string
	monthlyPrice    float64
}

type PreferenceRequest struct {
	Title   string
	Company string
	Months  int
	Price   float64
	Schema  string
	ID      string
}

type PreferenceResponse struct {
	ID               string `json:"id"`
	InitPoint        string `json:"init_point"`
	SandboxInitPoint string `json:"sandbox_init_point"`
}

type paymentResponse struct {
	ID                int64           `json:"id"`
	Status            string          `json:"status"`
	StatusDetail      string          `json:"status_detail"`
	CurrencyID        string          `json:"currency_id"`
	TransactionAmount float64         `json:"transaction_amount"`
	DateApproved      *time.Time      `json:"date_approved"`
	ExternalReference string          `json:"external_reference"`
	Metadata          json.RawMessage `json:"metadata"`
}

type PaymentMetadata struct {
	CompanyID  string `json:"company_id"`
	SchemaName string `json:"schema_name"`
	Months     int    `json:"months"`
}

type PaymentDetails struct {
	ID                string
	Status            string
	CurrencyID        string
	TransactionAmount float64
	DateApproved      *time.Time
	ExternalReference string
	Metadata          PaymentMetadata
}

func NewClient() *Client {
	accessToken := os.Getenv("MP_ACCESS_TOKEN")
	baseURL := os.Getenv("MP_API_BASE_URL")
	if baseURL == "" {
		baseURL = defaultAPIBaseURL
	}
	notificationURL := os.Getenv("MP_WEBHOOK_URL")
	webhookSecret := os.Getenv("MP_WEBHOOK_SECRET")
	integratorID := os.Getenv("MP_INTEGRATOR_ID")

	monthlyPrice := defaultMonthlyPriceBRL
	if raw := os.Getenv("MP_SUBSCRIPTION_PRICE"); raw != "" {
		if parsed, err := strconv.ParseFloat(raw, 64); err == nil && parsed > 0 {
			monthlyPrice = parsed
		}
	}

	return &Client{
		httpClient:      &http.Client{Timeout: 15 * time.Second},
		baseURL:         baseURL,
		accessToken:     accessToken,
		integratorID:    integratorID,
		notificationURL: notificationURL,
		webhookSecret:   webhookSecret,
		monthlyPrice:    monthlyPrice,
	}
}

func (c *Client) MonthlyPrice() float64 {
	return c.monthlyPrice
}

func (c *Client) WebhookSecret() string {
	return c.webhookSecret
}

func (c *Client) NotificationURL() string {
	return c.notificationURL
}

func (c *Client) haveCredentials() bool {
	return c != nil && c.accessToken != ""
}

func (c *Client) CreateSubscriptionPreference(ctx context.Context, req *PreferenceRequest) (*PreferenceResponse, error) {
	if c == nil || !c.haveCredentials() {
		return nil, fmt.Errorf("mercado pago credentials missing")
	}

	months := req.Months
	if months <= 0 {
		months = 1
	}

	price := req.Price
	if price <= 0 {
		price = c.monthlyPrice
	}

	item := map[string]interface{}{
		"title":       fmt.Sprintf("Mensalidade %s", req.Company),
		"description": fmt.Sprintf("Plano de assinatura (%d mês(es))", months),
		"quantity":    1,
		"unit_price":  price * float64(months),
		"currency_id": "BRL",
	}

	body := map[string]interface{}{
		"items":              []interface{}{item},
		"external_reference": req.ID,
		"notification_url":   c.notificationURL,
		"auto_return":        "approved",
		"metadata": map[string]interface{}{
			"company_id":  req.ID,
			"schema_name": req.Schema,
			"months":      months,
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/checkout/preferences", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	c.applyHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, c.bodyError(resp)
	}

	out := &PreferenceResponse{}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) GetPayment(ctx context.Context, id string) (*PaymentDetails, error) {
	if c == nil || !c.haveCredentials() {
		return nil, fmt.Errorf("mercado pago credentials missing")
	}

	url := fmt.Sprintf("%s/v1/payments/%s", c.baseURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	c.applyHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, c.bodyError(resp)
	}

	raw := &paymentResponse{}
	if err := json.NewDecoder(resp.Body).Decode(raw); err != nil {
		return nil, err
	}

	meta := PaymentMetadata{}
	if len(raw.Metadata) > 0 {
		_ = json.Unmarshal(raw.Metadata, &meta)
	}

	return &PaymentDetails{
		ID:                fmt.Sprintf("%d", raw.ID),
		Status:            raw.Status,
		CurrencyID:        raw.CurrencyID,
		TransactionAmount: raw.TransactionAmount,
		DateApproved:      raw.DateApproved,
		ExternalReference: raw.ExternalReference,
		Metadata:          meta,
	}, nil
}

func (c *Client) applyHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")
	if c.integratorID != "" {
		req.Header.Set("X-Integrator-Id", c.integratorID)
	}
}

func (c *Client) bodyError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("mercado pago error: %s (%s)", resp.Status, string(body))
}
