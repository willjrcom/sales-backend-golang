package companyentity_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	companyentity "github.com/willjrcom/sales-backend-go/internal/domain/company"
)

func TestPreferences_GetString(t *testing.T) {
	entries := []companyentity.Preference{
		{Key: companyentity.EnableDelivery, Value: "true"},
	}
	prefs := companyentity.NewPreferences(entries)
	v, err := prefs.GetString(companyentity.EnableDelivery)
	require.NoError(t, err)
	require.Equal(t, "true", v)

	_, err = prefs.GetString(companyentity.TableTaxRate)
	require.Error(t, err)
}

func TestPreferences_GetDecimal(t *testing.T) {
	entries := []companyentity.Preference{
		{Key: companyentity.TableTaxRate, Value: "0.15"},
	}
	prefs := companyentity.NewPreferences(entries)
	d, err := prefs.GetDecimal(companyentity.TableTaxRate)
	require.NoError(t, err)
	require.True(t, d.Equal(decimal.RequireFromString("0.15")))

	// invalid decimal
	prefs = companyentity.NewPreferences([]companyentity.Preference{{Key: companyentity.MinDeliveryTax, Value: "abc"}})
	_, err = prefs.GetDecimal(companyentity.MinDeliveryTax)
	require.Error(t, err)
}

func TestPreferences_GetBool(t *testing.T) {
	entries := []companyentity.Preference{
		{Key: companyentity.EnableTable, Value: "false"},
	}
	prefs := companyentity.NewPreferences(entries)
	b, err := prefs.GetBool(companyentity.EnableTable)
	require.NoError(t, err)
	require.False(t, b)

	// invalid bool
	prefs = companyentity.NewPreferences([]companyentity.Preference{{Key: companyentity.EnableTable, Value: "notabool"}})
	_, err = prefs.GetBool(companyentity.EnableTable)
	require.Error(t, err)
}
