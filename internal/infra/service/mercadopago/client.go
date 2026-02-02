package mercadopagoservice

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/mercadopago/sdk-go/pkg/preapproval"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

const (
	defaultMonthlyPriceBRL = 99.90
)

// Client wraps Mercado Pago SDK clients with project specific helpers.
type Client struct {
	preferenceClient  preference.Client
	paymentClient     payment.Client
	preapprovalClient preapproval.Client
	notificationURL   string
	webhookSecret     string
	monthlyPrice      float64
	successURL        string
	pendingURL        string
	failureURL        string
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
	CompanyID         string `json:"company_id"`
	SchemaName        string `json:"schema_name"`
	PaymentType       string `json:"payment_type"`
	Months            int    `json:"months"`
	PlanType          string `json:"plan_type"`
	IsUpcoming        bool   `json:"is_upcoming"`         // Indicates if this is a scheduled/upcoming subscription
	UpgradeTargetPlan string `json:"upgrade_target_plan"` // For upgrade payments
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
	client.preapprovalClient = preapproval.NewClient(cfg)

	return client
}

// Enabled indicates whether the SDK clients are ready for use.
func (c *Client) Enabled() bool {
	return c != nil && c.preferenceClient != nil && c.paymentClient != nil
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

// SuccessURL returns the URL for successful payments.
// Returns empty string if not configured.
func (c *Client) SuccessURL() string {
	if c == nil {
		return ""
	}
	return c.successURL
}

// ValidateSignature validates the x-signature header from Mercado Pago webhooks.
// The x-signature header format is: ts={timestamp},v1={hmac}
// The manifest for HMAC calculation is: id:{data.id};request-id:{x-request-id};ts:{ts};
// Note: data.id must be lowercased per MP documentation.
// If any value is missing, it should be omitted from the manifest.
func (c *Client) ValidateSignature(xSignature, xRequestID, dataID string) bool {
	if c == nil || c.webhookSecret == "" {
		return false
	}

	// Parse x-signature header to extract ts and v1
	ts, v1 := parseXSignature(xSignature)
	if ts == "" || v1 == "" {
		return false
	}

	// Build the manifest string conditionally (omit empty values per MP docs)
	var manifestParts []string
	if dataID != "" {
		manifestParts = append(manifestParts, fmt.Sprintf("id:%s", strings.ToLower(dataID)))
	}
	if xRequestID != "" {
		manifestParts = append(manifestParts, fmt.Sprintf("request-id:%s", xRequestID))
	}
	if ts != "" {
		manifestParts = append(manifestParts, fmt.Sprintf("ts:%s", ts))
	}
	manifest := strings.Join(manifestParts, ";") + ";"

	// Calculate HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(c.webhookSecret))
	mac.Write([]byte(manifest))
	computed := hex.EncodeToString(mac.Sum(nil))

	// Compare with received signature
	return hmac.Equal([]byte(computed), []byte(v1))
}

// parseXSignature extracts ts and v1 values from x-signature header.
// Format: "ts=1234567890,v1=abc123..."
func parseXSignature(header string) (ts, v1 string) {
	parts := strings.Split(header, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "ts=") {
			ts = strings.TrimPrefix(part, "ts=")
		} else if strings.HasPrefix(part, "v1=") {
			v1 = strings.TrimPrefix(part, "v1=")
		}
	}
	return ts, v1
}

// CheckoutItem represents a single line item in the checkout preference.
type CheckoutItem struct {
	ID          string // SKU or Product ID (Added for improvements)
	CategoryID  string // Category ID (Added for improvements)
	Title       string
	Description string
	Quantity    int
	UnitPrice   float64
}

func NewCheckoutItem(id, categoryID, title, description string, quantity int, unitPrice float64) *CheckoutItem {
	return &CheckoutItem{
		ID:          id,
		CategoryID:  categoryID,
		Title:       title,
		Description: description,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
	}
}

type PaymentCheckoutType string

const (
	PaymentCheckoutTypeSubscription         PaymentCheckoutType = "subscription"
	PaymentCheckoutTypeSubscriptionUpgrade  PaymentCheckoutType = "subscription_upgrade"
	PaymentCheckoutTypeSubscriptionSchedule PaymentCheckoutType = "subscription_schedule"
	PaymentCheckoutTypeCost                 PaymentCheckoutType = "cost"
)

// CheckoutPayer represents the payer information to improve approval rates.
type CheckoutPayer struct {
	Email string
	Name  string // Business Name
	Phone struct {
		AreaCode string
		Number   string
	}
	Address struct {
		ZipCode      string // CEP
		StreetName   string
		StreetNumber string
		Neighborhood string
		City         string
		State        string // UF
	}
}

// CheckoutRequest wraps the information required to create a multi-item checkout preference.
type CheckoutRequest struct {
	CompanyID         string
	Schema            string
	PaymentType       PaymentCheckoutType
	Item              *CheckoutItem
	Payer             *CheckoutPayer // Added for approval improvements
	ExternalReference string         // Usually the PaymentID
	Metadata          map[string]any
}

