import { request } from '../lib/api/client';
import type { CouponItem } from '@/types/api';

export interface IAdminCouponRepository {
  getCoupons(): Promise<CouponItem[]>;
  createCoupon(dto: Partial<CouponItem>): Promise<CouponItem>;
  deleteCoupon(id: number): Promise<any>;
}

export class AdminCouponRepository implements IAdminCouponRepository {
  async getCoupons(): Promise<CouponItem[]> {
    return request<CouponItem[]>('/admin/coupons');
  }

  async createCoupon(dto: Partial<CouponItem>): Promise<CouponItem> {
    return request<CouponItem>('/admin/coupons', {
      method: 'POST',
      body: JSON.stringify(dto),
    });
  }

  async deleteCoupon(id: number): Promise<any> {
    return request<any>(`/admin/coupons/${id}`, {
      method: 'DELETE',
    });
  }
}

export const adminCouponRepository = new AdminCouponRepository();
