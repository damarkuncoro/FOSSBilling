import { AdminCouponRepository, adminCouponRepository, IAdminCouponRepository } from '../repositories/admin_coupon.repository';
import type { CouponItem } from '@/types/api';

export class AdminCouponService {
  constructor(private repo: IAdminCouponRepository = adminCouponRepository) {}

  async listCoupons(): Promise<CouponItem[]> {
    return this.repo.getCoupons();
  }

  async createCoupon(dto: Partial<CouponItem>): Promise<CouponItem> {
    if (!dto.code || !dto.code.trim()) {
      throw new Error('Coupon code is required');
    }
    if (!dto.value || Number(dto.value) <= 0) {
      throw new Error('Coupon discount value must be greater than zero');
    }
    return this.repo.createCoupon({
      ...dto,
      code: dto.code.trim().toUpperCase(),
      type: dto.type || 'percentage',
      value: Number(dto.value),
      max_uses: Number(dto.max_uses) || 0,
      is_active: dto.is_active ?? true,
    });
  }

  async deleteCoupon(id: number): Promise<any> {
    if (!id || id <= 0) {
      throw new Error('Valid coupon ID is required');
    }
    return this.repo.deleteCoupon(id);
  }

  generateRandomCode(prefix = 'PROMO'): string {
    const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789';
    let code = prefix;
    for (let i = 0; i < 5; i++) {
      code += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return code;
  }
}

export const adminCouponService = new AdminCouponService();
