package decimal_test

import (
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

func TestParsePeriod(t *testing.T) {
	tests := []struct {
		code      string
		wantUnit  string
		wantQty   int
		wantTitle string
		wantErr   bool
	}{
		{"1W", "W", 1, "Weekly", false},
		{"2W", "W", 2, "Every 2 Weeks", false},
		{"1M", "M", 1, "Monthly", false},
		{"3M", "M", 3, "Quarterly", false},
		{"6M", "M", 6, "Semi-Annually", false},
		{"1Y", "Y", 1, "Annually", false},
		{"2Y", "Y", 2, "Biennially", false},
		{"3Y", "Y", 3, "Triennially", false},
		{"FREE", "F", 0, "Free", false},
		{"ONETIME", "O", 1, "One Time", false},
		{"invalid", "", 0, "", true},
		{"0M", "", 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			p, err := decimal.ParsePeriod(tt.code)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParsePeriod(%s) expected error, got nil", tt.code)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParsePeriod(%s) unexpected error: %v", tt.code, err)
			}
			if p.Unit != tt.wantUnit || p.Qty != tt.wantQty || p.Title != tt.wantTitle {
				t.Errorf("ParsePeriod(%s) = {%s, %d, %s}; want {%s, %d, %s}",
					tt.code, p.Unit, p.Qty, p.Title, tt.wantUnit, tt.wantQty, tt.wantTitle)
			}
		})
	}
}

func TestPeriod_CalculateNextDueDate(t *testing.T) {
	baseDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	p1M, _ := decimal.ParsePeriod("1M")
	next1M := p1M.CalculateNextDueDate(baseDate)
	if expected := time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC); !next1M.Equal(expected) {
		t.Errorf("1M next due date = %v; want %v", next1M, expected)
	}

	p1Y, _ := decimal.ParsePeriod("1Y")
	next1Y := p1Y.CalculateNextDueDate(baseDate)
	if expected := time.Date(2027, 1, 15, 0, 0, 0, 0, time.UTC); !next1Y.Equal(expected) {
		t.Errorf("1Y next due date = %v; want %v", next1Y, expected)
	}
}
