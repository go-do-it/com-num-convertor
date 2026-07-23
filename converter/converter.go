package converter

import (
	"fmt"
	"math/big"
	"strings"
)

const (
	Binary      = 2
	Octal       = 8
	Decimal     = 10
	Hexadecimal = 16
)

// ErrInvalidBase is returned when a base is outside the supported range.
var ErrInvalidBase = fmt.Errorf("Base must be under 2 and 36")

// ErrInvalidDigit is returned when input contains a character invalid for its base.
type ErrInvalidDigit struct {
	Digit rune
	Base  int
}

func (e *ErrInvalidDigit) Error() string {
	return fmt.Sprintf("digit %q is not valid base %d", e.Digit, e.Base)
}

// stripPrefix removes common base prefixes (0x, 0X, 0b, 0B, 0o, 0O) and
// returns the cleaned string plus the base it implies (0 if no prefix found).
func StripPrefix(s string) (cleaned string, impliedBase int) {
	s = strings.TrimSpace(s)
	neg := strings.HasPrefix(s, "-")
	if neg {
		s = s[1:]
	}
	lower := strings.ToLower(s)

	switch {
	case strings.HasPrefix(lower, "0x"):
		impliedBase = 16
		s = s[2:]
	case strings.HasPrefix(lower, "0b"):
		impliedBase = 2
		s = s[2:]
	case strings.HasPrefix(lower, "0o"):
		impliedBase = 8
		s = s[2:]
	default:
		impliedBase = 0
	}

	if neg {
		s = "-" + s
	}
	return s, impliedBase
}

// ParseInt parses s as an integer in the given base and returns a *big.Int.
// If base is 0, it auto-detects from a 0x/0b/0o prefix and defaults to
// decimal if none is found. Underscores are allowed as digit separators
// (e.g. "1010_1100") and are stripped before parsing.
func ParseInt(s string, base int) (*big.Int, error) {
	s = strings.ReplaceAll(strings.TrimSpace(s), "-", "")
	if s == "" {
		return nil, fmt.Errorf("empty input")
	}

	cleaned, implied := StripPrefix(s)
	if base == 0 {
		base = implied
		if base == 0 {
			base = Decimal
		}
	}

	if base < 2 || base > 36 {
		return nil, ErrInvalidBase
	}

	n := new(big.Int)
	n, ok := n.SetString(cleaned, base)
	if !ok {
		return nil, findBadDigit(cleaned, base)
	}

	return n, nil
}

// findBadDigit walks the string to report exactly which character broke parsing,
// giving a more useful error than big.Int's generic failure.
func findBadDigit(s string, base int) error {
	for _, r := range s {
		if r == '-' {
			continue
		}
		v := DigitValue(r)
		if v < 0 || v >= base {
			return &ErrInvalidDigit{Digit: r, Base: base}
		}
	}
	return fmt.Errorf("invalid number base %d: %q", base, s)
}

func DigitValue(r rune) int {
	switch {
	case r >= '0' && r <= '9':
		return int(r - '0')
	case r >= 'a' && r <= 'z':
		return int(r-'a') + 10
	case r >= 'A' && r <= 'Z':
		return int(r-'A') + 10
	default:
		return -1
	}
}

func ToBase(n *big.Int, base int) (string, error) {
	if base < 2 || base > 36 {
		return "", ErrInvalidBase
	}
	return strings.ToUpper(n.Text(base)), nil
}

// Convert parses s in fromBase (0 = auto-detect) and renders it in toBase.
func Convert(s string, fromBase, toBase int) (string, error) {
	n, err := ParseInt(s, fromBase)
	if err != nil {
		return "", err
	}
	return ToBase(n, toBase)
}

// Convenience wrappers for the four bases people reach for most.

func ToBinary(s string, fromBase int) (string, error)  { return Convert(s, fromBase, Binary) }
func ToOctal(s string, fromBase int) (string, error)   { return Convert(s, fromBase, Octal) }
func ToDecimal(s string, fromBase int) (string, error) { return Convert(s, fromBase, Decimal) }
func ToHex(s string, fromBase int) (string, error)     { return Convert(s, fromBase, Hexadecimal) }