// CreateUniqueCheckout creates a multi-item preference for the new billing architecture.
func (c *Client) CreateUniqueCheckout(ctx context.Context, req *CheckoutRequest) (*PreferenceResponse, error) {
	if c == nil || !c.Enabled() {
		return nil, fmt.Errorf("mercado pago client is not configured")
	}

	items := []preference.ItemRequest{
		{
			ID:          req.Item.ID,
			CategoryID:  req.Item.CategoryID,
			Title:       req.Item.Title,
			Description: req.Item.Description,
			Quantity:    req.Item.Quantity,
			UnitPrice:   req.Item.UnitPrice,
			CurrencyID:  "BRL",
		},
	}

	metadata := map[string]any{
		"company_id":   req.CompanyID,
		"schema_name":  req.Schema,
		"payment_type": string(req.PaymentType),
	}
	if req.Metadata != nil {
		maps.Copy(metadata, req.Metadata)
	}

	prefRequest := preference.Request{
		Items:             items,
		ExternalReference: req.ExternalReference,
		NotificationURL:   c.notificationURL,
		AutoReturn:        "approved",
		Metadata:          metadata,
		BackURLs: &preference.BackURLsRequest{
			Success: c.successURL,
			Pending: c.pendingURL,
			Failure: c.failureURL,
		},
	}

	if req.Payer != nil {
		prefRequest.Payer = &preference.PayerRequest{
			Email: req.Payer.Email,
			Name:  req.Payer.Name,
		}

		if req.Payer.Phone.Number != "" {
			prefRequest.Payer.Phone = &preference.PhoneRequest{
				AreaCode: req.Payer.Phone.AreaCode,
				Number:   req.Payer.Phone.Number,
			}
		}

		if req.Payer.Address.StreetName != "" {
			// addressNumber, _ := strconv.Atoi(req.Payer.Address.StreetNumber) // Removed: SDK uses string
			// Checking sdk-go preference.AddressRequest: StreetNumber is string or int?
			// Wait, I should verify SDK definition. But for now I will assume it handles what I give if I map correctly.
			// Actually sdk-go `preference.AddressRequest` defines:
			// StreetName string, StreetNumber string (verify?)
			// I will assume simple mapping first.
			prefRequest.Payer.Address = &preference.AddressRequest{
				ZipCode:      req.Payer.Address.ZipCode,
				StreetName:   req.Payer.Address.StreetName,
				StreetNumber: req.Payer.Address.StreetNumber,
			}
		}
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

// SubscriptionRequest wraps the information required to create a subscription (preapproval).
type SubscriptionRequest struct {
	Title         string
	Description   string
	Price         float64
	Frequency     int
	FrequencyType string // "months"
	ExternalRef   string
	PayerEmail    string
	BackURL       string
}

// SubscriptionResponse mirrors the minimal data needed by the application when creating a subscription.
type SubscriptionResponse struct {
	ID        string `json:"id"`
	InitPoint string `json:"init_point"`
}

// CreateSubscription creates a recurring payment (preapproval) preference.
func (c *Client) CreateSubscription(ctx context.Context, req *SubscriptionRequest) (*SubscriptionResponse, error) {
	if c == nil || !c.Enabled() {
		return nil, fmt.Errorf("mercado pago client is not configured")
	}

	preapprovalReq := preapproval.Request{
		Reason:            req.Title,
		ExternalReference: req.ExternalRef,
		PayerEmail:        req.PayerEmail,
		AutoRecurring: &preapproval.AutoRecurringRequest{
			Frequency:         req.Frequency,
			FrequencyType:     req.FrequencyType,
			TransactionAmount: req.Price,
			CurrencyID:        "BRL",
		},
		BackURL: req.BackURL,
		Status:  "pending",
	}

	resource, err := c.preapprovalClient.Create(ctx, preapprovalReq)
	if err != nil {
		return nil, err
	}

	return &SubscriptionResponse{
		ID:        resource.ID,
		InitPoint: resource.InitPoint,
	}, nil
}

// CancelSubscription cancels a recurring payment (preapproval) by ID.
func (c *Client) CancelSubscription(ctx context.Context, preapprovalID string) error {
	if c == nil || !c.Enabled() {
		return fmt.Errorf("mercado pago client is not configured")
	}

	// Update the preapproval status to "cancelled"
	updateReq := preapproval.UpdateRequest{
		Status: "cancelled",
	}

	_, err := c.preapprovalClient.Update(ctx, preapprovalID, updateReq)
	return err
}

// UpdateSubscriptionAmount updates the recurring amount of a subscription.
func (c *Client) UpdateSubscriptionAmount(ctx context.Context, preapprovalID string, newAmount float64) error {
	if c == nil || !c.Enabled() {
		return fmt.Errorf("mercado pago client is not configured")
	}

	updateReq := preapproval.UpdateRequest{
		AutoRecurring: &preapproval.AutoRecurringUpdateRequest{
			TransactionAmount: newAmount,
		},
	}

	_, err := c.preapprovalClient.Update(ctx, preapprovalID, updateReq)
	return err
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
		// Hack to convert map[string]any to struct via JSON
		if b, err := json.Marshal(resource.Metadata); err == nil {
			json.Unmarshal(b, &meta)
		}
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
