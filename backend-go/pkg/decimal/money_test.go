package decimal

import (
	"testing"
)

func TestMoney_Conversions(t *testing.T) {
	tests := []struct {
		name     string
		floatVal float64
		expected Money
		strVal   string
	}{
		{"Zero", 0.0, 0, "0.00"},
		{"Standard price", 19.99, 199900, "19.99"},
		{"Fractional cent rounding", 10.0055, 100055, "10.01"},
		{"Large financial amount", 1250000.50, 12500005000, "1250000.50"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := FromFloat(tt.floatVal)
			if m != tt.expected {
				t.Errorf("FromFloat(%v) = %v; want %v", tt.floatVal, m, tt.expected)
			}

			if m.String() != tt.strVal {
				t.Errorf("Money.String() = %v; want %v", m.String(), tt.strVal)
			}

			parsed, err := FromString(tt.strVal)
			if err != nil {
				t.Fatalf("FromString(%s) returned unexpected error: %v", tt.strVal, err)
			}
			if parsed.String() != tt.strVal {
				t.Errorf("FromString(%s).String() = %s; want %s", tt.strVal, parsed.String(), tt.strVal)
			}
		})
	}
}

func TestMoney_FormatPrecision(t *testing.T) {
	m := FromFloat(12.3456)
	if got := m.FormatPrecision(4); got != "12.3456" {
		t.Errorf("FormatPrecision(4) = %v; want 12.3456", got)
	}
	if got := m.FormatPrecision(2); got != "12.35" {
		t.Errorf("FormatPrecision(2) = %v; want 12.35", got)
	}
}
