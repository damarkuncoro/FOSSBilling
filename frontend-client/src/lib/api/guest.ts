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
    request<{ token: string; client: ClientProfile }>('/guest/auth/login', {
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
    request<{ invoice_id: number; order_ids: number[]; total: number }>(
      '/guest/cart/checkout',
      {
        method: 'POST',
        body: JSON.stringify(payload),
      }
    ),

  // Public Company Info & Branding
  getCompany: () => request<PublicCompanyInfo>('/guest/company'),
};
