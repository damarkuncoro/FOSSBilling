import { request } from '../lib/api/client';
import type { CartCalculation } from '@/types/api';

export interface ICartRepository {
  calculateCart(items: any[], promoCode?: string): Promise<CartCalculation>;
  checkoutCart(payload: {
    client_id: number;
    items: any[];
    promo_code?: string;
    gateway?: string;
  }): Promise<{ invoice_id: number; order_ids: number[]; total: number }>;
}

export class CartRepository implements ICartRepository {
  async calculateCart(items: any[], promoCode?: string): Promise<CartCalculation> {
    return request<CartCalculation>('/guest/cart/calculate', {
      method: 'POST',
      body: JSON.stringify({ items, promo_code: promoCode }),
    });
  }

  async checkoutCart(payload: {
    client_id: number;
    items: any[];
    promo_code?: string;
    gateway?: string;
  }): Promise<{ invoice_id: number; order_ids: number[]; total: number }> {
    return request<{ invoice_id: number; order_ids: number[]; total: number }>('/guest/cart/checkout', {
      method: 'POST',
      body: JSON.stringify(payload),
    });
  }
}

export const cartRepository = new CartRepository();
