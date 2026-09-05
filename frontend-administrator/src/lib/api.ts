const API_BASE = '/api/v1';

export interface ApiResponse<T = any> {
  success: boolean;
  data: T;
  error?: {
    code: string;
    message: string;
    details?: any;
  } | null;
}

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

export function getStoredToken(): string | null {
  return localStorage.getItem('fossbilling_admin_token');
}

export function setStoredToken(token: string) {
  localStorage.setItem('fossbilling_admin_token', token);
}

export function removeStoredToken() {
  localStorage.removeItem('fossbilling_admin_token');
  localStorage.removeItem('fossbilling_admin_user');
}

async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const token = getStoredToken();
  const headers = new Headers(options.headers || {});

  if (!headers.has('Content-Type') && !(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json');
  }

  if (token && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  const url = endpoint.startsWith('http') ? endpoint : `${API_BASE}${endpoint.startsWith('/') ? '' : '/'}${endpoint}`;

  const response = await fetch(url, {
    ...options,
    headers,
  });

  const json: ApiResponse<T> = await response.json().catch(() => {
    throw new ApiError(`HTTP Error: ${response.status} ${response.statusText}`, 'HTTP_ERROR');
  });

  if (!response.ok || !json.success) {
    const errorMsg = json.error?.message || `Request failed with status ${response.status}`;
    const errorCode = json.error?.code || 'UNKNOWN_ERROR';
    throw new ApiError(errorMsg, errorCode, json.error?.details);
  }

  return json.data;
}

export const api = {
  // Auth
  login: (email: string, password: string) =>
    request<{ token: string; staff: any; group: any }>('/admin/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),

  // Dashboard Stats
  getDashboardStats: () =>
    request<{
      total_revenue: number;
      mrr: number;
      arr: number;
      total_clients: number;
      active_orders: number;
      suspended_orders: number;
      pending_orders: number;
      unpaid_invoices: number;
      paid_invoices: number;
      open_tickets: number;
      closed_tickets: number;
    }>('/admin/stats/dashboard'),

  // Clients
  getClients: () => request<any[]>('/admin/clients'),

  // Orders
  getOrders: () => request<any[]>('/admin/orders'),
  activateOrder: (id: number) =>
    request<any>(`/admin/orders/${id}/activate`, { method: 'POST' }),
  suspendOrder: (id: number, reason = 'Administrative suspension') =>
    request<any>(`/admin/orders/${id}/suspend`, {
      method: 'POST',
      body: JSON.stringify({ reason }),
    }),
  unsuspendOrder: (id: number) =>
    request<any>(`/admin/orders/${id}/unsuspend`, { method: 'POST' }),

  // Support
  getSupportTickets: () => request<any[]>('/admin/support/tickets'),
  replySupportTicket: (id: number, content: string) =>
    request<any>(`/admin/support/tickets/${id}/reply`, {
      method: 'POST',
      body: JSON.stringify({ content }),
    }),

  // Currencies
  getCurrencies: () => request<any[]>('/admin/currencies'),
  createCurrency: (currency: {
    code: string;
    title: string;
    conversion_rate: number;
    format: string;
  }) =>
    request<any>('/admin/currencies', {
      method: 'POST',
      body: JSON.stringify(currency),
    }),
  setDefaultCurrency: (code: string) =>
    request<any>(`/admin/currencies/${code}/default`, { method: 'POST' }),
  deleteCurrency: (code: string) =>
    request<any>(`/admin/currencies/${code}`, { method: 'DELETE' }),

  // News
  getNews: () => request<any[]>('/admin/news'),
  createNews: (article: { title: string; content: string }) =>
    request<any>('/admin/news', {
      method: 'POST',
      body: JSON.stringify(article),
    }),
  deleteNews: (id: number) =>
    request<any>(`/admin/news/${id}`, { method: 'DELETE' }),

  // Mass Mailer
  getMassMailCampaigns: () => request<any[]>('/admin/mass-mail'),
  createMassMailCampaign: (campaign: {
    subject: string;
    content: string;
    target_group?: string;
  }) =>
    request<any>('/admin/mass-mail', {
      method: 'POST',
      body: JSON.stringify(campaign),
    }),
  sendMassMailCampaign: (id: number) =>
    request<any>(`/admin/mass-mail/${id}/send`, { method: 'POST' }),

  // Audit Logs
  getAuditLogs: () => request<any[]>('/admin/audit-logs'),
};
