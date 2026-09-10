import { request } from '../lib/api/client';
import type { ProductItem, ProductCategory } from '@/types/api';

export class AdminProductRepository {
  list = () => request<ProductItem[]>('/admin/products');
  get = (id: number) => request<ProductItem>(`/admin/products/${id}`);
  create = (d: any) => request<ProductItem>('/admin/products', { method: 'POST', body: JSON.stringify(d) });
  update = (id: number, d: any) => request<ProductItem>(`/admin/products/${id}`, { method: 'PUT', body: JSON.stringify(d) });
  delete = (id: number) => request(`/admin/products/${id}`, { method: 'DELETE' });
  categories = () => request<ProductCategory[]>('/admin/products/categories');
}

export const adminProductRepository = new AdminProductRepository();
