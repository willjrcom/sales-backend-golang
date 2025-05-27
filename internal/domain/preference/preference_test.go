package preference_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/willjrcom/sales-backend-go/internal/domain/preference"
)

func TestPreferences_GetString(t *testing.T) {
	entries := []preference.Preference{
		{Key: preference.EnableDelivery, Value: "true"},
	}
	prefs := preference.NewPreferences(entries)
	v, err := prefs.GetString(preference.EnableDelivery)
	require.NoError(t, err)
	require.Equal(t, "true", v)

	_, err = prefs.GetString(preference.ConsumptionTaxRate)
	require.Error(t, err)
}

func TestPreferences_GetDecimal(t *testing.T) {
	entries := []preference.Preference{
		{Key: preference.ConsumptionTaxRate, Value: "0.15"},
	}
	prefs := preference.NewPreferences(entries)
	d, err := prefs.GetDecimal(preference.ConsumptionTaxRate)
	require.NoError(t, err)
	require.True(t, d.Equal(decimal.RequireFromString("0.15")))

	// invalid decimal
	prefs = preference.NewPreferences([]preference.Preference{{Key: preference.MinDeliveryFee, Value: "abc"}})
	_, err = prefs.GetDecimal(preference.MinDeliveryFee)
	require.Error(t, err)
}

func TestPreferences_GetBool(t *testing.T) {
	entries := []preference.Preference{
		{Key: preference.EnableTables, Value: "false"},
	}
	prefs := preference.NewPreferences(entries)
	b, err := prefs.GetBool(preference.EnableTables)
	require.NoError(t, err)
	require.False(t, b)

	// invalid bool
	prefs = preference.NewPreferences([]preference.Preference{{Key: preference.EnableTables, Value: "notabool"}})
	_, err = prefs.GetBool(preference.EnableTables)
	require.Error(t, err)
}

func TestPreferences_MustHelpers(t *testing.T) {
	entries := []preference.Preference{
		{Key: preference.MinOrderValueForFreeDelivery, Value: "100.00"},
		{Key: preference.EnableMinOrderValueForFreeDelivery, Value: "true"},
	}
	prefs := preference.NewPreferences(entries)
	d := prefs.MustDecimal(preference.MinOrderValueForFreeDelivery)
	require.True(t, d.Equal(decimal.RequireFromString("100.00")))
	b := prefs.MustBool(preference.EnableMinOrderValueForFreeDelivery)
	require.True(t, b)
}
