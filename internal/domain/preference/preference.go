package preference

import (
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
)

// Key defines the enumeration of available preference keys.
type Key string

const (
	// ConsumptionTaxRate is the percentage rate applied for table consumption (e.g., 0.10 for 10%).
	ConsumptionTaxRate Key = "consumption_tax_rate"
	// MinOrderValueForFreeDelivery is the minimum order value to qualify for free delivery.
	MinOrderValueForFreeDelivery Key = "min_order_value_for_free_delivery"
	// EnableMinOrderValueForFreeDelivery toggles the free delivery minimum order value rule.
	EnableMinOrderValueForFreeDelivery Key = "enable_min_order_value_for_free_delivery"
	// EnableDelivery toggles delivery availability.
	EnableDelivery Key = "enable_delivery"
	// EnableTables toggles table service availability.
	EnableTables Key = "enable_tables"
	// MinDeliveryFee is the minimum fee applied for delivery.
	MinDeliveryFee Key = "min_delivery_fee"
)

// Preference holds a single key-value pair.
type Preference struct {
	Key   Key
	Value string
}

// Preferences is a map of preference keys to their raw string values.
type Preferences map[Key]string

// NewPreferences builds a map from a slice of Preference entries.
func NewPreferences(entries []Preference) Preferences {
	prefs := make(Preferences, len(entries))
	for _, e := range entries {
		prefs[e.Key] = e.Value
	}
	return prefs
}

// GetString returns the raw string value for the given key, or an error if missing.
func (p Preferences) GetString(key Key) (string, error) {
	v, ok := p[key]
	if !ok {
		return "", fmt.Errorf("preference %q not found", key)
	}
	return v, nil
}

// GetDecimal parses the value for key as decimal.Decimal.
func (p Preferences) GetDecimal(key Key) (decimal.Decimal, error) {
	raw, err := p.GetString(key)
	if err != nil {
		return decimal.Zero, err
	}
	dec, err := decimal.NewFromString(raw)
	if err != nil {
		return decimal.Zero, fmt.Errorf("invalid decimal for %q: %w", key, err)
	}
	return dec, nil
}

// GetBool parses the value for key as bool.
func (p Preferences) GetBool(key Key) (bool, error) {
	raw, err := p.GetString(key)
	if err != nil {
		return false, err
	}
	b, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("invalid bool for %q: %w", key, err)
	}
	return b, nil
}

// MustDecimal returns the parsed decimal or panics if error.
func (p Preferences) MustDecimal(key Key) decimal.Decimal {
	d, err := p.GetDecimal(key)
	if err != nil {
		panic(err)
	}
	return d
}

// MustBool returns the parsed bool or panics if error.
func (p Preferences) MustBool(key Key) bool {
	b, err := p.GetBool(key)
	if err != nil {
		panic(err)
	}
	return b
}
