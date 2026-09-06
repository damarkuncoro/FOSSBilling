export interface ApiResponse<T = any> {
  success: boolean;
  data: T;
  error?: {
    code: string;
    message: string;
    details?: any;
  } | null;
}

export interface ClientProfile {
  id: number;
  first_name: string;
  last_name: string;
  email: string;
  company?: string;
  country?: string;
  currency?: string;
  role?: string;
  status?: string;
  created_at?: string;
  updated_at?: string;
}

export interface Order {
  id: number;
  client_id: number;
  product_id: number;
  server_id?: number;
  title?: string;
  status: 'active' | 'pending' | 'suspended' | 'canceled';
  period: string;
  price: number;
  currency?: string;
  config?: Record<string, any>;
  created_at?: string;
  updated_at?: string;
}

export interface Invoice {
  id: number;
  client_id: number;
  serie_nr?: string;
  status: 'paid' | 'unpaid' | 'canceled' | 'refunded';
  currency: string;
  subtotal: number;
  discount: number;
  tax: number;
  total: number;
  due_at?: string;
  created_at: string;
  updated_at?: string;
}

export interface SupportTicket {
  id: number;
  client_id: number;
  subject: string;
  content?: string;
  status: 'open' | 'answered' | 'closed' | 'on_hold';
  priority: 'low' | 'medium' | 'high' | 'urgent';
  created_at: string;
  updated_at?: string;
  replies?: Array<{
    id: number;
    author: string;
    content: string;
    created_at: string;
  }>;
}

export interface ApiKey {
  id: number;
  name: string;
  key?: string;
  created_at?: string;
}

export interface PublicCompanyInfo {
  name: string;
  email: string;
  phone?: string;
  address_1?: string;
  address_2?: string;
  city?: string;
  state?: string;
  postcode?: string;
  country?: string;
  vat_number?: string;
  logo_url?: string;
  logo_dark_url?: string;
  favicon_url?: string;
  terms_url?: string;
}

export interface HostingPlan {
  id: number;
  title: string;
  description: string;
  price: number;
  period: string;
  type: string;
  features: string[];
  popular?: boolean;
}

export interface DomainSearchResult {
  domain: string;
  tld?: string;
  available: boolean;
  price: number;
  currency: string;
}

export interface CartCalculation {
  subtotal: number;
  discount: number;
  tax: number;
  total: number;
  items: any[];
}
