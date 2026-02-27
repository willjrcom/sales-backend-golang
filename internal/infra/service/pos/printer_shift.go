package pos

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/shopspring/decimal"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	shiftentity "github.com/willjrcom/sales-backend-go/internal/domain/shift"
)

func FormatShift(shift *shiftentity.Shift) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(escInit)

	var raw bytes.Buffer
	raw.WriteString(escAlignCenter)
	raw.WriteString(escBoldOn)
	raw.WriteString("RELATORIO DE PEDIDOS")
	raw.WriteString(escBoldOff)
	raw.WriteString(newline)
	fmt.Fprintf(&raw, "Abertura: %s%s", shift.OpenedAt.Format("02/01/2006 15:04"), newline)

	raw.WriteString(escAlignLeft)
	raw.WriteString(strings.Repeat("-", 40) + newline)

	raw.WriteString(escAlignCenter)
	raw.WriteString(escBoldOn)

	raw.WriteString("PEDIDOS")
	raw.WriteString(escBoldOff)
	raw.WriteString(newline)

	raw.WriteString(escAlignLeft)
	raw.WriteString("PEDIDO\tVALOR")
	raw.WriteString(newline)

	totalVendas := decimal.NewFromFloat(0)
	for _, o := range shift.Orders {
		FormatOrderShift(&raw, &o)
		totalVendas.Add(o.SubTotal)
	}

	raw.WriteString(strings.Repeat("-", 40) + newline)
	raw.WriteString(escAlignCenter)
	raw.WriteString(escBoldOn)

	raw.WriteString("TOTAL DE VENDAS")
	raw.WriteString(escBoldOff)
	raw.WriteString(newline)

	raw.WriteString(escAlignLeft)
	raw.WriteString(fmt.Sprintf("TOTAL GERAL:\tR$ %7.2f", d2f(totalVendas)))
	tw := tabwriter.NewWriter(&buf, 6, 11, 2, ' ', 0)
	tw.Write(raw.Bytes())
	tw.Flush()

	buf.WriteString(escCut)
	return buf.Bytes(), nil
}

func FormatOrderShift(buf *bytes.Buffer, o *orderentity.Order) {
	buf.WriteString(fmt.Sprintf("%d\tR$ %7.2f%s", o.OrderNumber, d2f(o.SubTotal), newline))
}
