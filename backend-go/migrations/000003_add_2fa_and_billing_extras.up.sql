-- 33. Security Additions
ALTER TABLE clients ADD COLUMN IF NOT EXISTS two_factor_enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE clients ADD COLUMN IF NOT EXISTS two_factor_secret VARCHAR(100) NULL;

ALTER TABLE staff ADD COLUMN IF NOT EXISTS two_factor_enabled BOOLEAN DEFAULT FALSE;
ALTER TABLE staff ADD COLUMN IF NOT EXISTS two_factor_secret VARCHAR(100) NULL;

-- 32. Admin Notifications Table
CREATE TABLE IF NOT EXISTS admin_notifications (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'info', -- 'info', 'warning', 'danger', 'success'
    module VARCHAR(50) NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_admin_notifications_read ON admin_notifications(is_read);

-- 34. Tax Rules Table
CREATE TABLE IF NOT EXISTS tax_rules (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(2) NULL, -- NULL means all countries
    state VARCHAR(100) NULL,
    rate NUMERIC(6, 2) NOT NULL DEFAULT 0.00,
    is_active BOOLEAN DEFAULT TRUE,
    tax_exempt BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tax_rules_country ON tax_rules(country);

-- Seed Initial Tax Rule
INSERT INTO tax_rules (name, country, rate, is_active)
VALUES ('Indonesia PPN', 'ID', 11.00, true)
ON CONFLICT DO NOTHING;
