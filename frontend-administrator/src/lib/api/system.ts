import { request } from './client';
import {
  StaffMemberItem,
  SecuritySettings,
  SystemStatusInfo,
  CompanySettings,
} from '@/types/api';
import {
  CustomPageItem,
  KnowledgebaseArticle,
  ExtensionModuleItem,
} from '@/types/modules';

export const systemApi = {
  login: (email: string, password: string) =>
    request<{ token: string; staff: any; group: any }>('/admin/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),
  getDashboardStats: () => request<any>('/admin/stats/dashboard'),
  getSupportTickets: () => request<any[]>('/admin/support/tickets'),
  replySupportTicket: (id: number, content: string) =>
    request<any>(`/admin/support/tickets/${id}/reply`, { method: 'POST', body: JSON.stringify({ content }) }),
  getStaffMembers: () => request<StaffMemberItem[]>('/admin/staff'),
  createStaffMember: (staff: Partial<StaffMemberItem> & { password?: string }) =>
    request<StaffMemberItem>('/admin/staff', { method: 'POST', body: JSON.stringify(staff) }),
  deleteStaffMember: (id: number) => request<any>(`/admin/staff/${id}`, { method: 'DELETE' }),
  getSecuritySettings: () => request<SecuritySettings>('/admin/settings/security'),
  updateSecuritySettings: (settings: Partial<SecuritySettings>) =>
    request<SecuritySettings>('/admin/settings/security', { method: 'PUT', body: JSON.stringify(settings) }),
  getSystemStatus: () => request<SystemStatusInfo>('/admin/system/status'),
  triggerCron: () => request<{ success: boolean; message: string; timestamp: string }>('/admin/system/cron/run', { method: 'POST' }),
  clearSystemCache: () => request<{ success: boolean; message: string }>('/admin/system/cache/clear', { method: 'POST' }),
  getNews: () => request<any[]>('/admin/news'),
  createNews: (article: { title: string; content: string }) => request<any>('/admin/news', { method: 'POST', body: JSON.stringify(article) }),
  deleteNews: (id: number) => request<any>(`/admin/news/${id}`, { method: 'DELETE' }),
  getMassMailCampaigns: () => request<any[]>('/admin/mass-mail'),
  createMassMailCampaign: (campaign: any) => request<any>('/admin/mass-mail', { method: 'POST', body: JSON.stringify(campaign) }),
  sendMassMailCampaign: (id: number) => request<any>(`/admin/mass-mail/${id}/send`, { method: 'POST' }),
  getCompany: () => request<CompanySettings>('/admin/company'),
  updateCompany: (settings: Partial<CompanySettings>) => request<CompanySettings>('/admin/company', { method: 'PUT', body: JSON.stringify(settings) }),
  getAuditLogs: () => request<any[]>('/admin/audit-logs'),
  getPages: () => request<CustomPageItem[]>('/admin/pages'),
  savePage: (page: Partial<CustomPageItem>) => request<CustomPageItem>('/admin/pages', { method: 'POST', body: JSON.stringify(page) }),
  deletePage: (id: number) => request<any>(`/admin/pages/${id}`, { method: 'DELETE' }),
  getKnowledgebase: () => request<KnowledgebaseArticle[]>('/admin/knowledgebase'),
  saveKnowledgebase: (article: Partial<KnowledgebaseArticle>) =>
    request<KnowledgebaseArticle>('/admin/knowledgebase', { method: 'POST', body: JSON.stringify(article) }),
  deleteKnowledgebase: (id: number) => request<any>(`/admin/knowledgebase/${id}`, { method: 'DELETE' }),
  getExtensions: () => request<ExtensionModuleItem[]>('/admin/extensions'),
  toggleExtension: (id: string, enabled: boolean) =>
    request<any>(`/admin/extensions/${id}/toggle`, { method: 'POST', body: JSON.stringify({ enabled }) }),
};
