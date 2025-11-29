package mercadopagoservice

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

const (
	defaultMonthlyPriceBRL  = 99.90
	signatureMaxSkewMinutes = 5 * time.Minute
)

var (
	// ErrSignatureHeaderMissing indicates the webhook request did not include the signature header.
	ErrSignatureHeaderMissing = errors.New("mercado pago signature header is missing")
	// ErrSignatureTimestampInvalid indicates the timestamp supplied on the signature is invalid or expired.
	ErrSignatureTimestampInvalid = errors.New("mercado pago signature timestamp is invalid or expired")
	// ErrSignatureMismatch indicates the computed signature did not match the one supplied on the webhook.
	ErrSignatureMismatch = errors.New("mercado pago signature mismatch")
	// ErrWebhookSecretNotConfigured indicates the client is missing the webhook secret/notification URL configuration.
	ErrWebhookSecretNotConfigured = errors.New("mercado pago webhook secret is not configured")
)

// Client wraps Mercado Pago SDK clients with project specific helpers.
type Client struct {
	preferenceClient preference.Client
	paymentClient    payment.Client
	notificationURL  string
	webhookSecret    string
	monthlyPrice     float64
	successURL       string
	pendingURL       string
	failureURL       string
}

// PreferenceRequest wraps the information required to create subscription preferences.
type PreferenceRequest struct {
	Title   string
	Company string
	Months  int
	Price   float64
	Schema  string
	ID      string
}

// PreferenceResponse mirrors the minimal data needed by the application when creating a preference.
type PreferenceResponse struct {
	ID               string `json:"id"`
	InitPoint        string `json:"init_point"`
	SandboxInitPoint string `json:"sandbox_init_point"`
}

// PaymentMetadata stores metadata persisted on Mercado Pago payments.
type PaymentMetadata struct {
	CompanyID  string `json:"company_id"`
	SchemaName string `json:"schema_name"`
	Months     int    `json:"months"`
}

// PaymentDetails wraps the fields we use when reconciling payments.
type PaymentDetails struct {
	ID                string
	Status            string
	CurrencyID        string
	TransactionAmount float64
	DateApproved      *time.Time
	ExternalReference string
	Metadata          PaymentMetadata
}

// NewClient configures the Mercado Pago SDK clients using environment variables.
func NewClient() *Client {
	monthlyPrice := defaultMonthlyPriceBRL
	if raw := os.Getenv("MP_SUBSCRIPTION_PRICE"); raw != "" {
		if parsed, err := strconv.ParseFloat(raw, 64); err == nil && parsed > 0 {
			monthlyPrice = parsed
		}
	}

	client := &Client{
		notificationURL: os.Getenv("MP_WEBHOOK_URL"),
		webhookSecret:   os.Getenv("MP_WEBHOOK_SECRET"),
		monthlyPrice:    monthlyPrice,
		successURL:      os.Getenv("MP_SUCCESS_URL"),
		pendingURL:      os.Getenv("MP_PENDING_URL"),
		failureURL:      os.Getenv("MP_FAILURE_URL"),
	}

	accessToken := os.Getenv("MP_ACCESS_TOKEN")
	if accessToken == "" {
		return client
	}

	opts := []config.Option{}
	if integratorID := os.Getenv("MP_INTEGRATOR_ID"); integratorID != "" {
		opts = append(opts, config.WithIntegratorID(integratorID))
	}

	cfg, err := config.New(accessToken, opts...)
	if err != nil {
		log.Printf("mercadopago: failed to initialize SDK config: %v", err)
		return client
	}

	client.preferenceClient = preference.NewClient(cfg)
	client.paymentClient = payment.NewClient(cfg)

	return client
}

// Enabled indicates whether the SDK clients are ready for use.
func (c *Client) Enabled() bool {
	if c == nil {
		return false
	}
	if c.preferenceClient == nil || c.paymentClient == nil {
		return false
	}
	if c.notificationURL == "" || c.successURL == "" || c.webhookSecret == "" {
		return false
	}
	return true
}

// MonthlyPrice returns the configured default monthly price.
func (c *Client) MonthlyPrice() float64 {
	if c == nil {
		return defaultMonthlyPriceBRL
	}
	return c.monthlyPrice
}

// WebhookSecret returns the configured webhook secret.
func (c *Client) WebhookSecret() string {
	if c == nil {
		return ""
	}
	return c.webhookSecret
}

// NotificationURL returns the URL we registered for Mercado Pago webhooks.
func (c *Client) NotificationURL() string {
	if c == nil {
		return ""
	}
	return c.notificationURL
}

