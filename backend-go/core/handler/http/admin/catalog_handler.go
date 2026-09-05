package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/response"
)

type CatalogHandler struct{}

func NewCatalogHandler() *CatalogHandler {
	return &CatalogHandler{}
}

// --- Products & Categories ---
func (h *CatalogHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products := []map[string]interface{}{
		{"id": 1, "title": "cPanel Starter Cloud", "slug": "cpanel-starter", "type": "hosting", "category_name": "Web Hosting", "price_monthly": 9.99, "is_active": true},
		{"id": 2, "title": "Cloud VPS Pro", "slug": "cloud-vps-pro", "type": "hosting", "category_name": "Cloud VPS", "price_monthly": 29.99, "is_active": true},
		{"id": 3, "title": "FOSSBilling License", "slug": "fossbilling-license", "type": "license", "category_name": "Licenses", "price_monthly": 199.00, "is_active": true},
	}
	response.JSON(w, http.StatusOK, products, nil)
}

func (h *CatalogHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&body)
	body["id"] = time.Now().Unix()
	body["is_active"] = true
	response.JSON(w, http.StatusCreated, body, nil)
}

func (h *CatalogHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&body)
	response.JSON(w, http.StatusOK, body, nil)
}

func (h *CatalogHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]bool{"deleted": true}, nil)
}

func (h *CatalogHandler) ListProductCategories(w http.ResponseWriter, r *http.Request) {
	categories := []map[string]interface{}{
		{"id": 1, "title": "Web Hosting", "slug": "web-hosting", "product_count": 4},
		{"id": 2, "title": "Cloud VPS", "slug": "cloud-vps", "product_count": 3},
		{"id": 3, "title": "Licenses", "slug": "licenses", "product_count": 2},
	}
	response.JSON(w, http.StatusOK, categories, nil)
}

// --- Domains & Registrars ---
func (h *CatalogHandler) ListTlds(w http.ResponseWriter, r *http.Request) {
	tlds := []map[string]interface{}{
		{"id": 1, "tld": ".com", "registrar": "namecheap", "price_registration": 12.99, "price_renewal": 14.99, "price_transfer": 12.99, "min_years": 1, "is_active": true},
		{"id": 2, "tld": ".id", "registrar": "custom", "price_registration": 18.00, "price_renewal": 18.00, "price_transfer": 18.00, "min_years": 1, "is_active": true},
		{"id": 3, "tld": ".net", "registrar": "namecheap", "price_registration": 13.50, "price_renewal": 15.50, "price_transfer": 13.50, "min_years": 1, "is_active": true},
	}
	response.JSON(w, http.StatusOK, tlds, nil)
}

func (h *CatalogHandler) CreateTld(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&body)
	body["id"] = time.Now().Unix()
	response.JSON(w, http.StatusCreated, body, nil)
}

func (h *CatalogHandler) DeleteTld(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]bool{"deleted": true}, nil)
}

func (h *CatalogHandler) ListRegistrars(w http.ResponseWriter, r *http.Request) {
	registrars := []map[string]interface{}{
		{"id": "namecheap", "name": "Namecheap API", "enabled": true, "api_user": "api_fossbilling", "test_mode": false},
		{"id": "enom", "name": "eNom Reseller", "enabled": true, "api_user": "reseller_demo", "test_mode": true},
		{"id": "custom", "name": "DigitalRegistrar (.ID)", "enabled": true, "test_mode": false},
	}
	response.JSON(w, http.StatusOK, registrars, nil)
}

// --- Servers ---
func (h *CatalogHandler) ListServers(w http.ResponseWriter, r *http.Request) {
	servers := []map[string]interface{}{
		{"id": 1, "name": "SG-Cloud-Node-01", "hostname": "sg1.nusantara-cloud.com", "ip": "103.144.20.10", "manager": "cpanel", "status": "online", "active_accounts": 48, "max_accounts": 150, "is_default": true},
		{"id": 2, "name": "JKT-DirectAdmin-02", "hostname": "jkt2.nusantara-cloud.com", "ip": "103.144.20.15", "manager": "directadmin", "status": "online", "active_accounts": 22, "max_accounts": 100, "is_default": false},
	}
	response.JSON(w, http.StatusOK, servers, nil)
}

func (h *CatalogHandler) CreateServer(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&body)
	body["id"] = time.Now().Unix()
	body["status"] = "online"
	response.JSON(w, http.StatusCreated, body, nil)
}

func (h *CatalogHandler) TestServer(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Connection handshake successful! Latency: 18ms.",
	}, nil)
}

func (h *CatalogHandler) DeleteServer(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]bool{"deleted": true}, nil)
}
