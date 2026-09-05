export interface ApiResponse<T = any> {
  success: boolean;
  data: T;
  error?: {
    code: string;
    message: string;
    details?: any;
  } | null;
}

export interface ProductItem {
  id: number;
  title: string;
  slug: string;
  type: 'hosting' | 'domain' | 'license' | 'downloadable' | 'custom';
  category_id?: number;
  category_name?: string;
  description: string;
  price_monthly: number;
  price_annually?: number;
  setup_fee?: number;
  is_active: boolean;
  stock?: number;
  created_at?: string;
}

export interface ProductCategory {
  id: number;
  title: string;
  slug: string;
  description?: string;
  product_count?: number;
}

export interface ServerItem {
  id: number;
  name: string;
  hostname: string;
  ip: string;
  manager: 'cpanel' | 'hestiacp' | 'cwp' | 'directadmin' | 'plesk' | 'custom';
  status: 'online' | 'offline' | 'unreachable';
  active_accounts: number;
  max_accounts: number;
  nameserver_1?: string;
  nameserver_2?: string;
  is_default: boolean;
}

export interface PaymentGatewayItem {
  id: string;
  name: string;
  description: string;
  type: 'card' | 'wallet' | 'bank_transfer' | 'crypto';
  enabled: boolean;
  test_mode: boolean;
  public_key?: string;
  secret_key?: string;
  webhook_secret?: string;
  instructions?: string;
}

export interface TaxRuleItem {
  id: number;
  name: string;
  country: string;
  state?: string;
  rate: number;
  is_active: boolean;
  apply_to_all_clients: boolean;
}

export interface CouponItem {
  id: number;
  code: string;
  type: 'percentage' | 'fixed';
  value: number;
  max_uses: number;
  used_count: number;
  expires_at?: string;
  is_active: boolean;
}

export interface EmailTemplateItem {
  id: string;
  code: string;
  subject: string;
  description: string;
  category: 'client' | 'invoice' | 'service' | 'support';
  content: string;
  enabled: boolean;
  variables: string[];
}

export interface MailConfig {
  transport: 'smtp' | 'sendmail' | 'phpmail';
  smtp_host: string;
  smtp_port: number;
  smtp_username: string;
  smtp_password?: string;
  smtp_encryption: 'none' | 'tls' | 'ssl';
  from_email: string;
  from_name: string;
}

export interface StaffMemberItem {
  id: number;
  name: string;
  email: string;
  role: 'superadmin' | 'admin' | 'support' | 'billing';
  status: 'active' | 'inactive';
  last_login?: string;
  created_at?: string;
}

export interface SecuritySettings {
  recaptcha_enabled: boolean;
  recaptcha_provider: 'cloudflare_turnstile' | 'google_recaptcha';
  site_key: string;
  secret_key?: string;
  ip_blacklist: string[];
  max_login_attempts: number;
  lockout_time_minutes: number;
  force_ssl: boolean;
}

export interface SystemStatusInfo {
  engine_version: string;
  php_version?: string;
  go_version?: string;
  database_type: string;
  database_size: string;
  active_sessions: number;
  cron_last_run: string;
  cron_status: 'healthy' | 'warning' | 'error';
  system_load: string;
  memory_usage: string;
  uptime: string;
}

export interface CompanySettings {
  id?: number;
  name: string;
  email: string;
  phone: string;
  address_1: string;
  address_2?: string;
  city: string;
  state: string;
  postcode: string;
  country: string;
  vat_number?: string;
  logo_url?: string;
  logo_dark_url?: string;
  favicon_url?: string;
  terms_url?: string;
  email_signature?: string;
  updated_at?: string;
}
