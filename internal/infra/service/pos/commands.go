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
	newline = "\r\n"
	// escCodePageLatin1 selects code page 16 (Latin-1) for accented characters
	escCodePageLatin1 = "\x1bt\x10"
)

// ToLatin1 converts a UTF-8 string to ISO-8859-1 (Latin-1) bytes.
// This is necessary because most thermal printers use Latin-1 for accented characters.
func ToLatin1(s string) []byte {
	res := make([]byte, 0, len(s))
	for _, r := range s {
		if r < 128 {
			res = append(res, byte(r))
		} else if r >= 160 && r <= 255 {
			res = append(res, byte(r))
		} else {
			// Map some common Portuguese characters that might be outside the 160-255 range if they exist
			// although most are within 160-255 in Unicode/Latin1
			switch r {
			case 'ç':
				res = append(res, 0xE7)
			case 'Ç':
				res = append(res, 0xC7)
			case 'ã':
				res = append(res, 0xE3)
			case 'Ã':
				res = append(res, 0xC3)
			case 'õ':
				res = append(res, 0xF5)
			case 'Õ':
				res = append(res, 0xD5)
			default:
				res = append(res, '?')
			}
		}
	}
	return res
}
