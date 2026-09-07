-- Enable UUID extension if PostgreSQL
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Clients Table
CREATE TABLE IF NOT EXISTS clients (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT NULL,
    email VARCHAR(191) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    company VARCHAR(150),
    address_1 VARCHAR(255),
    address_2 VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    postcode VARCHAR(20),
    country VARCHAR(2) DEFAULT 'US',
    phone_cc VARCHAR(10),
    phone VARCHAR(30),
    currency VARCHAR(3) DEFAULT 'USD',
    tax_exempt BOOLEAN DEFAULT FALSE,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_clients_email ON clients(email);
CREATE INDEX idx_clients_status ON clients(status);

-- 2. Client Balance Ledger
CREATE TABLE IF NOT EXISTS client_balances (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL, -- 'credit', 'debit'
    amount BIGINT NOT NULL,    -- in hundredths of a cent (e.g. $10.00 = 100000)
    description TEXT,
    rel_id BIGINT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_client_balances_client_id ON client_balances(client_id);

-- 3. Products Table
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NULL,
    type VARCHAR(50) NOT NULL, -- 'hosting', 'domain', 'license', 'downloadable', 'custom'
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    status VARCHAR(50) DEFAULT 'enabled',
    setup_type VARCHAR(50) DEFAULT 'recurring',
    config JSONB NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_slug ON products(slug);
CREATE INDEX idx_products_type ON products(type);

-- 4. Invoices Table
CREATE TABLE IF NOT EXISTS invoices (
    id BIGSERIAL PRIMARY KEY,
    serie VARCHAR(50) NOT NULL,
    nr VARCHAR(50) NOT NULL,
    client_id BIGINT NOT NULL REFERENCES clients(id) ON DELETE RESTRICT,
    status VARCHAR(50) DEFAULT 'unpaid', -- 'unpaid', 'paid', 'refunded', 'canceled'
    currency VARCHAR(3) DEFAULT 'USD',
    currency_rate NUMERIC(18, 6) DEFAULT 1.000000,
    subtotal BIGINT DEFAULT 0,
    tax BIGINT DEFAULT 0,
    total BIGINT DEFAULT 0,
    tax_rate NUMERIC(6, 2) DEFAULT 0.00,
    due_at DATE NOT NULL,
    paid_at TIMESTAMP WITH TIME ZONE NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_invoices_client_id ON invoices(client_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_due_at ON invoices(due_at);

-- 5. Invoice Items Table
CREATE TABLE IF NOT EXISTS invoice_items (
    id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    order_id BIGINT NULL,
    title VARCHAR(255) NOT NULL,
    period VARCHAR(20) NULL,
    price BIGINT NOT NULL,
    quantity INT DEFAULT 1,
    unit VARCHAR(50) DEFAULT 'unit',
    taxable BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_invoice_items_invoice_id ON invoice_items(invoice_id);

-- 6. Orders Table
CREATE TABLE IF NOT EXISTS client_orders (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES clients(id) ON DELETE RESTRICT,
    product_id BIGINT NULL REFERENCES products(id) ON DELETE SET NULL,
    invoice_id BIGINT NULL,
    status VARCHAR(50) DEFAULT 'pending_setup', -- 'pending_setup', 'active', 'suspended', 'canceled', 'terminated'
    title VARCHAR(255) NOT NULL,
    period VARCHAR(20) NOT NULL,
    price BIGINT NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    config JSONB NULL,
    activated_at TIMESTAMP WITH TIME ZONE NULL,
    expires_at TIMESTAMP WITH TIME ZONE NULL,
    next_due_date DATE NULL,
    suspended_at TIMESTAMP WITH TIME ZONE NULL,
    suspension_reason TEXT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_orders_client_id ON client_orders(client_id);
CREATE INDEX idx_orders_status ON client_orders(status);
CREATE INDEX idx_orders_next_due_date ON client_orders(next_due_date);

-- 7. Transactions Table
CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NULL REFERENCES invoices(id) ON DELETE SET NULL,
    gateway_id VARCHAR(50) NOT NULL,
    txn_id VARCHAR(191) NOT NULL,
    type VARCHAR(50) DEFAULT 'payment',
    amount BIGINT NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    status VARCHAR(50) DEFAULT 'pending',
    raw_payload JSONB NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_invoice_id ON transactions(invoice_id);
CREATE INDEX idx_transactions_txn_id ON transactions(txn_id);

-- 8. Support Tickets Table
CREATE TABLE IF NOT EXISTS support_tickets (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    helpdesk_id BIGINT NOT NULL,
    subject VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'open',
    priority VARCHAR(20) DEFAULT 'medium',
    rel_type VARCHAR(50) NULL,
    rel_id BIGINT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tickets_client_id ON support_tickets(client_id);
CREATE INDEX idx_tickets_status ON support_tickets(status);

-- 9. Ticket Messages Table
CREATE TABLE IF NOT EXISTS support_ticket_messages (
    id BIGSERIAL PRIMARY KEY,
    ticket_id BIGINT NOT NULL REFERENCES support_tickets(id) ON DELETE CASCADE,
    admin_id BIGINT NULL,
    client_id BIGINT NULL,
    content TEXT NOT NULL,
    ip_address VARCHAR(45) NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ticket_messages_ticket_id ON support_ticket_messages(ticket_id);

-- 10. Admin Groups Table
CREATE TABLE IF NOT EXISTS admin_groups (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    permissions JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 11. Staff Table
CREATE TABLE IF NOT EXISTS staff (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL REFERENCES admin_groups(id) ON DELETE RESTRICT,
    email VARCHAR(191) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(150) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'admin', -- 'superadmin', 'admin', 'support', 'billing'
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_staff_email ON staff(email);

-- 12. Audit Logs Table
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    staff_id BIGINT NULL REFERENCES staff(id) ON DELETE SET NULL,
    client_id BIGINT NULL REFERENCES clients(id) ON DELETE SET NULL,
    module VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    details TEXT NOT NULL,
    ip_address VARCHAR(45) NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_staff_id ON audit_logs(staff_id);
CREATE INDEX idx_audit_logs_module ON audit_logs(module);

-- 12b. Activity Logs Table
CREATE TABLE IF NOT EXISTS activity_logs (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NULL REFERENCES clients(id) ON DELETE SET NULL,
    admin_id BIGINT NULL REFERENCES staff(id) ON DELETE SET NULL,
    type VARCHAR(50) NOT NULL,
    event VARCHAR(100) NOT NULL,
    message TEXT NOT NULL,
    ip_address VARCHAR(45) NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_activity_logs_client_id ON activity_logs(client_id);
CREATE INDEX idx_activity_logs_type ON activity_logs(type);

-- 13. Promos Table
CREATE TABLE IF NOT EXISTS promos (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    type VARCHAR(50) NOT NULL, -- 'percentage', 'absolute'
    value BIGINT NOT NULL,
    max_uses INT NOT NULL DEFAULT 0,
    used_count INT NOT NULL DEFAULT 0,
    once_per_client BOOLEAN NOT NULL DEFAULT FALSE,
    start_date TIMESTAMP WITH TIME ZONE NULL,
    end_date TIMESTAMP WITH TIME ZONE NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_promos_code ON promos(code);

-- 14. Promo Redemptions Table
CREATE TABLE IF NOT EXISTS promo_redemptions (
    id BIGSERIAL PRIMARY KEY,
    promo_id BIGINT NOT NULL REFERENCES promos(id) ON DELETE CASCADE,
    client_id BIGINT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    order_id BIGINT NULL REFERENCES client_orders(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 15. Currencies Table
CREATE TABLE IF NOT EXISTS currencies (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(3) NOT NULL UNIQUE,
    title VARCHAR(50) NOT NULL,
    conversion_rate NUMERIC(18, 6) NOT NULL DEFAULT 1.000000,
    format VARCHAR(50) DEFAULT '{{price}} {{code}}',
    price_format VARCHAR(50) DEFAULT '1',
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_currencies_default ON currencies(is_default);
CREATE INDEX idx_currencies_code ON currencies(code);

-- 16. News Posts Table
CREATE TABLE IF NOT EXISTS news_posts (
    id BIGSERIAL PRIMARY KEY,
    admin_id BIGINT NULL REFERENCES staff(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'published',
    published_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_news_slug ON news_posts(slug);
CREATE INDEX idx_news_status ON news_posts(status);

-- 17. Downloadable Files Table
CREATE TABLE IF NOT EXISTS downloadable_files (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL DEFAULT 0,
    content_type VARCHAR(100) NOT NULL DEFAULT 'application/octet-stream',
    version VARCHAR(50) NOT NULL DEFAULT '1.0.0',
    downloads INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_downloadable_product_id ON downloadable_files(product_id);

-- 18. API Keys Table
CREATE TABLE IF NOT EXISTS api_keys (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    key VARCHAR(64) NOT NULL UNIQUE,
    secret VARCHAR(64) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_api_keys_client_id ON api_keys(client_id);
CREATE INDEX idx_api_keys_key ON api_keys(key);

-- 19. Mass Mail Campaigns Table
CREATE TABLE IF NOT EXISTS mass_mail_campaigns (
    id BIGSERIAL PRIMARY KEY,
    admin_id BIGINT NULL REFERENCES staff(id) ON DELETE SET NULL,
    subject VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    sent_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    sent_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE INDEX idx_mass_mail_status ON mass_mail_campaigns(status);

-- 20. Notifications Table
CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'info',
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_client_id ON notifications(client_id);

-- 21. Product Categories Table
CREATE TABLE IF NOT EXISTS product_categories (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 22. Servers Table
CREATE TABLE IF NOT EXISTS servers (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    ip VARCHAR(45) NOT NULL,
    manager VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    is_default BOOLEAN DEFAULT FALSE,
    max_accounts INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 23. TLDs Table
CREATE TABLE IF NOT EXISTS tlds (
    id BIGSERIAL PRIMARY KEY,
    tld VARCHAR(50) NOT NULL UNIQUE,
    registrar_id VARCHAR(50) NOT NULL,
    price_registration BIGINT NOT NULL,
    price_renewal BIGINT NOT NULL,
    price_transfer BIGINT NOT NULL,
    min_years INT DEFAULT 1,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 24. Pages Table
CREATE TABLE IF NOT EXISTS pages (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    published BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 25. Company Settings Table
CREATE TABLE IF NOT EXISTS company_settings (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(191) NOT NULL,
    phone VARCHAR(50),
    address VARCHAR(255),
    logo_url VARCHAR(500),
    favicon_url VARCHAR(500),
    timezone VARCHAR(50) DEFAULT 'UTC',
    date_format VARCHAR(50) DEFAULT 'Y-m-d',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 26. System Settings Table
CREATE TABLE IF NOT EXISTS system_settings (
    id BIGSERIAL PRIMARY KEY,
    section VARCHAR(50) NOT NULL,
    key VARCHAR(100) NOT NULL,
    value TEXT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(section, key)
);

-- 27. Blocked IPs Table (Antispam)
CREATE TABLE IF NOT EXISTS blocked_ips (
    id BIGSERIAL PRIMARY KEY,
    ip VARCHAR(45) NOT NULL UNIQUE,
    reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_blocked_ips_ip ON blocked_ips(ip);

-- 28. Custom Forms Table (Formbuilder)
CREATE TABLE IF NOT EXISTS forms (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    style JSONB NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 29. Custom Form Fields Table (Formbuilder)
CREATE TABLE IF NOT EXISTS form_fields (
    id BIGSERIAL PRIMARY KEY,
    form_id BIGINT NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(255) NOT NULL,
    hide_label BOOLEAN DEFAULT FALSE,
    description TEXT,
    type VARCHAR(50) NOT NULL, -- 'text', 'url', 'select', 'radio', 'checkbox', 'textarea'
    default_value TEXT,
    required BOOLEAN DEFAULT FALSE,
    hidden BOOLEAN DEFAULT FALSE,
    readonly BOOLEAN DEFAULT FALSE,
    options JSONB NULL,
    prefix VARCHAR(50),
    suffix VARCHAR(50),
    text_size INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(form_id, name)
);
CREATE INDEX idx_form_fields_form_id ON form_fields(form_id);

-- 30. Extensions Table
CREATE TABLE IF NOT EXISTS extensions (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'mod', 'theme', 'gateway', 'registrar', 'plugin', 'service'
    version VARCHAR(50) NOT NULL,
    description TEXT,
    author VARCHAR(255) NOT NULL,
    author_url VARCHAR(500),
    icon VARCHAR(500),
    status VARCHAR(50) DEFAULT 'active', -- 'active', 'inactive', 'core'
    has_settings BOOLEAN DEFAULT FALSE,
    config JSONB NULL,
    manifest JSONB NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_extensions_type ON extensions(type);
CREATE INDEX idx_extensions_status ON extensions(status);

-- 31. Redirects Table
CREATE TABLE IF NOT EXISTS redirects (
    id BIGSERIAL PRIMARY KEY,
    path VARCHAR(500) NOT NULL UNIQUE,
    target VARCHAR(1000) NOT NULL,
    status_code INT DEFAULT 301,
    is_enabled BOOLEAN DEFAULT TRUE,
    hit_count BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_redirects_path ON redirects(path);

-- 20. Seed Data
-- Super Admin Group & Staff
INSERT INTO admin_groups (id, name, permissions)
VALUES (1, 'Super Administrator', '{"all": true}'::jsonb)
ON CONFLICT (id) DO NOTHING;

INSERT INTO staff (id, group_id, name, email, password_hash, role, status)
VALUES (1, 1, 'Administrator', 'admin@fossbilling.org', '$2a$12$pqYB3.w3AIe4Kih5bDvjI.lv39gEy0mWy//TGPMgvU2vysmqqgSZW', 'superadmin', 'active')
ON CONFLICT (email) DO NOTHING;

-- Default Currencies
INSERT INTO currencies (code, title, conversion_rate, format, price_format, is_default)
VALUES 
    ('USD', 'US Dollar', 1.000000, '$ {{price}}', '2', TRUE),
    ('EUR', 'Euro', 0.920000, '{{price}} €', '2', FALSE),
    ('IDR', 'Indonesian Rupiah', 15800.000000, 'Rp {{price}}', '0', FALSE)
ON CONFLICT (code) DO NOTHING;

-- Default Products
INSERT INTO products (id, category_id, type, name, slug, description, status, setup_type)
VALUES
    (1, NULL, 'hosting', 'cPanel Starter Cloud', 'cpanel-starter-cloud', 'Perfect for personal blogs and portfolios.', 'enabled', 'recurring'),
    (2, NULL, 'vps', 'Cloud VPS Pro (DirectAdmin)', 'cloud-vps-pro-directadmin', 'High performance dedicated computing.', 'enabled', 'recurring'),
    (3, NULL, 'license', 'FOSSBilling Enterprise License', 'fossbilling-enterprise-license', 'Self-hosted enterprise license with priority SLA.', 'enabled', 'recurring'),
    (4, NULL, 'downloadable', 'Nusantara Cloud OS Template', 'nusantara-cloud-os-template', 'Pre-hardened Linux image with automated docker deployments.', 'enabled', 'onetime'),
    (10, NULL, 'domain', 'Domain Registration', 'domain-registration', 'Standard domain registration service.', 'enabled', 'recurring'),
    (99, NULL, 'domain', 'Domain TLD Registration', 'domain-tld-registration', 'Instant domain name registration.', 'enabled', 'recurring'),
    (101, NULL, 'hosting', 'Cloud VPS cPanel Pro', 'cloud-vps-cpanel-pro', 'Advanced Cloud VPS with cPanel.', 'enabled', 'recurring'),
    (202, NULL, 'hosting', 'DirectAdmin Hosting', 'directadmin-hosting', 'Budget friendly Linux hosting with DirectAdmin.', 'enabled', 'recurring'),
    (303, NULL, 'license', 'FOSSBilling Enterprise', 'fossbilling-enterprise', 'Enterprise self-hosted software license.', 'enabled', 'recurring'),
    (404, NULL, 'downloadable', 'Nusantara Cloud OS', 'nusantara-cloud-os', 'Signed enterprise OS image.', 'enabled', 'onetime')
ON CONFLICT (id) DO NOTHING;

-- Reset sequences
SELECT setval('products_id_seq', (SELECT COALESCE(MAX(id), 1) FROM products));
SELECT setval('admin_groups_id_seq', (SELECT COALESCE(MAX(id), 1) FROM admin_groups));
SELECT setval('staff_id_seq', (SELECT COALESCE(MAX(id), 1) FROM staff));
SELECT setval('currencies_id_seq', (SELECT COALESCE(MAX(id), 1) FROM currencies));

-- Default Promos
INSERT INTO promos (code, description, type, value, max_uses, active)
VALUES 
    ('MERDEKA20', '20% Discount for Merdeka Promo', 'percentage', 2000, 1000, TRUE),
    ('WELCOME10', '10% Welcome Discount', 'percentage', 1000, 1000, TRUE)
ON CONFLICT (code) DO NOTHING;

-- Default News
INSERT INTO news_posts (admin_id, title, slug, content, status)
VALUES 
    (1, 'Welcome to Next-Gen FOSSBilling', 'welcome-to-next-gen-fossbilling', 'We are excited to launch the high performance Cloud-Native Go & React edition!', 'published')
ON CONFLICT (slug) DO NOTHING;

-- Default Extensions
INSERT INTO extensions (id, name, type, version, description, author, status, has_settings)
VALUES
    ('antispam', 'Anti-Spam & Abuse Shield', 'plugin', '2.0.0', 'StopForumSpam, Cloudflare Turnstile, and temporary email protection.', 'FOSSBilling Core Team', 'active', TRUE),
    ('formbuilder', 'Custom Order Formbuilder', 'mod', '2.0.0', 'Dynamic custom checkout fields for products and provisioning options.', 'FOSSBilling Core Team', 'active', TRUE),
    ('servicehosting', 'cPanel & DirectAdmin Hosting Provisioner', 'service', '2.1.0', 'Automated provisioning for shared and reseller hosting servers.', 'FOSSBilling Core Team', 'active', TRUE),
    ('midtrans', 'Midtrans Payment Gateway', 'gateway', '1.3.0', 'SNAP, QRIS, Virtual Account, and Credit Card payments via Midtrans.', 'Nusantara Developers', 'active', TRUE),
    ('stripe', 'Stripe Checkout & Elements', 'gateway', '2.0.0', 'Global credit card and debit payment processing.', 'FOSSBilling Core Team', 'active', TRUE),
    ('cookieconsent', 'EU Cookie Consent Compliance', 'mod', '1.1.0', 'GDPR compliant cookie consent banner with configurable categories.', 'FOSSBilling Core Team', 'active', TRUE)
ON CONFLICT (id) DO NOTHING;

-- Default Provisioning Servers
INSERT INTO servers (id, name, hostname, ip, manager, status, is_default, max_accounts)
VALUES
    (1, 'Primary cPanel Cluster', 'cpanel.fossbilling.org', '198.51.100.10', 'cpanel', 'active', TRUE, 250),
    (2, 'DirectAdmin Cloud Node', 'da.fossbilling.org', '198.51.100.11', 'directadmin', 'active', FALSE, 500),
    (3, 'HestiaCP Performance Node', 'hestia.fossbilling.org', '198.51.100.12', 'hestia', 'active', FALSE, 150),
    (4, 'Plesk Web Cluster', 'plesk.fossbilling.org', '198.51.100.13', 'plesk', 'active', FALSE, 300),
    (5, 'CentOS Web Panel Node', 'cwp.fossbilling.org', '198.51.100.14', 'cwp', 'active', FALSE, 200)
ON CONFLICT (id) DO NOTHING;

SELECT setval('servers_id_seq', (SELECT COALESCE(MAX(id), 1) FROM servers));




