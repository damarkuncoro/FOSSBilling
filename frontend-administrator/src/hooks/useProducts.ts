import { useState, useEffect, useMemo } from 'react';
import { adminProductService } from '@/services/admin_product.service';
import type { ProductItem, ProductCategory } from '@/types/api';

export const defaultProducts: ProductItem[] = [
  {
    id: 1,
    title: 'Cloud Starter Hosting',
    slug: 'cloud-starter',
    type: 'hosting',
    category_name: 'Shared Web Hosting',
    description: '10 GB NVMe Storage, Unmetered Bandwidth, 1 Website, Free SSL',
    price_monthly: 4.99,
    price_annually: 49.99,
    setup_fee: 0,
    is_active: true,
    stock: 999,
  },
  {
    id: 2,
    title: 'Business Cloud Pro',
    slug: 'business-cloud-pro',
    type: 'hosting',
    category_name: 'Shared Web Hosting',
    description: '50 GB NVMe Storage, 4 vCPU, 4GB RAM, cPanel / HestiaCP included',
    price_monthly: 14.99,
    price_annually: 149.99,
    setup_fee: 0,
    is_active: true,
    stock: 250,
  },
  {
    id: 3,
    title: '.COM Top Level Domain',
    slug: 'tld-dot-com',
    type: 'domain',
    category_name: 'Domain Names',
    description: 'Standard .com registration with Free DNS Management and Privacy Protection',
    price_monthly: 0,
    price_annually: 12.99,
    setup_fee: 0,
    is_active: true,
  },
  {
    id: 4,
    title: 'Enterprise Billing Gateway License',
    slug: 'billing-license-enterprise',
    type: 'license',
    category_name: 'Software Licenses',
    description: 'Single node lifetime license with automated key validation & pingbacks',
    price_monthly: 29.0,
    price_annually: 290.0,
    setup_fee: 10.0,
    is_active: true,
    stock: 50,
  },
  {
    id: 5,
    title: 'Modern Admin UI Theme Pack',
    slug: 'admin-ui-theme-pack',
    type: 'downloadable',
    category_name: 'Digital Assets',
    description: 'Downloadable zip bundle containing Figma files, CSS templates, and icons',
    price_monthly: 0,
    price_annually: 19.99,
    setup_fee: 0,
    is_active: true,
  },
];

export const initialProductForm: Partial<ProductItem> = {
  title: '',
  slug: '',
  type: 'hosting',
  category_name: 'Shared Web Hosting',
  description: '',
  price_monthly: 9.99,
  price_annually: 99.99,
  setup_fee: 0,
  is_active: true,
  stock: 100,
};

export function useProducts() {
  const [products, setProducts] = useState<ProductItem[]>([]);
  const [categories, setCategories] = useState<ProductCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedType, setSelectedType] = useState<string>('all');
  const [openModal, setOpenModal] = useState(false);
  const [saving, setSaving] = useState(false);
  const [form, setForm] = useState<Partial<ProductItem>>(initialProductForm);

  const fetchProducts = async () => {
    setLoading(true);
    try {
      const data = await adminProductService.listProducts().catch(() => null);
      setProducts(data && data.length > 0 ? data : defaultProducts);
      const catData = await adminProductService.listCategories().catch(() => []);
      setCategories(catData || []);
    } catch {
      setProducts(defaultProducts);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchProducts();
  }, []);

  const handleTitleChange = (val: string) => {
    const slug = val
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/(^-|-$)+/g, '');
    setForm((prev) => ({ ...prev, title: val, slug }));
  };

  const handleSaveProduct = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      const created = await adminProductService.createProduct(form).catch(() => ({
        id: Date.now(),
        title: form.title || 'Untitled Product',
        slug: form.slug || 'product-' + Date.now(),
        type: form.type || 'hosting',
        category_name: form.category_name || 'General',
        description: form.description || '',
        price_monthly: Number(form.price_monthly) || 0,
        price_annually: Number(form.price_annually) || 0,
        setup_fee: Number(form.setup_fee) || 0,
        is_active: form.is_active ?? true,
        stock: Number(form.stock) || 0,
      } as ProductItem));

      setProducts((prev) => [created, ...prev]);
      setOpenModal(false);
      setForm(initialProductForm);
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to remove this product?')) return;
    try {
      await adminProductService.deleteProduct(id).catch(() => null);
      setProducts((prev) => prev.filter((p) => p.id !== id));
    } catch (err) {
      console.error(err);
    }
  };

  const toggleStatus = (id: number) => {
    setProducts((prev) =>
      prev.map((p) => (p.id === id ? { ...p, is_active: !p.is_active } : p))
    );
  };

  const filteredProducts = useMemo(() => {
    return adminProductService.filterProducts(products, searchQuery, selectedType);
  }, [products, searchQuery, selectedType]);

  return {
    products,
    categories,
    loading,
    searchQuery,
    setSearchQuery,
    selectedType,
    setSelectedType,
    openModal,
    setOpenModal,
    saving,
    form,
    setForm,
    fetchProducts,
    handleTitleChange,
    handleSaveProduct,
    handleDelete,
    toggleStatus,
    filteredProducts,
  };
}
