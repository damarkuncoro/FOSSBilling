import {
  ApiResponse,
  ProductItem,
  ProductCategory,
  ServerItem,
  PaymentGatewayItem,
  TaxRuleItem,
  CouponItem,
  EmailTemplateItem,
  MailConfig,
  StaffMemberItem,
  SecuritySettings,
  SystemStatusInfo,
  CompanySettings,
} from '@/types/api';
import {
  TldPricingItem,
  RegistrarConfig,
  CustomPageItem,
  KnowledgebaseArticle,
  ExtensionModuleItem,
  FinancialReportSummary,
} from '@/types/modules';

export * from '@/types/api';
export * from '@/types/modules';

const API_BASE = '/api/v1';

export class ApiError extends Error {
  code: string;
  details?: any;

  constructor(message: string, code = 'API_ERROR', details?: any) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.details = details;
  }
}

export const getStoredToken = () => localStorage.getItem('fossbilling_admin_token');
export const setStoredToken = (t: string) => localStorage.setItem('fossbilling_admin_token', t);
export const removeStoredToken = () => {
  localStorage.removeItem('fossbilling_admin_token');
  localStorage.removeItem('fossbilling_admin_user');
};

async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const token = getStoredToken();
  const headers = new Headers(options.headers || {});
  if (!headers.has('Content-Type') && !(options.body instanceof FormData)) headers.set('Content-Type', 'application/json');
  if (token && !headers.has('Authorization')) headers.set('Authorization', `Bearer ${token}`);

  const url = endpoint.startsWith('http') ? endpoint : `${API_BASE}${endpoint.startsWith('/') ? '' : '/'}${endpoint}`;
  const response = await fetch(url, { ...options, headers });
  const json: ApiResponse<T> = await response.json().catch(() => {
    throw new ApiError(`HTTP Error: ${response.status} ${response.statusText}`, 'HTTP_ERROR');
  });

  if (!response.ok || !json.success) {
    throw new ApiError(json.error?.message || `Status ${response.status}`, json.error?.code || 'UNKNOWN_ERROR', json.error?.details);
  }
  return json.data;
}

