import { request } from '../lib/api/client';
import type { PaymentGatewayItem, TaxRuleItem } from '@/types/api';

export interface IAdminGatewayRepository {
  getPaymentGateways(): Promise<PaymentGatewayItem[]>;
  updatePaymentGateway(id: string, dto: Partial<PaymentGatewayItem>): Promise<PaymentGatewayItem>;
  getTaxRules(): Promise<TaxRuleItem[]>;
  createTaxRule(dto: Partial<TaxRuleItem>): Promise<TaxRuleItem>;
  deleteTaxRule(id: number): Promise<any>;
}

export class AdminGatewayRepository implements IAdminGatewayRepository {
  async getPaymentGateways(): Promise<PaymentGatewayItem[]> {
    return request<PaymentGatewayItem[]>('/admin/gateways');
  }

  async updatePaymentGateway(id: string, dto: Partial<PaymentGatewayItem>): Promise<PaymentGatewayItem> {
    return request<PaymentGatewayItem>(`/admin/gateways/${id}`, {
      method: 'PUT',
      body: JSON.stringify(dto),
    });
  }

  async getTaxRules(): Promise<TaxRuleItem[]> {
    return request<TaxRuleItem[]>('/admin/taxes');
  }

  async createTaxRule(dto: Partial<TaxRuleItem>): Promise<TaxRuleItem> {
    return request<TaxRuleItem>('/admin/taxes', {
      method: 'POST',
      body: JSON.stringify(dto),
    });
  }

  async deleteTaxRule(id: number): Promise<any> {
    return request<any>(`/admin/taxes/${id}`, {
      method: 'DELETE',
    });
  }
}

export const adminGatewayRepository = new AdminGatewayRepository();
