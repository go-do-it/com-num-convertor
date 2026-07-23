package converter

import (
	"testing"
)

func TestConvertKownValues(t *testing.T) {
	cases := []struct {
		name 		string
		input		string
		from		int
		to 			int
		expected	string
	} {

		{"hex to binary", "0x25B9D2", 0, Binary, "1001011011100111010010"},
		{"binary to hex", "1010111001001001", Binary, Hexadecimal, "AE49"},
		{"hex to binary 2", "0xA8B3D", 0, Binary, "10101000101100111101"},
		{"binary to hex 2", "1100100010110110010110", Binary, Hexadecimal, "322D96"},
		{"decimal to hex", "255", Decimal, Hexadecimal, "FF"},
		{"hex to decimal", "0xFF", 0, Decimal, "255"},
		{"negative decimal to binary", "-10", Decimal, Binary, "-1010"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := Convert(c.input, c.from, c.to)
			if err := nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.expected {
				t.Errorf("Convert(%q, %d, %d) =  %q, want %q" c.input, c.from, c.to, got, c.expected)
			}
		})
	}

}

func TestInvalidDigit(t *testing.T) {
	_, err := Convert("120", Binary, Decimal)
	if err == nil {
		t.Fatal("expected error for invalid digit, got nil")
	}
}

func TestInvalidBase(t *testing.T) {
	_, err := Convert("10", 1, Decimal)
	if err != ErrInvalidBase {
		t.Errorf("expected ErrInvalidBase, got %v", err)
	}
}

func TestAutoDetectPrefix(t *testing.T) {
	got, err := ToDecimal("0b1010", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "10" {
		t.Errorf("got %q, want %q", got, "10")
	}
}

