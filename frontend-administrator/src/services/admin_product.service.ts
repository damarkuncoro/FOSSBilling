import { adminProductRepository, AdminProductRepository } from '../repositories/admin_product.repository';
import type { ProductItem } from '@/types/api';

export class AdminProductService {
  constructor(private repo: AdminProductRepository = adminProductRepository) {}

  listProducts = () => this.repo.list();
  getProduct = (id: number) => this.repo.get(id);
  createProduct = (d: any) => {
    if (!d.title || !d.title.trim()) throw new Error('Product title is required');
    return this.repo.create(d);
  };
  updateProduct = (id: number, d: any) => this.repo.update(id, d);
  deleteProduct = (id: number) => this.repo.delete(id);
  listCategories = () => this.repo.categories();

  filterProducts(ps: ProductItem[], q: string, t: string) {
    const lq = q.toLowerCase();
    return ps.filter(p =>
      ((p.title || '').toLowerCase().includes(lq) || (p.description || '').toLowerCase().includes(lq)) &&
      (t === 'all' || p.type === t)
    );
  }
}

export const adminProductService = new AdminProductService();
