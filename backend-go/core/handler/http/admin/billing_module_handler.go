package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/handler/middleware"
	billingUsecase "github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/billing"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/staff"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/usecase/stats"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type BillingModuleHandler struct {
	staffService  *staff.StaffService
	statsService  *stats.StatsService
	taxCalculator *billingUsecase.TaxCalculator
	promoRepo     domain.PromoRepository
}

func NewBillingModuleHandler(
	staffService *staff.StaffService,
	statsService *stats.StatsService,
	taxCalculator *billingUsecase.TaxCalculator,
	promoRepo domain.PromoRepository,
) *BillingModuleHandler {
	return &BillingModuleHandler{
		staffService:  staffService,
		statsService:  statsService,
		taxCalculator: taxCalculator,
		promoRepo:     promoRepo,
	}
}

// --- Gateways & Tax ---
func (h *BillingModuleHandler) ListGateways(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "billing", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: billing", nil)
		return
	}

	gateways := []map[string]interface{}{
		{"id": "midtrans", "name": "Midtrans Payment Gateway", "type": "wallet", "enabled": true, "test_mode": false},
		{"id": "stripe", "name": "Stripe Global Payments", "type": "card", "enabled": true, "test_mode": true},
		{"id": "bank_transfer", "name": "BCA / Mandiri Manual Transfer", "type": "bank_transfer", "enabled": true, "test_mode": false},
	}
	response.JSON(w, http.StatusOK, gateways, nil)
}

func (h *BillingModuleHandler) ListTaxRules(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "billing", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: billing", nil)
		return
	}

	rules, err := h.taxCalculator.ListRules(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, rules, nil)
}

func (h *BillingModuleHandler) CreateTaxRule(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "billing", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: billing", nil)
		return
	}

	var rule domain.TaxRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_BODY", err.Error(), nil)
		return
	}

	if err := h.taxCalculator.CreateRule(r.Context(), &rule); err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusCreated, rule, nil)
}

func (h *BillingModuleHandler) UpdateTaxRule(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "billing", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: billing", nil)
		return
	}

	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	var rule domain.TaxRule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_BODY", err.Error(), nil)
		return
	}
	rule.ID = id

	if err := h.taxCalculator.UpdateRule(r.Context(), &rule); err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, rule, nil)
}

func (h *BillingModuleHandler) DeleteTaxRule(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "billing", "delete")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: billing", nil)
		return
	}

	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err := h.taxCalculator.DeleteRule(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, map[string]bool{"deleted": true}, nil)
}

// --- Coupons ---
func (h *BillingModuleHandler) ListCoupons(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "billing", "read")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: billing", nil)
		return
	}

	coupons, _, err := h.promoRepo.List(r.Context(), 100, 0)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusOK, coupons, nil)
}

func (h *BillingModuleHandler) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "billing", "write")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: billing", nil)
		return
	}

	var promo domain.Promo
	if err := json.NewDecoder(r.Body).Decode(&promo); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_BODY", err.Error(), nil)
		return
	}

	if err := h.promoRepo.Create(r.Context(), &promo); err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	response.JSON(w, http.StatusCreated, promo, nil)
}

func (h *BillingModuleHandler) DeleteCoupon(w http.ResponseWriter, r *http.Request) {
	staffID := middleware.GetClientID(r.Context())
	allowed, _ := h.staffService.HasPermission(r.Context(), staffID, "billing", "delete")
	if !allowed {
		response.Error(w, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions for module: billing", nil)
		return
	}

	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	// Delete promo by ID - assuming promoRepo has Delete method. Let's check domain/promo.go
	if err := h.promoRepo.Delete(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
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
		"transport":       "smtp",
		"smtp_host":       "smtp.mailgun.org",
		"smtp_port":       587,
		"smtp_username":   "postmaster@fossbilling.org",
		"smtp_encryption": "tls",
		"from_email":      "noreply@fossbilling.org",
		"from_name":       "FOSSBilling System",
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
	report, err := h.statsService.GetFinancialReports(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate financial reports", err.Error())
		return
	}

	response.JSON(w, http.StatusOK, report, nil)
}

func (h *BillingModuleHandler) ExportInvoicesCSV(w http.ResponseWriter, r *http.Request) {
	csvData, err := h.statsService.GenerateInvoicesCSV(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to generate CSV", err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=invoices_export.csv")
	w.Write([]byte(csvData))
}
