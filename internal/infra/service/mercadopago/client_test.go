package mercadopagoservice

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

func TestValidateSignature(t *testing.T) {
	// Test secret key
	secret := "test-secret-key-12345"

	// Create a client with the secret
	client := &Client{
		webhookSecret: secret,
	}

	// Test data
	dataID := "12345678901"
	requestID := "abc123-request-id"
	ts := "1704908010"

	// Build expected manifest (same as MP docs)
	// Template: id:[data.id];request-id:[x-request-id];ts:[ts];
	expectedManifest := fmt.Sprintf("id:%s;request-id:%s;ts:%s;", strings.ToLower(dataID), requestID, ts)
	t.Logf("Expected manifest: %s", expectedManifest)

	// Calculate expected HMAC
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(expectedManifest))
	expectedV1 := hex.EncodeToString(mac.Sum(nil))
	t.Logf("Expected v1 signature: %s", expectedV1)

	// Build x-signature header
	xSignature := fmt.Sprintf("ts=%s,v1=%s", ts, expectedV1)
	t.Logf("x-signature header: %s", xSignature)

	// Test validation
	valid := client.ValidateSignature(xSignature, requestID, dataID)
	if !valid {
		t.Errorf("Expected signature to be valid, but got invalid")
	}

	// Test with wrong signature
	wrongSignature := "ts=1704908010,v1=wrongsignature123"
	invalidResult := client.ValidateSignature(wrongSignature, requestID, dataID)
	if invalidResult {
		t.Errorf("Expected wrong signature to be invalid, but got valid")
	}

	// Test with empty secret
	clientNoSecret := &Client{}
	noSecretResult := clientNoSecret.ValidateSignature(xSignature, requestID, dataID)
	if noSecretResult {
		t.Errorf("Expected validation to fail with no secret, but got valid")
	}

	// Test with uppercase dataID (should be lowercased)
	uppercaseDataID := "ABC123DEF456"
	uppercaseManifest := fmt.Sprintf("id:%s;request-id:%s;ts:%s;", strings.ToLower(uppercaseDataID), requestID, ts)
	uppercaseMac := hmac.New(sha256.New, []byte(secret))
	uppercaseMac.Write([]byte(uppercaseManifest))
	uppercaseV1 := hex.EncodeToString(uppercaseMac.Sum(nil))
	uppercaseSignature := fmt.Sprintf("ts=%s,v1=%s", ts, uppercaseV1)

	uppercaseResult := client.ValidateSignature(uppercaseSignature, requestID, uppercaseDataID)
	if !uppercaseResult {
		t.Errorf("Expected uppercase dataID to be valid (lowercased), but got invalid")
	}

	t.Log("All signature validation tests passed!")
}

func TestParseXSignature(t *testing.T) {
	// Test normal case
	ts, v1 := parseXSignature("ts=1704908010,v1=618c85345248dd820d5fd456117c2ab2ef8eda45a0282ff693eac24131a5e839")
	if ts != "1704908010" {
		t.Errorf("Expected ts=1704908010, got ts=%s", ts)
	}
	if v1 != "618c85345248dd820d5fd456117c2ab2ef8eda45a0282ff693eac24131a5e839" {
		t.Errorf("Unexpected v1 value: %s", v1)
	}

	// Test with spaces
	ts2, v12 := parseXSignature("ts=123, v1=abc")
	if ts2 != "123" || v12 != "abc" {
		t.Errorf("Expected ts=123, v1=abc, got ts=%s, v1=%s", ts2, v12)
	}

	// Test empty
	ts3, v13 := parseXSignature("")
	if ts3 != "" || v13 != "" {
		t.Errorf("Expected empty values for empty header")
	}

	t.Log("All parseXSignature tests passed!")
}
