import { request } from './client';
import {
  ProductItem,
  ProductCategory,
  ServerItem,
} from '@/types/api';
import {
  TldPricingItem,
  RegistrarConfig,
} from '@/types/modules';

export const catalogApi = {
  getProducts: () => request<ProductItem[]>('/admin/products'),
  createProduct: (product: Partial<ProductItem>) => request<ProductItem>('/admin/products', { method: 'POST', body: JSON.stringify(product) }),
  updateProduct: (id: number, product: Partial<ProductItem>) => request<ProductItem>(`/admin/products/${id}`, { method: 'PUT', body: JSON.stringify(product) }),
  deleteProduct: (id: number) => request<any>(`/admin/products/${id}`, { method: 'DELETE' }),
  getProductCategories: () => request<ProductCategory[]>('/admin/product-categories'),
  getServers: () => request<ServerItem[]>('/admin/servers'),
  createServer: (server: Partial<ServerItem>) => request<ServerItem>('/admin/servers', { method: 'POST', body: JSON.stringify(server) }),
  updateServer: (id: number, server: Partial<ServerItem>) => request<ServerItem>(`/admin/servers/${id}`, { method: 'PUT', body: JSON.stringify(server) }),
  testServerConnection: (id: number) => request<{ success: boolean; message: string }>(`/admin/servers/${id}/test`, { method: 'POST' }),
  deleteServer: (id: number) => request<any>(`/admin/servers/${id}`, { method: 'DELETE' }),
  getTlds: () => request<TldPricingItem[]>('/admin/domains/tlds'),
  createTld: (tld: Partial<TldPricingItem>) => request<TldPricingItem>('/admin/domains/tlds', { method: 'POST', body: JSON.stringify(tld) }),
  deleteTld: (id: number) => request<any>(`/admin/domains/tlds/${id}`, { method: 'DELETE' }),
  getRegistrars: () => request<RegistrarConfig[]>('/admin/domains/registrars'),
  updateRegistrar: (id: string, config: Partial<RegistrarConfig>) =>
    request<RegistrarConfig>(`/admin/domains/registrars/${id}`, { method: 'PUT', body: JSON.stringify(config) }),
};
