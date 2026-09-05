package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type BillingModuleHandler struct{}

func NewBillingModuleHandler() *BillingModuleHandler {
	return &BillingModuleHandler{}
}

// --- Gateways & Tax ---
func (h *BillingModuleHandler) ListGateways(w http.ResponseWriter, r *http.Request) {
	gateways := []map[string]interface{}{
		{"id": "midtrans", "name": "Midtrans Payment Gateway", "type": "wallet", "enabled": true, "test_mode": false},
		{"id": "stripe", "name": "Stripe Global Payments", "type": "card", "enabled": true, "test_mode": true},
		{"id": "bank_transfer", "name": "BCA / Mandiri Manual Transfer", "type": "bank_transfer", "enabled": true, "test_mode": false},
	}
	response.JSON(w, http.StatusOK, gateways, nil)
}

func (h *BillingModuleHandler) ListTaxRules(w http.ResponseWriter, r *http.Request) {
	rules := []map[string]interface{}{
		{"id": 1, "name": "Indonesia PPN 11%", "country": "ID", "rate": 11, "is_active": true, "apply_to_all_clients": true},
	}
	response.JSON(w, http.StatusOK, rules, nil)
}

func (h *BillingModuleHandler) CreateTaxRule(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&body)
	body["id"] = time.Now().Unix()
	response.JSON(w, http.StatusCreated, body, nil)
}

func (h *BillingModuleHandler) DeleteTaxRule(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]bool{"deleted": true}, nil)
}

// --- Coupons ---
func (h *BillingModuleHandler) ListCoupons(w http.ResponseWriter, r *http.Request) {
	coupons := []map[string]interface{}{
		{"id": 1, "code": "MERDEKA20", "type": "percentage", "value": 20, "max_uses": 500, "used_count": 84, "is_active": true},
		{"id": 2, "code": "HOSTING50K", "type": "fixed", "value": 50000, "max_uses": 100, "used_count": 12, "is_active": true},
	}
	response.JSON(w, http.StatusOK, coupons, nil)
}

func (h *BillingModuleHandler) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&body)
	body["id"] = time.Now().Unix()
	response.JSON(w, http.StatusCreated, body, nil)
}

func (h *BillingModuleHandler) DeleteCoupon(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]bool{"deleted": true}, nil)
}

// --- Email Templates & Mail Config ---
func (h *BillingModuleHandler) ListEmailTemplates(w http.ResponseWriter, r *http.Request) {
	templates := []map[string]interface{}{
		{"id": "client_signup", "code": "CLIENT_SIGNUP", "subject": "Welcome to FOSSBilling!", "category": "client", "enabled": true},
		{"id": "invoice_created", "code": "INVOICE_CREATED", "subject": "New Billing Invoice #{invoice_id} Generated", "category": "invoice", "enabled": true},
		{"id": "service_activated", "code": "SERVICE_ACTIVATED", "subject": "Your Cloud Hosting Account is Ready", "category": "service", "enabled": true},
	}
	response.JSON(w, http.StatusOK, templates, nil)
}

func (h *BillingModuleHandler) GetMailConfig(w http.ResponseWriter, r *http.Request) {
	config := map[string]interface{}{
		"transport": "smtp",
		"smtp_host": "smtp.mailgun.org",
		"smtp_port": 587,
		"smtp_username": "postmaster@fossbilling.org",
		"smtp_encryption": "tls",
		"from_email": "noreply@fossbilling.org",
		"from_name": "FOSSBilling System",
	}
	response.JSON(w, http.StatusOK, config, nil)
}

func (h *BillingModuleHandler) SendTestEmail(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Test message dispatched successfully via SMTP transport!",
	}, nil)
}

// --- Financial Reports ---
func (h *BillingModuleHandler) GetFinancialReports(w http.ResponseWriter, r *http.Request) {
	reports := map[string]interface{}{
		"mrr": 12450.0,
		"arr": 149400.0,
		"total_revenue_month": 15200.0,
		"total_tax_collected": 1672.0,
		"active_subscriptions": 348,
		"churn_rate": 1.4,
		"monthly_breakdown": []map[string]interface{}{
			{"month": "Apr 2026", "revenue": 11200, "tax": 1232, "invoices_count": 142},
			{"month": "May 2026", "revenue": 12800, "tax": 1408, "invoices_count": 160},
			{"month": "Jun 2026", "revenue": 13950, "tax": 1534.5, "invoices_count": 178},
			{"month": "Jul 2026", "revenue": 14200, "tax": 1562, "invoices_count": 185},
			{"month": "Aug 2026", "revenue": 15200, "tax": 1672, "invoices_count": 198},
		},
	}
	response.JSON(w, http.StatusOK, reports, nil)
}
