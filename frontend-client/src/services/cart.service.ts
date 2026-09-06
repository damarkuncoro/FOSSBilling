import { CartRepository, cartRepository, ICartRepository } from '../repositories/cart.repository';
import type { CartCalculation } from '@/types/api';

export class CartService {
  constructor(private repo: ICartRepository = cartRepository) {}

  async calculateCart(items: any[], promoCode?: string): Promise<CartCalculation> {
    return this.repo.calculateCart(items, promoCode);
  }

  async checkoutCart(payload: {
    client_id: number;
    items: any[];
    promo_code?: string;
    gateway?: string;
  }): Promise<{ invoice_id: number; order_ids: number[]; total: number }> {
    if (!payload.client_id) {
      throw new Error('Client ID is required for checkout');
    }
    if (!payload.items || payload.items.length === 0) {
      throw new Error('Cart cannot be empty');
    }
    return this.repo.checkoutCart(payload);
  }
}

export const cartService = new CartService();
