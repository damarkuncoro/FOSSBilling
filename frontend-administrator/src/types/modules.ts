export interface TldPricingItem {
  id: number;
  tld: string; // e.g. .com, .id, .net
  registrar: string; // namecheap, enom, custom
  price_registration: number;
  price_renewal: number;
  price_transfer: number;
  min_years: number;
  is_active: boolean;
}

export interface RegistrarConfig {
  id: string;
  name: string;
  enabled: boolean;
  api_user?: string;
  api_key?: string;
  test_mode: boolean;
}

export interface CustomPageItem {
  id: number;
  title: string;
  slug: string;
  content: string;
  published: boolean;
  meta_title?: string;
  meta_description?: string;
  updated_at?: string;
}

export interface KnowledgebaseArticle {
  id: number;
  category: string;
  title: string;
  slug: string;
  content: string;
  views: number;
  published: boolean;
  updated_at?: string;
}

export interface ExtensionModuleItem {
  id: string;
  name: string;
  version: string;
  author: string;
  description: string;
  type: 'service' | 'gateway' | 'plugin' | 'theme';
  is_installed: boolean;
  is_enabled: boolean;
  has_update?: boolean;
}

export interface FinancialReportSummary {
  mrr: number;
  arr: number;
  total_revenue_month: number;
  total_tax_collected: number;
  active_subscriptions: number;
  churn_rate: number;
  monthly_breakdown: Array<{
    month: string;
    revenue: number;
    tax: number;
    invoices_count: number;
  }>;
}
