import { request } from './client';
import {
  ClientProfile,
  Order,
  Invoice,
  SupportTicket,
  ApiKey,
  Notification,
} from '@/types/api';

export const clientPortalApi = {
  // Client Profile
  getProfile: () =>
    request<{ client: ClientProfile; balance: number }>('/client/profile'),
  updateProfile: (profile: Partial<ClientProfile>) =>
    request<{ client: ClientProfile }>('/client/profile', {
      method: 'PUT',
      body: JSON.stringify(profile),
    }),
  changePassword: (dto: { current_password: string; new_password: string }) =>
    request<{ success: boolean; message: string }>('/client/profile/change-password', {
      method: 'POST',
      body: JSON.stringify(dto),
    }),
  setupTwoFactor: () =>
    request<{ secret: string; qr_url: string }>('/client/profile/2fa/setup', {
      method: 'POST',
    }),
  enableTwoFactor: (code: string) =>
    request<any>('/client/profile/2fa/enable', {
      method: 'POST',
      body: JSON.stringify({ code }),
    }),
  disableTwoFactor: () =>
    request<any>('/client/profile/2fa/disable', {
      method: 'POST',
    }),

  // Client Orders & Services
  getOrders: () => request<Order[]>('/client/orders'),
  getOrder: (id: number) => request<Order>(`/client/orders/${id}`),
  syncServiceStatus: (id: number) =>
    request<any>(`/client/orders/${id}/sync`, { method: 'POST' }),
  changeServicePassword: (id: number, password: string) =>
    request<any>(`/client/orders/${id}/change-password`, {
      method: 'POST',
      body: JSON.stringify({ password }),
    }),

  // Client Domains
  getDomains: () => request<any[]>('/client/domains'),
  updateNameservers: (id: number, nameservers: string[]) =>
    request<any>(`/client/domains/${id}/nameservers`, {
      method: 'PUT',
      body: JSON.stringify({ nameservers }),
    }),
  toggleDomainAutoRenew: (id: number) =>
    request<any>(`/client/domains/${id}/toggle-autorenew`, {
      method: 'POST',
    }),

  // Client Invoices & Funds
  getInvoices: () => request<Invoice[]>('/client/invoices'),
  getInvoice: (id: number) => request<Invoice>(`/client/invoices/${id}`),
  payWithBalance: (id: number) =>
    request<any>(`/client/invoices/${id}/pay-balance`, { method: 'POST' }),
  depositFunds: (amount: number, currency = 'USD', gateway = '') =>
    request<any>('/client/funds/deposit', {
      method: 'POST',
      body: JSON.stringify({ amount, currency, gateway }),
    }),

  // Client Support
  getTickets: () => request<SupportTicket[]>('/client/support/tickets'),
  getTicket: (id: number) => request<SupportTicket>(`/client/support/tickets/${id}`),
  openTicket: (ticket: { subject: string; message: string; priority?: string }) =>
    request<SupportTicket>('/client/support/tickets', {
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

  // Software Licenses
  getLicenses: () => request<any[]>('/client/licenses'),
  resetLicenseLock: (id: number) =>
    request<any>(`/client/licenses/${id}/reset`, { method: 'POST' }),

  // Digital Downloads
  getDownloads: () => request<any[]>('/client/downloads'),
  getDownloadLink: (id: number) =>
    request<{ download_url: string; expires_at: number }>(
      `/client/downloads/${id}/link`
    ),

  // API Keys
  getApiKeys: () => request<ApiKey[]>('/client/api-keys'),
  generateApiKey: (name: string) =>
    request<ApiKey>('/client/api-keys', {
      method: 'POST',
      body: JSON.stringify({ name }),
    }),
  revokeApiKey: (id: number) =>
    request<any>(`/client/api-keys/${id}`, { method: 'DELETE' }),

  // Notifications
  getNotifications: () => request<Notification[]>('/client/notifications'),
  markNotificationRead: (id: number) =>
    request<any>(`/client/notifications/${id}/read`, { method: 'PUT' }),
  markAllNotificationsRead: () =>
    request<any>('/client/notifications/mark-all-read', { method: 'POST' }),
};
