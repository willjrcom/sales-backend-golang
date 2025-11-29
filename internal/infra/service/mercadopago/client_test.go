package mercadopagoservice

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestValidateWebhookSignatureSuccess(t *testing.T) {
	client := &Client{
		notificationURL: "https://example.com/webhook",
		webhookSecret:   "super-secret",
	}

	payload := []byte(`{"data":{"id":"123"}}`)
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	signature := hex.EncodeToString(computeSignature(client.webhookSecret, client.notificationURL, ts, payload))
	header := fmt.Sprintf("ts=%s,v1=%s", ts, signature)

	if err := client.ValidateWebhookSignature(header, payload); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestValidateWebhookSignatureMismatch(t *testing.T) {
	client := &Client{
		notificationURL: "https://example.com/webhook",
		webhookSecret:   "super-secret",
	}

	payload := []byte(`{"data":{"id":"123"}}`)
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	header := fmt.Sprintf("ts=%s,v1=%s", ts, strings.Repeat("0", 64))

	err := client.ValidateWebhookSignature(header, payload)
	if !errors.Is(err, ErrSignatureMismatch) {
		t.Fatalf("expected ErrSignatureMismatch, got %v", err)
	}
}

func TestValidateWebhookSignatureExpired(t *testing.T) {
	client := &Client{
		notificationURL: "https://example.com/webhook",
		webhookSecret:   "super-secret",
	}

	payload := []byte(`{"data":{"id":"123"}}`)
	ts := strconv.FormatInt(time.Now().Add(-10*time.Minute).Unix(), 10)
	signature := hex.EncodeToString(computeSignature(client.webhookSecret, client.notificationURL, ts, payload))
	header := fmt.Sprintf("ts=%s,v1=%s", ts, signature)

	err := client.ValidateWebhookSignature(header, payload)
	if !errors.Is(err, ErrSignatureTimestampInvalid) {
		t.Fatalf("expected ErrSignatureTimestampInvalid, got %v", err)
	}
}
