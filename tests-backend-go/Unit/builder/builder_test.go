package builder_test

import (
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/builder"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
)

func TestInvoiceBuilder_Calculation(t *testing.T) {
	promo := &domain.Promo{
		ID:     1,
		Code:   "DISCOUNT10",
		Type:   domain.PromoTypePercentage,
		Value:  100000, // 10%
		Active: true,
	}

	invoice, items, err := builder.NewInvoiceBuilder().
		ForClient(10).
		WithSerieAndNr("INV", "2026-9999").
		WithCurrency("USD", 1.0).
		WithDueDays(14).
		WithTaxRate(10.0). // 10% tax
		ApplyPromo(promo).
		AddItem("cPanel Hosting Pro", decimal.FromFloat(100.0), 1, true).
		AddItem("Domain Registration .com", decimal.FromFloat(15.0), 1, false).
		Build()

	if err != nil {
		t.Fatalf("InvoiceBuilder failed: %v", err)
	}

	if invoice.ClientID != 10 {
		t.Errorf("expected client ID 10, got %d", invoice.ClientID)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
	if invoice.Subtotal != decimal.FromFloat(115.0) {
		t.Errorf("expected subtotal 115.00, got %s", invoice.Subtotal.String())
	}
	if invoice.Total <= 0 {
		t.Errorf("expected total > 0, got %s", invoice.Total.String())
	}
}

func TestMessageBuilder(t *testing.T) {
	msg, err := mailer.NewMessageBuilder().
		From("FOSSBilling Support", "support@fossbilling.org").
		To("customer@example.com").
		Subject("Pemberitahuan Layanan").
		HTMLBody("<p>Layanan Anda aktif</p>").
		AttachPDF("invoice.pdf", []byte("%PDF-1.4 mock pdf")).
		Build()

	if err != nil {
		t.Fatalf("MessageBuilder failed: %v", err)
	}

	if msg.Subject != "Pemberitahuan Layanan" {
		t.Errorf("expected subject 'Pemberitahuan Layanan', got '%s'", msg.Subject)
	}
	if len(msg.Attachments) != 1 {
		t.Errorf("expected 1 attachment, got %d", len(msg.Attachments))
	}
}

func TestOrderBuilder(t *testing.T) {
	order, err := builder.NewOrderBuilder().
		ForClient(1).
		ForProduct(101, "Cloud VPS Dedicated 8GB").
		WithPeriod("1M").
		WithPrice(decimal.FromFloat(45.0), "USD").
		WithConfig("hostname", "vps.solusinusantara.com").
		AsActive().
		Build()

	if err != nil {
		t.Fatalf("OrderBuilder failed: %v", err)
	}

	if order.Status != domain.OrderStatusActive {
		t.Errorf("expected status active, got %s", order.Status)
	}
	if order.ActivatedAt == nil {
		t.Errorf("expected activatedAt to be set")
	}
}

func TestClientBuilder(t *testing.T) {
	client, err := builder.NewClientBuilder().
		WithEmail("client.builder@example.com").
		WithPassword("SecurePassword999!").
		WithName("Eko", "Patrio").
		WithCompany("PT Sukses Makmur").
		WithCountryAndCurrency("ID", "IDR").
		Build()

	if err != nil {
		t.Fatalf("ClientBuilder failed: %v", err)
	}

	if client.Email != "client.builder@example.com" {
		t.Errorf("expected email client.builder@example.com, got %s", client.Email)
	}
	if client.PasswordHash == "" {
		t.Errorf("expected passwordHash to be set")
	}
}

