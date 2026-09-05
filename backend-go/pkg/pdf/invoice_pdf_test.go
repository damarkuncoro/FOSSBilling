package pdf_test

import (
	"strings"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/pdf"
)

func TestGenerateInvoiceHTML(t *testing.T) {
	client := &domain.Client{
		FirstName: "Ahmad",
		LastName:  "Dhani",
		Email:     "ahmad.dhani@example.com",
		Company:   "Republik Cinta Management",
		Country:   "ID",
	}

	inv := &domain.Invoice{
		ID:       555,
		Nr:       "INV-00555",
		Status:   domain.InvoiceStatusPaid,
		Currency: "USD",
		Subtotal: 1000000,
		Tax:      110000,
		Total:    1110000,
		DueAt:    time.Now().AddDate(0, 0, 14),
		Items: []domain.InvoiceItem{
			{
				Title:    "Dedicated Server Bare Metal",
				Price:    1000000,
				Quantity: 1,
			},
		},
	}

	htmlBytes, err := pdf.GenerateInvoiceHTML(inv, client, "Nusantara Cloud", "", "")
	if err != nil {
		t.Fatalf("GenerateInvoiceHTML failed: %v", err)
	}

	htmlStr := string(htmlBytes)
	if !strings.Contains(htmlStr, "INV-00555") {
		t.Errorf("expected invoice number in output")
	}
	if !strings.Contains(htmlStr, "LUNAS / PAID") {
		t.Errorf("expected paid status stamp")
	}
	if !strings.Contains(htmlStr, "Dedicated Server Bare Metal") {
		t.Errorf("expected item title in table")
	}
	if !strings.Contains(htmlStr, "Ahmad Dhani") {
		t.Errorf("expected client name")
	}
}
