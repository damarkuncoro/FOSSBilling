package pdf_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/pdf"
)

func TestGenerateInvoicePDF(t *testing.T) {
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

	pdfBytes, err := pdf.GenerateInvoicePDF(inv, client, "Nusantara Cloud", "", "")
	if err != nil {
		t.Fatalf("GenerateInvoicePDF failed: %v", err)
	}

	// 1. Verify standard PDF 1.4 header
	if !bytes.HasPrefix(pdfBytes, []byte("%PDF-1.4")) {
		t.Errorf("expected PDF header %%PDF-1.4, got %s", string(pdfBytes[:8]))
	}

	// 2. Verify standard EOF marker
	if !bytes.HasSuffix(bytes.TrimSpace(pdfBytes), []byte("%%EOF")) {
		t.Errorf("expected PDF document to end with %%%%EOF")
	}

	// 3. Verify content contains invoice details
	pdfStr := string(pdfBytes)
	if !strings.Contains(pdfStr, "INV-00555") {
		t.Errorf("expected invoice number in PDF output")
	}
	if !strings.Contains(pdfStr, "(PAID)") {
		t.Errorf("expected paid status badge in PDF output")
	}
	if !strings.Contains(pdfStr, "Dedicated Server Bare Metal") {
		t.Errorf("expected item title in PDF output")
	}
	if !strings.Contains(pdfStr, "Ahmad Dhani") {
		t.Errorf("expected client name in PDF output")
	}
}
