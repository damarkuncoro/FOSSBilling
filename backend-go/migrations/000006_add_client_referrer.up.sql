ALTER TABLE clients ADD COLUMN IF NOT EXISTS referrer_id BIGINT REFERENCES affiliates(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_clients_referrer_id ON clients(referrer_id);
