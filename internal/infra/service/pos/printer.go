package pos

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

// d2f converts decimal.Decimal to float64 for formatting.
func d2f(d decimal.Decimal) float64 {
	return d.InexactFloat64()
}

const (
	escInit        = "\x1b@"     // Initialize printer
	escBoldOn      = "\x1bE\x01" // Bold on
	escBoldOff     = "\x1bE\x00" // Bold off
	escAlignLeft   = "\x1ba\x00" // Align left
	escAlignCenter = "\x1ba\x01" // Align center
	escCut         = "\x1dV\x00" // Full cut
	newline        = "\n"
)

// FormatOrder generates ESC/POS bytes for a 40-column receipt of the given order.
// It initializes the printer, prints the header, item groups, footer, and cuts the paper.
func FormatOrder(o *orderentity.Order) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(escInit)
	formatHeader(&buf, o)
	printGroupItemsSection(&buf, o.GroupItems)
	formatFooter(&buf, o)
	buf.WriteString(newline + newline + newline)
	buf.WriteString(escCut)
	return buf.Bytes(), nil
}

func formatHeader(buf *bytes.Buffer, o *orderentity.Order) {
	buf.WriteString(escAlignCenter)
	buf.WriteString(escBoldOn)
	buf.WriteString(fmt.Sprintf("PEDIDO %d%s", o.OrderNumber, newline))
	buf.WriteString(escBoldOff)
	buf.WriteString(escAlignLeft)
	if o.PendingAt != nil {
		buf.WriteString(fmt.Sprintf("Gerado: %s%s", o.PendingAt.Format("02/01/2006 15:04"), newline))
	}

	if o.Attendant != nil && o.Attendant.User != nil {
		buf.WriteString(fmt.Sprintf("Atendente: %s%s", o.Attendant.User.Name, newline))
	}

	buf.WriteString(strings.Repeat("-", 40) + newline)
}

// printGroupItemsSection prints all item groups of the order.
func printGroupItemsSection(buf *bytes.Buffer, groups []orderentity.GroupItem) {
	for _, grp := range groups {
		printGroupItem(buf, &grp)
	}
	buf.WriteString(strings.Repeat("-", 40) + newline)
}

// printGroupItem prints group header, its items, any complement, and subtotal.
func printGroupItem(buf *bytes.Buffer, group *orderentity.GroupItem) {
	// Header: category and size
	buf.WriteString(escAlignLeft)
	buf.WriteString(escBoldOn)
	var parts []string
	if c := group.Category; c != nil && c.Name != "" {
		parts = append(parts, c.Name)
	}
	if s := group.Size; s != "" {
		parts = append(parts, s)
	}
	header := strings.Join(parts, " - ")
	if header == "" {
		header = "Grupo"
	}
	buf.WriteString(fmt.Sprintf("%-40s%s", header, newline))
	buf.WriteString(escBoldOff)
	// Quantity of group
	buf.WriteString(fmt.Sprintf("Qtd:%36.0f%s", group.Quantity, newline))
	// Observation
	if obs := group.Observation; obs != "" {
		buf.WriteString(fmt.Sprintf("Obs: %s%s", obs, newline))
	}
	// Items
	for _, item := range group.Items {
		printItem(buf, &item)
	}
	// Complement item
	if comp := group.ComplementItem; comp != nil {
		buf.WriteString(escBoldOn)
		buf.WriteString(fmt.Sprintf("%2.0f x %-20s %7.2f%s", group.Quantity, comp.Name, d2f(comp.TotalPrice), newline))
		buf.WriteString(escBoldOff)
	}
	// Subtotal for this group
	buf.WriteString(fmt.Sprintf("Subtotal:%31.2f%s", d2f(group.TotalPrice), newline))
}

// printItem writes a single order item and its additional items to the buffer.
func printItem(buf *bytes.Buffer, item *orderentity.Item) {
	name := item.Name
	if len(name) > 20 {
		name = name[:20]
	}
	buf.WriteString(fmt.Sprintf("%2.0f x %-20s %7.2f%s", item.Quantity, name, d2f(item.TotalPrice), newline))
	for _, add := range item.AdditionalItems {
		printAdditionalItem(buf, &add)
	}
	// Observation for item
	if obs := item.Observation; obs != "" {
		buf.WriteString(fmt.Sprintf("   Obs: %s%s", obs, newline))
	}
	// Removed items for item
	if len(item.RemovedItems) > 0 {
		for _, rm := range item.RemovedItems {
			buf.WriteString(fmt.Sprintf("   Removido: %s%s", rm, newline))
		}
	}
}

// printAdditionalItem writes a single additional item to the buffer.
func printAdditionalItem(buf *bytes.Buffer, add *orderentity.Item) {
	name := add.Name
	if len(name) > 17 {
		name = name[:17]
	}
	buf.WriteString(fmt.Sprintf("   + %-17s %7.2f%s", name, d2f(add.TotalPrice), newline))
}

// formatFooter writes the total payable amount to the buffer.
func formatFooter(buf *bytes.Buffer, o *orderentity.Order) {
	buf.WriteString(fmt.Sprintf("TOTAL:%31.2f%s", d2f(o.TotalPayable), newline))
}

// FormatGroupItem generates ESC/POS bytes for a 40-column receipt section for a single group of items.
// It prints each item of the group followed by a separator line.
func FormatGroupItem(group *orderentity.GroupItem) ([]byte, error) {
	var buf bytes.Buffer
	printGroupItem(&buf, group)
	buf.WriteString(strings.Repeat("-", 40) + newline)
	return buf.Bytes(), nil
}

// FormatItem generates ESC/POS bytes for a single order item, including its additional items.
func FormatItem(item *orderentity.Item) ([]byte, error) {
	var buf bytes.Buffer
	printItem(&buf, item)
	return buf.Bytes(), nil
}
