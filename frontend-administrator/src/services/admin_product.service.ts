import { AdminProductRepository, adminProductRepository, IAdminProductRepository } from '../repositories/admin_product.repository';
import type { ProductItem, ProductCategory } from '@/types/api';

export class AdminProductService {
  constructor(private repo: IAdminProductRepository = adminProductRepository) {}

  async listProducts(): Promise<ProductItem[]> {
    return this.repo.getProducts();
  }

  async getProduct(id: number): Promise<ProductItem> {
    if (!id || id <= 0) {
      throw new Error('Valid product ID is required');
    }
    return this.repo.getProduct(id);
  }

  async createProduct(dto: Partial<ProductItem>): Promise<ProductItem> {
    if (!dto.title || !dto.title.trim()) {
      throw new Error('Product title is required');
    }
    const slug = dto.slug || dto.title.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)+/g, '');
    return this.repo.createProduct({
      ...dto,
      slug,
      price_monthly: Number(dto.price_monthly) || 0,
      price_annually: Number(dto.price_annually) || 0,
      setup_fee: Number(dto.setup_fee) || 0,
    });
  }

  async updateProduct(id: number, dto: Partial<ProductItem>): Promise<ProductItem> {
    if (!id || id <= 0) {
      throw new Error('Valid product ID is required');
    }
    return this.repo.updateProduct(id, dto);
  }

  async deleteProduct(id: number): Promise<any> {
    if (!id || id <= 0) {
      throw new Error('Valid product ID is required');
    }
    return this.repo.deleteProduct(id);
  }

  async listCategories(): Promise<ProductCategory[]> {
    return this.repo.getProductCategories();
  }

  filterProducts(products: ProductItem[], query: string, type: string): ProductItem[] {
    const q = query.toLowerCase().trim();
    return products.filter((p) => {
      const matchesSearch =
        !q ||
        p.title.toLowerCase().includes(q) ||
        p.description.toLowerCase().includes(q) ||
        p.slug.toLowerCase().includes(q);
      const matchesType = !type || type === 'all' || p.type === type;
      return matchesSearch && matchesType;
    });
  }
}

export const adminProductService = new AdminProductService();
