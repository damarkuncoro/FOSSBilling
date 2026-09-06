import { request } from '../lib/api/client';
import type { ProductItem, ProductCategory } from '@/types/api';

export interface IAdminProductRepository {
  getProducts(): Promise<ProductItem[]>;
  getProduct(id: number): Promise<ProductItem>;
  createProduct(dto: Partial<ProductItem>): Promise<ProductItem>;
  updateProduct(id: number, dto: Partial<ProductItem>): Promise<ProductItem>;
  deleteProduct(id: number): Promise<any>;
  getProductCategories(): Promise<ProductCategory[]>;
}

export class AdminProductRepository implements IAdminProductRepository {
  async getProducts(): Promise<ProductItem[]> {
    return request<ProductItem[]>('/admin/products');
  }

  async getProduct(id: number): Promise<ProductItem> {
    return request<ProductItem>(`/admin/products/${id}`);
  }

  async createProduct(dto: Partial<ProductItem>): Promise<ProductItem> {
    return request<ProductItem>('/admin/products', {
      method: 'POST',
      body: JSON.stringify(dto),
    });
  }

  async updateProduct(id: number, dto: Partial<ProductItem>): Promise<ProductItem> {
    return request<ProductItem>(`/admin/products/${id}`, {
      method: 'PUT',
      body: JSON.stringify(dto),
    });
  }

  async deleteProduct(id: number): Promise<any> {
    return request<any>(`/admin/products/${id}`, {
      method: 'DELETE',
    });
  }

  async getProductCategories(): Promise<ProductCategory[]> {
    return request<ProductCategory[]>('/admin/products/categories');
  }
}

export const adminProductRepository = new AdminProductRepository();
