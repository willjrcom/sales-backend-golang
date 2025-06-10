package pos

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

// truncate returns a string truncated to at most max runes, preserving UTF-8 boundaries
func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) > max {
		return string(runes[:max])
	}
	return s
}

// FormatOrder generates ESC/POS bytes for a 40-column receipt of the given order.
// It initializes the printer, selects Latin-1 code page, prints the header, item groups, footer, and cuts the paper.
func FormatOrder(o *orderentity.Order) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(escInit)
	buf.WriteString(escCodePageLatin1)

	// Build raw text for receipt
	var raw bytes.Buffer
	switch {
	case o.Delivery != nil:
		formatDeliverySection(&raw, o)
	case o.Pickup != nil:
		formatPickupSection(&raw, o)
	case o.Table != nil:
		formatTableSection(&raw, o)
	}
	formatOrderDetailSection(&raw, o)
	printGroupItemsSection(&raw, o.GroupItems)
	formatPaymentsSection(&raw, o)
	formatTotalFooter(&raw, o)
	raw.WriteString(strings.Repeat(newline, 3))

	// Align columns using tabwriter into main buffer, checking for errors
	tw := tabwriter.NewWriter(&buf, 6, 11, 2, ' ', 0)
	if _, err := tw.Write(raw.Bytes()); err != nil {
		return nil, err
	}
	if err := tw.Flush(); err != nil {
		return nil, err
	}

	buf.WriteString(escCut)
	return buf.Bytes(), nil
}

