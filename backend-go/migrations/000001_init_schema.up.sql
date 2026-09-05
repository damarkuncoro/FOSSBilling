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
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
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

CREATE INDEX idx_promo_redemptions_client_id ON promo_redemptions(client_id);

