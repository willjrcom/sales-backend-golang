package pos

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/shopspring/decimal"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
)

// RenderShiftHTML returns the rendered HTML for a shift report
func RenderShiftHTML(shift *shiftentity.Shift) ([]byte, error) {
	funcMap := template.FuncMap{
		"formatDate": func(t *time.Time) string {
			if t == nil {
				return ""
			}
			return t.Format("02/01/2006 15:04")
		},
		"formatMoney": func(d decimal.Decimal) string {
			return "R$ " + d.StringFixed(2)
		},
		"now": func() string {
			return time.Now().Format("02/01/2006 15:04:05")
		},
	}

	tmpl, err := template.New("shift").Funcs(funcMap).Parse(ShiftReportTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, shift); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// RenderGroupItemKitchenHTML returns the rendered HTML for a kitchen ticket
func RenderGroupItemKitchenHTML(group *orderentity.GroupItem, company *companydto.CompanyDTO) ([]byte, error) {
	tmpl, err := template.New("kitchen").Parse(KitchenTicketTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	data := struct {
		*orderentity.GroupItem
		Company *companydto.CompanyDTO
	}{
		GroupItem: group,
		Company:   company,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// RenderOrderHTML returns the rendered HTML for a full order receipt
func RenderOrderHTML(order *orderentity.Order, company *companydto.CompanyDTO) ([]byte, error) {
	// Helper function for date formatting
	funcMap := template.FuncMap{
		"formatDate": func(t *time.Time) string {
			if t == nil {
				return ""
			}
			return t.Format("02/01/2006 15:04")
		},
		"formatMoney": func(d decimal.Decimal) string {
			return "R$ " + d.StringFixed(2)
		},
		"multiply": func(d decimal.Decimal, q float64) decimal.Decimal {
			return d.Mul(decimal.NewFromFloat(q))
		},
		"printf": fmt.Sprintf,
	}

	tmpl, err := template.New("receipt").Funcs(funcMap).Parse(OrderReceiptTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	data := struct {
		*orderentity.Order
		Company *companydto.CompanyDTO
	}{
		Order:   order,
		Company: company,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
