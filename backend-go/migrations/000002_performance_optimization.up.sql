-- Performance Optimization Migration

-- 1. Invoices: Speed up listing by client and sorting by date
CREATE INDEX IF NOT EXISTS idx_invoices_client_created ON invoices(client_id, created_at DESC);

-- 2. Orders: Speed up listing by client and sorting by date
CREATE INDEX IF NOT EXISTS idx_orders_client_created ON client_orders(client_id, created_at DESC);

-- 3. Support Tickets: Speed up listing and sorting
CREATE INDEX IF NOT EXISTS idx_tickets_client_updated ON support_tickets(client_id, updated_at DESC);

-- 4. Activity & Audit Logs: High volume tables need date-based sorting indexes
CREATE INDEX IF NOT EXISTS idx_activity_logs_created ON activity_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs(created_at DESC);

-- 5. Admin Notifications: Dashboard alerts
CREATE INDEX IF NOT EXISTS idx_admin_notifications_created ON admin_notifications(created_at DESC);

-- 6. Foreign Key Lookups (Missing from init)
CREATE INDEX IF NOT EXISTS idx_invoice_items_order_id ON invoice_items(order_id);
CREATE INDEX IF NOT EXISTS idx_client_orders_product_id ON client_orders(product_id);
CREATE INDEX IF NOT EXISTS idx_client_orders_invoice_id ON client_orders(invoice_id);

-- 7. Search Optimization (B-Tree for prefix matching)
CREATE INDEX IF NOT EXISTS idx_clients_name ON clients(last_name, first_name);
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
