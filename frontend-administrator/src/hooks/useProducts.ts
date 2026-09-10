import { useState, useMemo } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminProductService } from '@/services/admin_product.service';
import { adminFormBuilderService } from '@/services/admin_form_builder.service';
import type { ProductItem } from '@/types/api';

export const defaultProducts: ProductItem[] = [
  { id: 1, title: 'Cloud Starter Hosting', slug: 'cloud-starter', type: 'hosting', category_name: 'Shared Web Hosting', description: '10 GB NVMe Storage, Unmetered Bandwidth, 1 Website, Free SSL', price_monthly: 4.99, price_annually: 49.99, setup_fee: 0, is_active: true, stock: 999 },
  { id: 2, title: 'Business Cloud Pro', slug: 'business-cloud-pro', type: 'hosting', category_name: 'Shared Web Hosting', description: '50 GB NVMe Storage, 4 vCPU, 4GB RAM, cPanel / HestiaCP included', price_monthly: 14.99, price_annually: 149.99, setup_fee: 0, is_active: true, stock: 250 },
];

export const initialProductForm: Partial<ProductItem> = {
  title: '', slug: '', type: 'hosting', category_name: 'Shared Web Hosting', description: '', price_monthly: 9.99, price_annually: 99.99, setup_fee: 0, is_active: true, stock: 100,
};

export function useProducts() {
  const queryClient = useQueryClient();
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedType, setSelectedType] = useState<string>('all');
  const [openModal, setOpenModal] = useState(false);
  const [form, setForm] = useState<Partial<ProductItem>>(initialProductForm);

  const { data: products = [], isLoading: productsLoading, refetch: fetchProducts } = useQuery({
    queryKey: ['admin', 'products'],
    queryFn: async () => {
      const data = await adminProductService.listProducts();
      return data && data.length > 0 ? data : defaultProducts;
    },
  });

  const { data: categories = [] } = useQuery({
    queryKey: ['admin', 'products', 'categories'],
    queryFn: () => adminProductService.listCategories(),
  });

  const { data: availableForms = [] } = useQuery({
    queryKey: ['admin', 'forms'],
    queryFn: () => adminFormBuilderService.listForms(),
  });

  const createMutation = useMutation({
    mutationFn: (input: Partial<ProductItem>) => adminProductService.createProduct(input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin', 'products'] });
      setOpenModal(false);
      setForm(initialProductForm);
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, input }: { id: number; input: Partial<ProductItem> }) =>
      adminProductService.updateProduct(id, input),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'products'] }),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => adminProductService.deleteProduct(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['admin', 'products'] }),
  });

  const handleTitleChange = (val: string) => {
    const slug = val.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)+/g, '');
    setForm((prev) => ({ ...prev, title: val, slug }));
  };

  const filteredProducts = useMemo(() => {
    return adminProductService.filterProducts(products, searchQuery, selectedType);
  }, [products, searchQuery, selectedType]);

  return {
    products,
    categories,
    availableForms,
    loading: productsLoading,
    searchQuery,
    setSearchQuery,
    selectedType,
    setSelectedType,
    openModal,
    setOpenModal,
    saving: createMutation.isPending || updateMutation.isPending || deleteMutation.isPending,
    form,
    setForm,
    fetchProducts,
    handleTitleChange,
    handleSaveProduct: (e: React.FormEvent) => {
      e.preventDefault();
      createMutation.mutate(form);
    },
    handleDelete: (id: number) => {
      if (confirm('Are you sure?')) deleteMutation.mutate(id);
    },
    toggleStatus: (id: number) => {
      const p = products.find((x) => x.id === id);
      if (p) updateMutation.mutate({ id, input: { is_active: !p.is_active } });
    },
    filteredProducts,
  };
}
