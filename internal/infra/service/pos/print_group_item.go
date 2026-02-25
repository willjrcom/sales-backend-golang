package pos

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	companydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/company"
)

// FormatGroupItemKitchen generates ESC/POS bytes for a kitchen print of a group of items,
// showing only item names, and complements, without prices or totals.
func FormatGroupItemKitchen(group *orderentity.GroupItem, company *companydto.CompanyDTO) ([]byte, error) {
	var final bytes.Buffer
	// Initialize printer and select Latin-1 code page
	final.WriteString(escInit)
	final.WriteString(escCodePageLatin1)

	// --- CABEÇALHO COZINHA (Centralizado e Negrito) ---
	var headerRaw bytes.Buffer
	printGroupItemKitchenHeader(&headerRaw, company)
	final.WriteString(escAlignCenter)
	final.WriteString(escBoldOn)
	final.Write(ToLatin1(headerRaw.String()))
	final.WriteString(escBoldOff)
	final.WriteString(escAlignLeft)

	// --- CORPO COZINHA (Tabelado) ---
	var bodyRaw bytes.Buffer
	printGroupItemKitchen(&bodyRaw, group)
	bodyRaw.WriteString(strings.Repeat(newline, 2))

	var bodyAligned bytes.Buffer
	tw := tabwriter.NewWriter(&bodyAligned, 0, 0, 1, ' ', 0)
	tw.Write(bodyRaw.Bytes())
	tw.Flush()

	final.Write(ToLatin1(bodyAligned.String()))

	// cut ticket
	final.WriteString(escCut)
	return final.Bytes(), nil
}

// printGroupItemKitchen writes group header, items, additions, removals, and group complement,
// ignoring price values for kitchen production.
func printGroupItemKitchen(buf *bytes.Buffer, group *orderentity.GroupItem) {
	categoryName := ""
	if c := group.Category; c != nil && c.Name != "" {
		categoryName = c.Name
	}
	categoryLabel := truncate(categoryName, 15)
	sizeLabel := truncate(group.Size, 10)
	quantityStr := fmt.Sprintf("Qtd:%.0f", group.Quantity)

	buf.WriteString(fmt.Sprintf("%-15s\t%-10s\t%s%s", categoryLabel, sizeLabel, quantityStr, newline))

	// Items and their observations
	for _, item := range group.Items {
		name := truncate(item.Name, 20)
		buf.WriteString(fmt.Sprintf("%-4.1f\t%s%s", item.Quantity, name, newline))
		for _, add := range item.AdditionalItems {
			addName := truncate(add.Name, 17)
			buf.WriteString(fmt.Sprintf("+\t%.0fx %s%s", add.Quantity, addName, newline))
		}
		for _, rm := range item.RemovedItems {
			buf.WriteString(fmt.Sprintf("-\t%s%s", rm, newline))
		}
		if obs := item.Observation; obs != "" {
			buf.WriteString(fmt.Sprintf("OBS:\t%s%s", obs, newline))
		}
	}

	// Group complement item
	if comp := group.ComplementItem; comp != nil {
		compName := truncate(comp.Name, 20)
		buf.WriteString(fmt.Sprintf("COMPLEMENTO:\t%s%s", compName, newline))
	}
}

func printGroupItemKitchenHeader(buf *bytes.Buffer, company *companydto.CompanyDTO) {
	if company != nil {
		fmt.Fprintf(buf, "%s%s", company.TradeName, newline)
	}
	buf.WriteString("COZINHA" + newline)
	buf.WriteString(newline)
}
