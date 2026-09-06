import { AdminGatewayRepository, adminGatewayRepository, IAdminGatewayRepository } from '../repositories/admin_gateway.repository';
import type { PaymentGatewayItem, TaxRuleItem } from '@/types/api';

export class AdminGatewayService {
  constructor(private repo: IAdminGatewayRepository = adminGatewayRepository) {}

  async listPaymentGateways(): Promise<PaymentGatewayItem[]> {
    return this.repo.getPaymentGateways();
  }

  async updatePaymentGateway(id: string, dto: Partial<PaymentGatewayItem>): Promise<PaymentGatewayItem> {
    if (!id || !id.trim()) {
      throw new Error('Gateway ID is required');
    }
    return this.repo.updatePaymentGateway(id, dto);
  }

  async listTaxRules(): Promise<TaxRuleItem[]> {
    return this.repo.getTaxRules();
  }

  async createTaxRule(dto: Partial<TaxRuleItem>): Promise<TaxRuleItem> {
    if (!dto.name || !dto.name.trim()) {
      throw new Error('Tax rule name is required');
    }
    return this.repo.createTaxRule({
      ...dto,
      country: dto.country || 'GLOBAL',
      rate: Number(dto.rate) || 0,
      is_active: dto.is_active ?? true,
      apply_to_all_clients: dto.apply_to_all_clients ?? false,
    });
  }

  async deleteTaxRule(id: number): Promise<any> {
    if (!id || id <= 0) {
      throw new Error('Valid tax rule ID is required');
    }
    return this.repo.deleteTaxRule(id);
  }
}

export const adminGatewayService = new AdminGatewayService();
