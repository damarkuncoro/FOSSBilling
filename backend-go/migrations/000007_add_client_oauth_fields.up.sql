ALTER TABLE clients ADD COLUMN IF NOT EXISTS oauth_provider VARCHAR(50);
ALTER TABLE clients ADD COLUMN IF NOT EXISTS oauth_id VARCHAR(255);
CREATE INDEX IF NOT EXISTS idx_clients_oauth ON clients(oauth_provider, oauth_id);