// CreateSubscriptionPreference creates a checkout preference using the official SDK.
func (c *Client) CreateSubscriptionPreference(ctx context.Context, req *PreferenceRequest) (*PreferenceResponse, error) {
	if c == nil || !c.Enabled() {
		return nil, fmt.Errorf("mercado pago client is not configured")
	}

	if c.notificationURL == "" {
		return nil, fmt.Errorf("mercado pago notification url is not configured")
	}

	if c.successURL == "" {
		return nil, fmt.Errorf("mercado pago success url is not configured (set MP_SUCCESS_URL)")
	}

	months := req.Months
	if months <= 0 {
		months = 1
	}

	price := req.Price
	if price <= 0 {
		price = c.monthlyPrice
	}

	item := preference.ItemRequest{
		Title:       fmt.Sprintf("Mensalidade %s", req.Company),
		Description: fmt.Sprintf("Plano de assinatura (%d mês(es))", months),
		Quantity:    1,
		UnitPrice:   price * float64(months),
		CurrencyID:  "BRL",
	}

	metadata := map[string]any{
		"company_id":  req.ID,
		"schema_name": req.Schema,
		"months":      months,
	}

	prefRequest := preference.Request{
		Items:             []preference.ItemRequest{item},
		ExternalReference: req.ID,
		NotificationURL:   c.notificationURL,
		AutoReturn:        "approved",
		Metadata:          metadata,
		BackURLs: &preference.BackURLsRequest{
			Success: c.successURL,
			Pending: c.pendingURL,
			Failure: c.failureURL,
		},
	}

	resource, err := c.preferenceClient.Create(ctx, prefRequest)
	if err != nil {
		return nil, err
	}

	return &PreferenceResponse{
		ID:               resource.ID,
		InitPoint:        resource.InitPoint,
		SandboxInitPoint: resource.SandboxInitPoint,
	}, nil
}

// GetPayment fetches a payment and maps essential fields for the domain layer.
func (c *Client) GetPayment(ctx context.Context, id string) (*PaymentDetails, error) {
	if c == nil || !c.Enabled() {
		return nil, fmt.Errorf("mercado pago client is not configured")
	}

	paymentID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid payment id: %w", err)
	}

	resource, err := c.paymentClient.Get(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	meta := PaymentMetadata{}
	if resource.Metadata != nil {
		meta.CompanyID = stringFromAny(resource.Metadata["company_id"])
		meta.SchemaName = stringFromAny(resource.Metadata["schema_name"])
		meta.Months = intFromAny(resource.Metadata["months"])
	}

	var approved *time.Time
	if !resource.DateApproved.IsZero() {
		t := resource.DateApproved
		approved = &t
	}

	return &PaymentDetails{
		ID:                strconv.Itoa(resource.ID),
		Status:            resource.Status,
		CurrencyID:        resource.CurrencyID,
		TransactionAmount: resource.TransactionAmount,
		DateApproved:      approved,
		ExternalReference: resource.ExternalReference,
		Metadata:          meta,
	}, nil
}

func stringFromAny(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case float64:
		return strconv.FormatInt(int64(v), 10)
	case float32:
		return strconv.FormatInt(int64(v), 10)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case json.Number:
		return v.String()
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

func intFromAny(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case float32:
		return int(v)
	case string:
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	case json.Number:
		if parsed, err := v.Int64(); err == nil {
			return int(parsed)
		}
		if parsedFloat, err := v.Float64(); err == nil {
			return int(parsedFloat)
		}
	}
	return 0
}

// ValidateWebhookSignature verifies whether the webhook payload matches the Mercado Pago signature scheme.
func (c *Client) ValidateWebhookSignature(signatureHeader string, payload []byte) error {
	if c == nil || c.webhookSecret == "" || c.notificationURL == "" {
		return ErrWebhookSecretNotConfigured
	}

	ts, providedSignature, err := parseSignatureHeader(signatureHeader)
	if err != nil {
		return err
	}

	if err := validateSignatureTimestamp(ts); err != nil {
		return err
	}

	expected := computeSignature(c.webhookSecret, c.notificationURL, ts, payload)
	received, err := hex.DecodeString(providedSignature)
	if err != nil {
		return ErrSignatureMismatch
	}

	if !hmac.Equal(expected, received) {
		return ErrSignatureMismatch
	}

	return nil
}

func parseSignatureHeader(header string) (string, string, error) {
	header = strings.TrimSpace(header)
	if header == "" {
		return "", "", ErrSignatureHeaderMissing
	}

	var timestamp string
	var signature string

	segments := strings.FieldsFunc(header, func(r rune) bool {
		return r == ',' || r == ';'
	})

	for _, segment := range segments {
		segment = strings.TrimSpace(segment)
		if segment == "" {
			continue
		}
		kv := strings.SplitN(segment, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(kv[0]))
		value := strings.TrimSpace(kv[1])
		switch key {
		case "ts", "t", "timestamp":
			timestamp = value
		case "v1", "signature", "sha256":
			value = strings.TrimPrefix(value, "sha256=")
			value = strings.TrimPrefix(value, "sha256:")
			signature = strings.ToLower(strings.TrimSpace(value))
		}
	}

	if timestamp == "" || signature == "" {
		return "", "", ErrSignatureHeaderMissing
	}

	return timestamp, signature, nil
}

func validateSignatureTimestamp(raw string) error {
	parsed, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return ErrSignatureTimestampInvalid
	}

	ts := time.Unix(parsed, 0)
	now := time.Now()

	if ts.After(now.Add(signatureMaxSkewMinutes)) || ts.Before(now.Add(-signatureMaxSkewMinutes)) {
		return ErrSignatureTimestampInvalid
	}

	return nil
}

func computeSignature(secret, notificationURL, timestamp string, payload []byte) []byte {
	if payload == nil {
		payload = []byte{}
	}
	data := strings.Join([]string{timestamp, notificationURL, string(payload)}, ":")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return mac.Sum(nil)
}