export const api = {
  login: (email: string, password: string) =>
    request<{ token: string; staff: any; group: any }>('/admin/auth/login', { method: 'POST', body: JSON.stringify({ email, password }) }),
  getDashboardStats: () => request<any>('/admin/stats/dashboard'),

  // Clients & Orders
  getClients: () => request<any[]>('/admin/clients'),
  getClient: (id: number) => request<any>(`/admin/clients/${id}`),
  createClient: (data: any) => request<any>('/admin/clients', { method: 'POST', body: JSON.stringify(data) }),
  updateClient: (id: number, data: any) => request<any>(`/admin/clients/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteClient: (id: number) => request<any>(`/admin/clients/${id}`, { method: 'DELETE' }),
  getOrders: () => request<any[]>('/admin/orders'),
  activateOrder: (id: number) => request<any>(`/admin/orders/${id}/activate`, { method: 'POST' }),
  suspendOrder: (id: number, reason = 'Administrative suspension') =>
    request<any>(`/admin/orders/${id}/suspend`, { method: 'POST', body: JSON.stringify({ reason }) }),
  unsuspendOrder: (id: number) => request<any>(`/admin/orders/${id}/unsuspend`, { method: 'POST' }),
  // Invoices & Billing
  getInvoices: () => request<any[]>('/admin/invoices'),
  createInvoice: (data: any) => request<any>('/admin/invoices', { method: 'POST', body: JSON.stringify(data) }),

  // Support & Currencies
  getSupportTickets: () => request<any[]>('/admin/support/tickets'),
  replySupportTicket: (id: number, content: string) =>
    request<any>(`/admin/support/tickets/${id}/reply`, { method: 'POST', body: JSON.stringify({ content }) }),
  getCurrencies: () => request<any[]>('/admin/currencies'),
  createCurrency: (currency: any) => request<any>('/admin/currencies', { method: 'POST', body: JSON.stringify(currency) }),
  setDefaultCurrency: (code: string) => request<any>(`/admin/currencies/${code}/default`, { method: 'POST' }),
  deleteCurrency: (code: string) => request<any>(`/admin/currencies/${code}`, { method: 'DELETE' }),

  // Products & Servers
  getProducts: () => request<ProductItem[]>('/admin/products'),
  createProduct: (product: Partial<ProductItem>) => request<ProductItem>('/admin/products', { method: 'POST', body: JSON.stringify(product) }),
  updateProduct: (id: number, product: Partial<ProductItem>) => request<ProductItem>(`/admin/products/${id}`, { method: 'PUT', body: JSON.stringify(product) }),
  deleteProduct: (id: number) => request<any>(`/admin/products/${id}`, { method: 'DELETE' }),
  getProductCategories: () => request<ProductCategory[]>('/admin/product-categories'),
  getServers: () => request<ServerItem[]>('/admin/servers'),
  createServer: (server: Partial<ServerItem>) => request<ServerItem>('/admin/servers', { method: 'POST', body: JSON.stringify(server) }),
  updateServer: (id: number, server: Partial<ServerItem>) => request<ServerItem>(`/admin/servers/${id}`, { method: 'PUT', body: JSON.stringify(server) }),
  testServerConnection: (id: number) => request<{ success: boolean; message: string }>(`/admin/servers/${id}/test`, { method: 'POST' }),
  deleteServer: (id: number) => request<any>(`/admin/servers/${id}`, { method: 'DELETE' }),

  // Gateways & Tax
  getPaymentGateways: () => request<PaymentGatewayItem[]>('/admin/gateways'),
  updatePaymentGateway: (id: string, config: Partial<PaymentGatewayItem>) =>
    request<PaymentGatewayItem>(`/admin/gateways/${id}`, { method: 'PUT', body: JSON.stringify(config) }),
  getTaxRules: () => request<TaxRuleItem[]>('/admin/tax-rules'),
  createTaxRule: (rule: Partial<TaxRuleItem>) => request<TaxRuleItem>('/admin/tax-rules', { method: 'POST', body: JSON.stringify(rule) }),
  deleteTaxRule: (id: number) => request<any>(`/admin/tax-rules/${id}`, { method: 'DELETE' }),

  // Coupons & Email
  getCoupons: () => request<CouponItem[]>('/admin/coupons'),
  createCoupon: (coupon: Partial<CouponItem>) => request<CouponItem>('/admin/coupons', { method: 'POST', body: JSON.stringify(coupon) }),
  deleteCoupon: (id: number) => request<any>(`/admin/coupons/${id}`, { method: 'DELETE' }),
  getEmailTemplates: () => request<EmailTemplateItem[]>('/admin/email-templates'),
  updateEmailTemplate: (id: string, template: Partial<EmailTemplateItem>) =>
    request<EmailTemplateItem>(`/admin/email-templates/${id}`, { method: 'PUT', body: JSON.stringify(template) }),
  getMailConfig: () => request<MailConfig>('/admin/settings/mail'),
  updateMailConfig: (config: Partial<MailConfig>) => request<MailConfig>('/admin/settings/mail', { method: 'PUT', body: JSON.stringify(config) }),
  sendTestEmail: (toEmail: string) => request<{ success: boolean; message: string }>('/admin/settings/mail/test', { method: 'POST', body: JSON.stringify({ email: toEmail }) }),

  // Staff & Security
  getStaffMembers: () => request<StaffMemberItem[]>('/admin/staff'),
  createStaffMember: (staff: Partial<StaffMemberItem> & { password?: string }) =>
    request<StaffMemberItem>('/admin/staff', { method: 'POST', body: JSON.stringify(staff) }),
  deleteStaffMember: (id: number) => request<any>(`/admin/staff/${id}`, { method: 'DELETE' }),
  getSecuritySettings: () => request<SecuritySettings>('/admin/settings/security'),
  updateSecuritySettings: (settings: Partial<SecuritySettings>) =>
    request<SecuritySettings>('/admin/settings/security', { method: 'PUT', body: JSON.stringify(settings) }),

  // System & Logs
  getSystemStatus: () => request<SystemStatusInfo>('/admin/system/status'),
  triggerCron: () => request<{ success: boolean; message: string; timestamp: string }>('/admin/system/cron/run', { method: 'POST' }),
  clearSystemCache: () => request<{ success: boolean; message: string }>('/admin/system/cache/clear', { method: 'POST' }),
  getNews: () => request<any[]>('/admin/news'),
  createNews: (article: { title: string; content: string }) => request<any>('/admin/news', { method: 'POST', body: JSON.stringify(article) }),
  deleteNews: (id: number) => request<any>(`/admin/news/${id}`, { method: 'DELETE' }),
  getMassMailCampaigns: () => request<any[]>('/admin/mass-mail'),
  createMassMailCampaign: (campaign: any) => request<any>('/admin/mass-mail', { method: 'POST', body: JSON.stringify(campaign) }),
  sendMassMailCampaign: (id: number) => request<any>(`/admin/mass-mail/${id}/send`, { method: 'POST' }),
  getCompany: () => request<CompanySettings>('/admin/company'),
  updateCompany: (settings: Partial<CompanySettings>) => request<CompanySettings>('/admin/company', { method: 'PUT', body: JSON.stringify(settings) }),
  getAuditLogs: () => request<any[]>('/admin/audit-logs'),

  // Domains & Registrars
  getTlds: () => request<TldPricingItem[]>('/admin/domains/tlds'),
  createTld: (tld: Partial<TldPricingItem>) => request<TldPricingItem>('/admin/domains/tlds', { method: 'POST', body: JSON.stringify(tld) }),
  deleteTld: (id: number) => request<any>(`/admin/domains/tlds/${id}`, { method: 'DELETE' }),
  getRegistrars: () => request<RegistrarConfig[]>('/admin/domains/registrars'),
  updateRegistrar: (id: string, config: Partial<RegistrarConfig>) =>
    request<RegistrarConfig>(`/admin/domains/registrars/${id}`, { method: 'PUT', body: JSON.stringify(config) }),

  // Pages & Knowledgebase
  getPages: () => request<CustomPageItem[]>('/admin/pages'),
  savePage: (page: Partial<CustomPageItem>) => request<CustomPageItem>('/admin/pages', { method: 'POST', body: JSON.stringify(page) }),
  deletePage: (id: number) => request<any>(`/admin/pages/${id}`, { method: 'DELETE' }),
  getKnowledgebase: () => request<KnowledgebaseArticle[]>('/admin/knowledgebase'),
  saveKnowledgebase: (article: Partial<KnowledgebaseArticle>) =>
    request<KnowledgebaseArticle>('/admin/knowledgebase', { method: 'POST', body: JSON.stringify(article) }),
  deleteKnowledgebase: (id: number) => request<any>(`/admin/knowledgebase/${id}`, { method: 'DELETE' }),

  // Extensions & Reports
  getExtensions: () => request<ExtensionModuleItem[]>('/admin/extensions'),
  toggleExtension: (id: string, enabled: boolean) =>
    request<any>(`/admin/extensions/${id}/toggle`, { method: 'POST', body: JSON.stringify({ enabled }) }),
  getFinancialReports: () => request<FinancialReportSummary>('/admin/reports/financial'),
};
