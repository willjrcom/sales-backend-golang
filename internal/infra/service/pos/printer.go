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

	// Print section based on order type: Delivery, Pickup, or Table
	switch {
	case o.Delivery != nil:
		formatDeliverySection(&buf, o)
	case o.Pickup != nil:
		formatPickupSection(&buf, o)
	case o.Table != nil:
		formatTableSection(&buf, o)
	}

	formatOrderDetailSection(&buf, o)
	printGroupItemsSection(&buf, o.GroupItems)
	formatPaymentsSection(&buf, o)
	formatFooter(&buf, o)
	buf.WriteString(strings.Repeat(newline, 3))
	buf.WriteString(escCut)
	return buf.Bytes(), nil
}

func formatHeader(buf *bytes.Buffer, o *orderentity.Order) {
	if o.PendingAt != nil {
		buf.WriteString(fmt.Sprintf("Gerado: %s%s", o.PendingAt.Format("02/01/2006 15:04"), newline))
	}

	if o.Attendant != nil && o.Attendant.User != nil {
		buf.WriteString(fmt.Sprintf("Atendente: %s%s", o.Attendant.User.Name, newline))
	}
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

	buf.WriteString(fmt.Sprintf("%-40s%s", header, newline))
	buf.WriteString(escBoldOff)

	// Quantity of group
	buf.WriteString(fmt.Sprintf("Qtd:%36.0f%s", group.Quantity, newline))

	// Group time logs
	if group.StartAt != nil {
		buf.WriteString(fmt.Sprintf("Agendado: %s%s", group.StartAt.Format("02/01/2006 15:04"), newline))
	}
	// if group.PendingAt != nil {
	// 	buf.WriteString(fmt.Sprintf("Pendente: %s%s", group.PendingAt.Format("02/01/2006 15:04"), newline))
	// }
	// if group.StartedAt != nil {
	// 	buf.WriteString(fmt.Sprintf("Iniciado: %s%s", group.StartedAt.Format("02/01/2006 15:04"), newline))
	// }
	if group.ReadyAt != nil {
		buf.WriteString(fmt.Sprintf("Pronto: %s%s", group.ReadyAt.Format("02/01/2006 15:04"), newline))
	}
	// if group.CanceledAt != nil {
	// 	buf.WriteString(fmt.Sprintf("Cancelado: %s%s", group.CanceledAt.Format("02/01/2006 15:04"), newline))
	// }

	// Observation
	if obs := group.Observation; obs != "" {
		buf.WriteString(escBoldOn)
		buf.WriteString(fmt.Sprintf("Obs: %s%s", obs, newline))
		buf.WriteString(escBoldOff)
	}

	// Items
	for _, item := range group.Items {
		printItem(buf, &item)
	}

	// Complement item
	if comp := group.ComplementItem; comp != nil {
		printComplementItem(buf, comp, group)
	}

	// Subtotal for this group
	buf.WriteString(escBoldOn)
	buf.WriteString(fmt.Sprintf("Subtotal:%31.2f%s", d2f(group.TotalPrice), newline))
	buf.WriteString(escBoldOff)
}

func printComplementItem(buf *bytes.Buffer, comp *orderentity.Item, group *orderentity.GroupItem) {
	buf.WriteString(escBoldOn)
	buf.WriteString(fmt.Sprintf("%4.1f x %-20s %7.2f%s", group.Quantity, comp.Name, d2f(comp.TotalPrice), newline))
	buf.WriteString(escBoldOff)
}

