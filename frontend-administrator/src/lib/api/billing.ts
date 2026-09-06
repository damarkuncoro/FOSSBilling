import { request } from './client';
import {
  PaymentGatewayItem,
  TaxRuleItem,
  CouponItem,
  EmailTemplateItem,
  MailConfig,
} from '@/types/api';
import { FinancialReportSummary } from '@/types/modules';

export const billingApi = {
  getInvoices: () => request<any[]>('/admin/invoices'),
  createInvoice: (data: any) => request<any>('/admin/invoices', { method: 'POST', body: JSON.stringify(data) }),
  getCurrencies: () => request<any[]>('/admin/currencies'),
  createCurrency: (currency: any) => request<any>('/admin/currencies', { method: 'POST', body: JSON.stringify(currency) }),
  setDefaultCurrency: (code: string) => request<any>(`/admin/currencies/${code}/default`, { method: 'POST' }),
  deleteCurrency: (code: string) => request<any>(`/admin/currencies/${code}`, { method: 'DELETE' }),
  getPaymentGateways: () => request<PaymentGatewayItem[]>('/admin/gateways'),
  updatePaymentGateway: (id: string, config: Partial<PaymentGatewayItem>) =>
    request<PaymentGatewayItem>(`/admin/gateways/${id}`, { method: 'PUT', body: JSON.stringify(config) }),
  getTaxRules: () => request<TaxRuleItem[]>('/admin/tax-rules'),
  createTaxRule: (rule: Partial<TaxRuleItem>) => request<TaxRuleItem>('/admin/tax-rules', { method: 'POST', body: JSON.stringify(rule) }),
  deleteTaxRule: (id: number) => request<any>(`/admin/tax-rules/${id}`, { method: 'DELETE' }),
  getCoupons: () => request<CouponItem[]>('/admin/coupons'),
  createCoupon: (coupon: Partial<CouponItem>) => request<CouponItem>('/admin/coupons', { method: 'POST', body: JSON.stringify(coupon) }),
  deleteCoupon: (id: number) => request<any>(`/admin/coupons/${id}`, { method: 'DELETE' }),
  getEmailTemplates: () => request<EmailTemplateItem[]>('/admin/email-templates'),
  updateEmailTemplate: (id: string, template: Partial<EmailTemplateItem>) =>
    request<EmailTemplateItem>(`/admin/email-templates/${id}`, { method: 'PUT', body: JSON.stringify(template) }),
  getMailConfig: () => request<MailConfig>('/admin/settings/mail'),
  updateMailConfig: (config: Partial<MailConfig>) => request<MailConfig>('/admin/settings/mail', { method: 'PUT', body: JSON.stringify(config) }),
  sendTestEmail: (toEmail: string) =>
    request<{ success: boolean; message: string }>('/admin/settings/mail/test', {
      method: 'POST',
      body: JSON.stringify({ email: toEmail }),
    }),
  getFinancialReports: () => request<FinancialReportSummary>('/admin/reports/financial'),
};
