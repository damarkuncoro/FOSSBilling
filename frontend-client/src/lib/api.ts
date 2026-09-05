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

export function getStoredClientToken(): string | null {
  return localStorage.getItem('fossbilling_client_token');
}

export function setStoredClientToken(token: string) {
  localStorage.setItem('fossbilling_client_token', token);
}

export function removeStoredClientToken() {
  localStorage.removeItem('fossbilling_client_token');
  localStorage.removeItem('fossbilling_client_user');
}

async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const token = getStoredClientToken();
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
  // Guest / Public Endpoints
  guestCurrencies: () => request<any[]>('/guest/currencies'),
  guestNews: () => request<any[]>('/guest/news'),
  guestNewsBySlug: (slug: string) => request<any>(`/guest/news/${slug}`),

  // Guest Auth
  login: (email: string, password: string) =>
    request<{ token: string; client: any }>('/guest/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),
  register: (dto: {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    currency?: string;
  }) =>
    request<{ token: string; client: any }>('/guest/auth/register', {
      method: 'POST',
      body: JSON.stringify(dto),
    }),

  // Cart & Checkout
  calculateCart: (items: any[], promoCode?: string) =>
    request<{
      subtotal: number;
      discount: number;
      tax: number;
      total: number;
      items: any[];
    }>('/guest/cart/calculate', {
      method: 'POST',
      body: JSON.stringify({ items, promo_code: promoCode }),
    }),
  checkoutCart: (payload: {
    client_id: number;
    items: any[];
    promo_code?: string;
    gateway?: string;
  }) =>
    request<{ invoice_id: number; order_ids: number[]; total: number }>(
      '/guest/cart/checkout',
      {
        method: 'POST',
        body: JSON.stringify(payload),
      }
    ),

  // Client Profile
  getProfile: () =>
    request<{ client: any; balance: number }>('/client/profile'),
  updateProfile: (profile: any) =>
    request<{ client: any }>('/client/profile', {
      method: 'PUT',
      body: JSON.stringify(profile),
    }),

  // Client Orders & Services
  getOrders: () => request<any[]>('/client/orders'),
  getOrder: (id: number) => request<any>(`/client/orders/${id}`),

  // Client Invoices
  getInvoices: () => request<any[]>('/client/invoices'),
  getInvoice: (id: number) => request<any>(`/client/invoices/${id}`),
  payWithBalance: (id: number) =>
    request<any>(`/client/invoices/${id}/pay-balance`, { method: 'POST' }),

  // Client Support
  getTickets: () => request<any[]>('/client/support/tickets'),
  getTicket: (id: number) => request<any>(`/client/support/tickets/${id}`),
  openTicket: (ticket: { subject: string; message: string; priority?: string }) =>
    request<any>('/client/support/tickets', {
      method: 'POST',
      body: JSON.stringify(ticket),
    }),
  replyTicket: (id: number, content: string) =>
    request<any>(`/client/support/tickets/${id}/reply`, {
      method: 'POST',
      body: JSON.stringify({ content }),
    }),
  closeTicket: (id: number) =>
    request<any>(`/client/support/tickets/${id}/close`, { method: 'POST' }),

  // Digital Downloads
  getDownloadLink: (id: number) =>
    request<{ download_url: string; expires_at: number }>(
      `/client/downloads/${id}/link`
    ),

  // API Keys
  getApiKeys: () => request<any[]>('/client/api-keys'),
  generateApiKey: (name: string) =>
    request<any>('/client/api-keys', {
      method: 'POST',
      body: JSON.stringify({ name }),
    }),
  revokeApiKey: (id: number) =>
    request<any>(`/client/api-keys/${id}`, { method: 'DELETE' }),
};