// printItem writes a single order item and its additional items to the buffer.
func printItem(buf *bytes.Buffer, item *orderentity.Item) {
	name := item.Name
	if len(name) > 20 {
		name = name[:20]
	}

	buf.WriteString(fmt.Sprintf("%4.1f x %-20s %7.2f%s", item.Quantity, name, d2f(item.TotalPrice), newline))
	for _, add := range item.AdditionalItems {
		printAdditionalItem(buf, &add)
	}

	// Removed items for item
	if len(item.RemovedItems) > 0 {
		for _, rm := range item.RemovedItems {
			buf.WriteString(fmt.Sprintf("   - %s%s", rm, newline))
		}
	}

	// Observation for item
	if obs := item.Observation; obs != "" {
		buf.WriteString(escBoldOn)
		buf.WriteString(fmt.Sprintf("   Obs: %s%s", obs, newline))
		buf.WriteString(escBoldOff)
	}

	buf.WriteString(newline)
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
	buf.WriteString(escAlignLeft)

	// Delivery time logs
	// if pa := o.Delivery.PendingAt; pa != nil {
	// 	buf.WriteString(fmt.Sprintf("Pendente: %s%s", pa.Format("02/01/2006 15:04"), newline))
	// }
	if sa := o.Delivery.ShippedAt; sa != nil {
		buf.WriteString(fmt.Sprintf("Despachado: %s%s", sa.Format("02/01/2006 15:04"), newline))
	}
	// Client name
	if o.Delivery.Client != nil && o.Delivery.Client.Name != "" {
		buf.WriteString(fmt.Sprintf("Cliente: %s%s", o.Delivery.Client.Name, newline))
	}
	// Address
	if a := o.Delivery.Address; a != nil {
		buf.WriteString(fmt.Sprintf("Endereço: %s, %s%s", a.Street, a.Number, newline))

		if a.Complement != "" {
			buf.WriteString(fmt.Sprintf("Complemento: %s%s", a.Complement, newline))
		}

		if a.Reference != "" {
			buf.WriteString(fmt.Sprintf("Ref: %s%s", a.Reference, newline))
		}

		buf.WriteString(fmt.Sprintf("Bairro: %s%s", a.Neighborhood, newline))
		buf.WriteString(fmt.Sprintf("Cidade: %s - %s%s", a.City, a.UF, newline))
		buf.WriteString(fmt.Sprintf("CEP: %s%s", a.Cep, newline))
	}

	// Delivery driver
	if d := o.Delivery.Driver; d != nil && d.Employee != nil && d.Employee.User != nil {
		buf.WriteString(fmt.Sprintf("Motoboy: %s%s", d.Employee.User.Name, newline))
	}

	// Delivery tax
	if t := o.Delivery.DeliveryTax; t != nil {
		buf.WriteString(fmt.Sprintf("Taxa entrega: %7.2f%s", d2f(*t), newline))
	}

	// Change for delivery
	buf.WriteString(fmt.Sprintf("Troco entrega: %7.2f%s", d2f(o.Delivery.Change), newline))

	// Payment method for delivery
	buf.WriteString(fmt.Sprintf("Pagamento entrega: %s%s", o.Delivery.PaymentMethod, newline))
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
	// 	buf.WriteString(fmt.Sprintf("Pendente: %s%s", pa.Format("02/01/2006 15:04"), newline))
	// }
	if ra := o.Pickup.ReadyAt; ra != nil {
		buf.WriteString(fmt.Sprintf("Pronto: %s%s", ra.Format("02/01/2006 15:04"), newline))
	}
	if name := o.Pickup.Name; name != "" {
		buf.WriteString(fmt.Sprintf("Cliente: %s%s", name, newline))
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
	// 	buf.WriteString(fmt.Sprintf("Pendente: %s%s", pa.Format("02/01/2006 15:04"), newline))
	// }
	// if ca := o.Table.ClosedAt; ca != nil {
	// 	buf.WriteString(fmt.Sprintf("Fechado: %s%s", ca.Format("02/01/2006 15:04"), newline))
	// }
	if name := o.Table.Name; name != "" {
		buf.WriteString(fmt.Sprintf("Nome: %s%s", name, newline))
	}
	if contact := o.Table.Contact; contact != "" {
		buf.WriteString(fmt.Sprintf("Contato: %s%s", contact, newline))
	}
	buf.WriteString(strings.Repeat("-", 40) + newline)
}

// formatOrderDetailSection prints order detail fields: observation, items count, paid and change totals.
func formatOrderDetailSection(buf *bytes.Buffer, o *orderentity.Order) {
	buf.WriteString(escAlignLeft)
	// Observation
	if obs := o.Observation; obs != "" {
		buf.WriteString(escBoldOn)
		buf.WriteString(fmt.Sprintf("Observação: %s%s", obs, newline))
		buf.WriteString(escBoldOff)
	}

	// Total items
	if o.QuantityItems > 0 {
		buf.WriteString(fmt.Sprintf("Itens:%36.0f%s", o.QuantityItems, newline))
	}

	// Total paid
	buf.WriteString(fmt.Sprintf("Pago:%31.2f%s", d2f(o.TotalPaid), newline))

	// Total change
	buf.WriteString(fmt.Sprintf("Troco:%31.2f%s", d2f(o.TotalChange), newline))
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
