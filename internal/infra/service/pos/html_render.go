package pos

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/shopspring/decimal"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

// RenderGroupItemKitchenHTML returns the rendered HTML for a kitchen ticket
func RenderGroupItemKitchenHTML(group *orderentity.GroupItem) ([]byte, error) {
	tmpl, err := template.New("kitchen").Parse(KitchenTicketTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, group); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// RenderOrderHTML returns the rendered HTML for a full order receipt
func RenderOrderHTML(order *orderentity.Order) ([]byte, error) {
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
		"printf": fmt.Sprintf,
	}

	tmpl, err := template.New("receipt").Funcs(funcMap).Parse(OrderReceiptTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, order); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
