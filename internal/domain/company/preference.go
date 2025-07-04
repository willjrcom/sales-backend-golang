package companyentity

import (
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
)

// Key defines the enumeration of available preference keys.
type Key string

const (
	// TableTaxRate is the percentage rate applied for table consumption (e.g., 0.10 for 10%).
	TableTaxRate Key = "table_tax_rate"
	// MinOrderValueForFreeDelivery is the minimum order value to qualify for free delivery.
	MinOrderValueForFreeDelivery Key = "min_order_value_for_free_delivery"
	// EnableMinOrderValueForFreeDelivery toggles the free delivery minimum order value rule.
	EnableMinOrderValueForFreeDelivery Key = "enable_min_order_value_for_free_delivery"
	// EnableDelivery toggles delivery availability.
	EnableDelivery Key = "enable_delivery"
	// EnableTable toggles table service availability.
	EnableTable Key = "enable_table"
	// MinDeliveryTax is the minimum fee applied for delivery.
	MinDeliveryTax Key = "min_delivery_tax"

	// EnablePrintOrderOnShipOrder toggles printing when launching an order.
	EnablePrintOrderOnShipOrder Key = "enable_print_order_on_pend_order"
	// PrinterOrderOnPendOrder is the printer used for orders.
	PrinterOrderOnPendOrder Key = "printer_order_on_pend_order"
	// EnablePrintOrderOnShipDelivery toggles printing when launching a delivery.
	EnablePrintOrderOnShipDelivery Key = "enable_print_order_on_ship_delivery"
	// PrinterDeliveryOnShipDelivery is the printer used for deliveries.
	PrinterDeliveryOnShipDelivery Key = "printer_delivery_on_ship_delivery"
	// EnablePrintOrderOnShipDelivery toggles printing when finish a process.
	EnablePrintItemsOnFinishProcess Key = "enable_print_items_on_finish_process"
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

func NewDefaultPreferences() Preferences {
	return Preferences{
		EnableDelivery:                     "true",
		EnableTable:                        "true",
		EnableMinOrderValueForFreeDelivery: "true",
		TableTaxRate:                       "0.10",
		MinDeliveryTax:                     "0.00",
		MinOrderValueForFreeDelivery:       "0.00",
	}
}

// GetString returns the raw string value for the given key, or an error if missing.
func (p Preferences) GetString(key Key) (string, error) {
	v, ok := p[key]
	if !ok {
		return "", fmt.Errorf("preference %q not found", key)
	}

	return v, nil
}

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
