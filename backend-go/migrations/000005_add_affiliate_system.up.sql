-- Affiliate System Tables
CREATE TABLE IF NOT EXISTS affiliates (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL UNIQUE REFERENCES clients(id) ON DELETE CASCADE,
    commission_rate DOUBLE PRECISION NOT NULL DEFAULT 10.0,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    balance BIGINT NOT NULL DEFAULT 0,
    total_earned BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS affiliate_referrals (
    id BIGSERIAL PRIMARY KEY,
    affiliate_id BIGINT NOT NULL REFERENCES affiliates(id) ON DELETE CASCADE,
    client_id BIGINT NOT NULL REFERENCES clients(id),
    order_id BIGINT NOT NULL REFERENCES client_orders(id),
    amount BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Index for performance
CREATE INDEX idx_affiliate_referrals_affiliate_id ON affiliate_referrals(affiliate_id);
CREATE INDEX idx_affiliate_referrals_client_id ON affiliate_referrals(client_id);
