package pos

import "github.com/shopspring/decimal"

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
	// newline is carriage return + line feed for ESC/POS
	newline        = "\r\n"
	// escCodePageLatin1 selects code page 16 (Latin-1) for accented characters
	escCodePageLatin1 = "\x1bt\x10"
)
