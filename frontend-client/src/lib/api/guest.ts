import { request } from './client';
import { ClientProfile, PublicCompanyInfo, CartCalculation } from '@/types/api';

export const guestApi = {
  // Guest / Public Endpoints
  guestCurrencies: () => request<any[]>('/guest/currencies'),
  guestNews: () => request<any[]>('/guest/news'),
  guestNewsBySlug: (slug: string) => request<any>(`/guest/news/${slug}`),
  checkDomain: (domain: string) =>
    request<{ domain: string; tld: string; available: boolean; price: number; currency: string }>(
      `/guest/domains/check?domain=${encodeURIComponent(domain)}`
    ),

  // Guest Auth
  login: (email: string, password: string) =>
    request<{ token: string; client: ClientProfile; two_factor_required?: boolean }>('/guest/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),
  verifyTwoFactor: (email: string, code: string) =>
    request<{ token: string; client: ClientProfile }>('/guest/auth/verify-2fa', {
      method: 'POST',
      body: JSON.stringify({ email, code }),
    }),
  register: (dto: {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    currency?: string;
  }) =>
    request<{ token: string; client: ClientProfile }>('/guest/auth/register', {
      method: 'POST',
      body: JSON.stringify(dto),
    }),

  // Cart & Checkout
  calculateCart: (items: any[], promoCode?: string) =>
    request<CartCalculation>('/guest/cart/calculate', {
      method: 'POST',
      body: JSON.stringify({ items, promo_code: promoCode }),
    }),
  checkoutCart: (payload: {
    client_id: number;
    items: any[];
    promo_code?: string;
    gateway?: string;
  }) =>
    request<{ invoice?: { id: number }; invoice_id: number; order_ids: number[]; total: number }>(
      '/guest/cart/checkout',
      {
        method: 'POST',
        body: JSON.stringify(payload),
      }
    ),

  // Public Company Info & Branding
  getCompany: () => request<PublicCompanyInfo>('/guest/company'),

  // Knowledgebase
  getKBCategories: () => request<any[]>('/guest/kb/categories'),
  getKBArticles: (catID?: number) =>
    request<any[]>(`/guest/kb/articles${catID ? `?category_id=${catID}` : ''}`),
  getKBArticle: (slug: string) =>
    request<any>(`/guest/kb/articles/detail?slug=${slug}`),

  // Forms
  getForm: (id: number) => request<any>(`/guest/forms/${id}`),

  // Redirects
  lookupRedirect: (path: string) => request<{ path: string; target: string; status_code: number }>(`/guest/redirects/lookup?path=${encodeURIComponent(path)}`),
};