func formatHeader(buf *bytes.Buffer, o *orderentity.Order) {
	buf.WriteString(escAlignLeft)
	if o.PendingAt != nil {
		buf.WriteString(fmt.Sprintf("Gerado:\t\t%s%s", o.PendingAt.Format("15:04"), newline))
	}

	if o.Attendant != nil && o.Attendant.User != nil {
		buf.WriteString(fmt.Sprintf("Atendente:\t\t%s%s", o.Attendant.User.Name, newline))
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

	parts = append(parts, "Qtd:"+fmt.Sprintf("%.0f", group.Quantity))
	header := strings.Join(parts, "\t")

	buf.WriteString(fmt.Sprintf("%-40s%s", header, newline))
	buf.WriteString(escBoldOff)

	// Group time logs
	if group.StartAt != nil {
		buf.WriteString(fmt.Sprintf("Agendado:\t%s%s", group.StartAt.Format("15:04"), newline))
	}
	// if group.PendingAt != nil {
	// 	buf.WriteString(fmt.Sprintf("Pendente:\t\t%s%s", group.PendingAt.Format("15:04"), newline))
	// }
	// if group.StartedAt != nil {
	// 	buf.WriteString(fmt.Sprintf("Iniciado:\t\t%s%s", group.StartedAt.Format("15:04"), newline))
	// }
	// if group.CanceledAt != nil {
	// 	buf.WriteString(fmt.Sprintf("Cancelado:\t\t%s%s", group.CanceledAt.Format("15:04"), newline))
	// }

	// Items: quantidade | produto | valor
	for _, item := range group.Items {
		printItem(buf, &item)
	}

	// Complement item
	if comp := group.ComplementItem; comp != nil {
		printComplementItem(buf, comp, group)
	}

	// Observation
	if obs := group.Observation; obs != "" {
		buf.WriteString(escBoldOn)
		buf.WriteString(fmt.Sprintf("Obs: %s%s", obs, newline))
		buf.WriteString(escBoldOff)
	}

	if group.ReadyAt != nil {
		buf.WriteString(fmt.Sprintf("Pronto as:\t\t%s%s", group.ReadyAt.Format("15:04"), newline))
	}

	// Subtotal for this group
	buf.WriteString(escBoldOn)
	buf.WriteString(fmt.Sprintf("Subtotal:\t\t%.2f%s", d2f(group.TotalPrice), newline))
	buf.WriteString(escBoldOff)
}

func printComplementItem(buf *bytes.Buffer, comp *orderentity.Item, group *orderentity.GroupItem) {
	buf.WriteString(escBoldOn)
	// truncate complement name to 20 runes to avoid breaking UTF-8
	name := truncate(comp.Name, 20)
	buf.WriteString(fmt.Sprintf("%4.1f\t%-20s\t%7.2f%s", group.Quantity, name, d2f(comp.TotalPrice), newline))
	buf.WriteString(escBoldOff)
}

// printItem writes a single order item and its additional items to the buffer.
// printItem writes a single order item and its additional items to the buffer.
func printItem(buf *bytes.Buffer, item *orderentity.Item) {
	// truncate item name to 20 runes to avoid breaking UTF-8
	name := truncate(item.Name, 20)
	buf.WriteString(fmt.Sprintf("%.1f\t%-20s\t%.2f%s", item.Quantity, name, d2f(item.TotalPrice), newline))

	for _, add := range item.AdditionalItems {
		printAdditionalItem(buf, &add)
	}

	// Removed items for item
	if len(item.RemovedItems) > 0 {
		for _, rm := range item.RemovedItems {
			buf.WriteString(fmt.Sprintf("-\t%s\t%s", rm, newline))
		}
	}

	// Observation for item
	if obs := item.Observation; obs != "" {
		buf.WriteString(escBoldOn)
		buf.WriteString(fmt.Sprintf("Obs: %s%s", obs, newline))
		buf.WriteString(escBoldOff)
	}

	buf.WriteString(newline)
}

// printAdditionalItem writes a single additional item to the buffer.
func printAdditionalItem(buf *bytes.Buffer, add *orderentity.Item) {
	// truncate additional item name to 17 runes to avoid breaking UTF-8
	name := truncate(add.Name, 17)
	buf.WriteString(fmt.Sprintf("+\t%-17s\t%.2f%s", name, d2f(add.TotalPrice), newline))
}

// formatTotalFooter writes the total payable amount to the buffer.
func formatTotalFooter(buf *bytes.Buffer, o *orderentity.Order) {
	buf.WriteString(fmt.Sprintf("TOTAL:\t\t%.2f%s", d2f(o.TotalPayable), newline))
}

// formatDeliverySection prints delivery-related details if present.
func formatDeliverySection(buf *bytes.Buffer, o *orderentity.Order) {
	if o.Delivery == nil {
		return
	}

	buf.WriteString(escAlignCenter)
	buf.WriteString(escBoldOn)
	buf.WriteString(fmt.Sprintf("PEDIDO DE ENTREGA %d%s", o.OrderNumber, newline))
	buf.WriteString(escBoldOff)

	formatHeader(buf, o)

	buf.WriteString(escAlignCenter)
	buf.WriteString("ENTREGA" + newline)
	buf.WriteString(escAlignLeft)

	// Delivery time logs
	// if pa := o.Delivery.PendingAt; pa != nil {
	// 	buf.WriteString(fmt.Sprintf("Pendente: %s%s", pa.Format("15:04"), newline))
	// }

	if sa := o.Delivery.ShippedAt; sa != nil {
		buf.WriteString(fmt.Sprintf("Despachado:\t\t%s%s", sa.Format("15:04"), newline))
	}

	// Client name
	if o.Delivery.Client != nil && o.Delivery.Client.Name != "" {
		// truncate client name to 20 runes to avoid breaking UTF-8
		clientName := truncate(o.Delivery.Client.Name, 20)
		buf.WriteString(fmt.Sprintf("Cliente:\t\t%-20s%s", clientName, newline))
	}

	// Address
	if a := o.Delivery.Address; a != nil {
		buf.WriteString(escAlignCenter)
		buf.WriteString("Endereço")
		buf.WriteString(newline)

		buf.WriteString(escAlignLeft)
		buf.WriteString(fmt.Sprintf("%s, %s%s", a.Street, a.Number, newline))

		if a.Complement != "" {
			buf.WriteString(fmt.Sprintf("Complemento:\t%s%s", a.Complement, newline))
		}

		if a.Reference != "" {
			buf.WriteString(fmt.Sprintf("Ref:\t%s%s", a.Reference, newline))
		}

		buf.WriteString(fmt.Sprintf("Bairro:\t%s%s", a.Neighborhood, newline))
		buf.WriteString(fmt.Sprintf("Cidade:\t%s - %s%s", a.City, a.UF, newline))
		buf.WriteString(fmt.Sprintf("CEP:\t%s%s", a.Cep, newline))
	}

	// Delivery driver
	if d := o.Delivery.Driver; d != nil && d.Employee != nil && d.Employee.User != nil {
		buf.WriteString(fmt.Sprintf("Entregador:\t%s%s", d.Employee.User.Name, newline))
	}

	// Delivery tax
	if t := o.Delivery.DeliveryTax; t != nil {
		buf.WriteString(fmt.Sprintf("Taxa entrega:\t\t%.2f%s", d2f(*t), newline))
	}

	// Change for delivery
	buf.WriteString(fmt.Sprintf("Troco:\t\t%.2f%s", d2f(o.Delivery.Change), newline))

	// Payment method for delivery
	buf.WriteString(fmt.Sprintf("Forma de pagamento:\t%s%s", o.Delivery.PaymentMethod, newline))
	buf.WriteString(strings.Repeat("-", 40) + newline)
}

// formatPickupSection prints pickup-related details if present.
func formatPickupSection(buf *bytes.Buffer, o *orderentity.Order) {
	if o.Pickup == nil {
		return
	}

	buf.WriteString(escAlignCenter)
	buf.WriteString(escBoldOn)
	buf.WriteString(fmt.Sprintf("PEDIDO DE RETIRADA %d%s", o.OrderNumber, newline))
	buf.WriteString(escBoldOff)

	formatHeader(buf, o)
	buf.WriteString(escAlignLeft)

	// if pa := o.Pickup.PendingAt; pa != nil {
	// 	buf.WriteString(fmt.Sprintf("Pendente: %s%s", pa.Format("15:04"), newline))
	// }
	if ra := o.Pickup.ReadyAt; ra != nil {
		buf.WriteString(fmt.Sprintf("Pronto:\t\t%s%s", ra.Format("15:04"), newline))
	}
	if name := o.Pickup.Name; name != "" {
		buf.WriteString(fmt.Sprintf("Cliente:\t\t%s%s", name, newline))
	}
	buf.WriteString(strings.Repeat("-", 40) + newline)
}

// formatTableSection prints table-related details if present.
func formatTableSection(buf *bytes.Buffer, o *orderentity.Order) {
	if o.Table == nil {
		return
	}

	buf.WriteString(escAlignCenter)
	buf.WriteString(escBoldOn)
	buf.WriteString(fmt.Sprintf("PEDIDO DE MESA %d%s", o.OrderNumber, newline))
	buf.WriteString(escBoldOff)

	formatHeader(buf, o)
	buf.WriteString(escAlignLeft)

	// if pa := o.Table.PendingAt; pa != nil {
	// 	buf.WriteString(fmt.Sprintf("Pendente: %s%s", pa.Format("15:04"), newline))
	// }
	// if ca := o.Table.ClosedAt; ca != nil {
	// 	buf.WriteString(fmt.Sprintf("Fechado: %s%s", ca.Format("15:04"), newline))
	// }
	if name := o.Table.Name; name != "" {
		buf.WriteString(fmt.Sprintf("Nome:\t%s%s", name, newline))
	}
	if contact := o.Table.Contact; contact != "" {
		buf.WriteString(fmt.Sprintf("Contato:\t%s%s", contact, newline))
	}
	buf.WriteString(strings.Repeat("-", 40) + newline)
}

// formatOrderDetailSection prints order detail fields: observation, items count, paid and change totals.
func formatOrderDetailSection(buf *bytes.Buffer, o *orderentity.Order) {
	buf.WriteString(escAlignLeft)
	// Observation
	if obs := o.Observation; obs != "" {
		buf.WriteString(escBoldOn)
		buf.WriteString("Observação do pedido")
		buf.WriteString(newline)
		buf.WriteString(obs)
		buf.WriteString(escBoldOff)
		buf.WriteString(newline)
	}

	// Total items
	buf.WriteString(fmt.Sprintf("Total de itens:\t\t%.0f%s", o.QuantityItems, newline))

	// Total paid
	buf.WriteString(fmt.Sprintf("Pago:\t\t%.2f%s", d2f(o.TotalPaid), newline))

	// Total change
	buf.WriteString(fmt.Sprintf("Troco:\t\t%.2f%s", d2f(o.TotalChange), newline))
	buf.WriteString(strings.Repeat("-", 40) + newline)
}

// formatPaymentsSection writes each payment entry of the order.
func formatPaymentsSection(buf *bytes.Buffer, o *orderentity.Order) {
	if len(o.Payments) == 0 {
		return
	}
	buf.WriteString(escAlignLeft)
	for _, p := range o.Payments {
		buf.WriteString(fmt.Sprintf("%s:%7.2f%s", p.Method, d2f(p.TotalPaid), newline))
	}
	buf.WriteString(strings.Repeat("-", 40) + newline)
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
