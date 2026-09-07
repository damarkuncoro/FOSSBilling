package period_test

import (
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/period"
)

func TestPeriod_ParseAndCalculate(t *testing.T) {
	tests := []struct {
		code         string
		expectedQty  int
		expectedUnit string
		shouldError  bool
	}{
		{"1M", 1, "M", false},
		{"3M", 3, "M", false},
		{"1Y", 1, "Y", false},
		{"2W", 2, "W", false},
		{"30D", 30, "D", false},
		{"99Y", 0, "", true},  // exceeds max 5 years
		{"100D", 0, "", true}, // exceeds max 90 days
		{"invalid", 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			p, err := period.Parse(tt.code)
			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for code '%s', got nil", tt.code)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error for code '%s': %v", tt.code, err)
			}

			if p.Qty != tt.expectedQty || p.Unit != tt.expectedUnit {
				t.Errorf("Parsed period = %+v; want qty=%d, unit=%s", p, tt.expectedQty, tt.expectedUnit)
			}
		})
	}

	// Test Expiration Calculation
	baseDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	pMonth, _ := period.Parse("1M")
	expMonth := pMonth.CalculateExpiration(baseDate)
	if expMonth.Month() != time.February || expMonth.Day() != 1 {
		t.Errorf("Expected 2026-02-01, got %v", expMonth)
	}

	pYear, _ := period.Parse("1Y")
	expYear := pYear.CalculateExpiration(baseDate)
	if expYear.Year() != 2027 {
		t.Errorf("Expected year 2027, got %v", expYear)
	}
}
