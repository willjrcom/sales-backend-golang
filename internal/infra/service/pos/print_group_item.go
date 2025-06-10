package pos

import (
	"bytes"
	"fmt"
	"strings"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

// FormatGroupItemKitchen generates ESC/POS bytes for a kitchen print of a group of items,
// showing only item names, quantities, and complements, without prices or totals.
func FormatGroupItemKitchen(group *orderentity.GroupItem) ([]byte, error) {
	var buf bytes.Buffer
	// Initialize printer and select Latin-1 code page
	buf.WriteString(escInit)
	buf.WriteString(escCodePageLatin1)

	// print group items for kitchen
	printGroupItemKitchen(&buf, group)

	// extra line feeds after group
	buf.WriteString(strings.Repeat(newline, 2))

	// cut ticket
	buf.WriteString(escCut)
	return buf.Bytes(), nil
}

// printGroupItemKitchen writes group header, items, additions, removals, and group complement,
// ignoring price values for kitchen production.
func printGroupItemKitchen(buf *bytes.Buffer, group *orderentity.GroupItem) {
	buf.WriteString(escAlignLeft)
	// Header: category, size, quantity
	var parts []string
	if c := group.Category; c != nil && c.Name != "" {
		parts = append(parts, c.Name)
	}
	if s := group.Size; s != "" {
		parts = append(parts, s)
	}
	parts = append(parts, fmt.Sprintf("Qtd:%.0f", group.Quantity))
	buf.WriteString(strings.Join(parts, " ") + newline)

	// Items and their observations
	for _, item := range group.Items {
		name := truncate(item.Name, 20)
		buf.WriteString(fmt.Sprintf("%.1f\t%s%s", item.Quantity, name, newline))
		for _, add := range item.AdditionalItems {
			addName := truncate(add.Name, 17)
			buf.WriteString(fmt.Sprintf("+\t%s%s", addName, newline))
		}
		for _, rm := range item.RemovedItems {
			buf.WriteString(fmt.Sprintf("-\t%s%s", rm, newline))
		}
		if obs := item.Observation; obs != "" {
			buf.WriteString(fmt.Sprintf("Obs: %s%s", obs, newline))
		}
	}

	// Group complement item
	if comp := group.ComplementItem; comp != nil {
		compName := truncate(comp.Name, 20)
		buf.WriteString(fmt.Sprintf("Complemento: %s%s", compName, newline))
	}
}
